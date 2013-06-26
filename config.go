package dogo

import (
	"os"
	"io"
	"bufio"
	"strconv"
	"strings"
	"errors"
)

type Config struct {
	fileName string
	data map[string]map[string]string
}

func NewConfig(fileName string) (c *Config, err error) {
	var file *os.File

	if file, err = os.Open(fileName); err != nil {
		return nil, err
	}

	c = &Config{fileName, make(map[string]map[string]string)}

	if err = c.read(bufio.NewReader(file)); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) AddSection(section string) bool {
	if _, ok := c.data[section]; ok {
		return false
	}
	c.data[section] = make(map[string]string)
	return true
}

func (c *Config) AddOption(section string, option string, value string) bool {
	c.AddSection(section) // make sure section exists

	section = strings.ToLower(section)
	option = strings.ToLower(option)

	_, ok := c.data[section][option]
	c.data[section][option] = value

	return !ok
}

func (c *Config) HasSection(section string) bool {
	if _, ok := c.data[strings.ToLower(section)]; ok {
		return true
	}
	return false
}


func (c *Config) GetOptions(section string) (options []string, err error) {
	section = strings.ToLower(section)

	if _, ok := c.data[section]; !ok {
		return nil, errors.New("section not found")
	}

	i := 0
	options = make([]string, len(c.data[section]))
	for s, _ := range c.data[section] {
		options[i] = s
		i++
	}

	return options, nil
}


func (c *Config) HasOption(section string, option string) bool {
	section = strings.ToLower(section)
	option = strings.ToLower(option)

	if _, ok := c.data[section]; !ok {
		return false
	}
	_, ok := c.data[section][option]
	return ok
}

func (c *Config) String(section string, option string) (string, error) {
	if _, ok := c.data[section]; !ok {
		return "", errors.New("setion not fond")
	}
	
	if value, ok := c.data[section][option]; ok {		
		return value, nil
	}
	return "", errors.New("option not found")
}


func (c *Config) Bool(section string, option string) (bool, error) {
	value, err := c.String(section, option)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(value)
}

func (c *Config) Int(section string, option string) (int, error) {
	value, err := c.String(section, option)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

func (c *Config) Float(section string, option string) (float64, error) {
	value, err := c.String(section, option)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(value, 64)
}

func stripComments(l string) string {
	for _, c := range []string{" ;", "\t;", " #", "\t#"} {
		if i := strings.Index(l, c); i != -1 {
			l = l[0:i]
		}
	}
	return l
}

func (c *Config) read(buf *bufio.Reader) (err error) {
	var section, option string

	for {
		l, err := buf.ReadString('\n') // parse line-by-line

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		
		l = strings.TrimSpace(l)

		switch {
		case len(l) == 0, l[0] == '#', l[0] == ';':
			continue
		case len(l) >= 3 && strings.ToLower(l[0:3]) == "rem":
			continue
		case l[0] == '[' && l[len(l)-1] == ']':
			option = "" 
			section = strings.TrimSpace(l[1 : len(l)-1])
			c.AddSection(section)
		default:
			i := strings.IndexAny(l, "=:")
			switch {
			case i > 0:
				
				i := strings.IndexAny(l, "=:")
				option = strings.TrimSpace(l[0:i])
				value := strings.TrimSpace(stripComments(l[i+1:]))
				c.AddOption(section, option, value)
			case section != "" && option != "":
				prev, _ := c.String(section, option)
				value := strings.TrimSpace(stripComments(l))
				c.AddOption(section, option, prev+"\n"+value)
			default:
				return errors.New("could not parse line: " + l)
			}
		}
	}
	return nil
}