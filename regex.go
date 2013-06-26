package dogo

import (
	"regexp"
	"strings"
)

type RegexRoute struct {
	regex   *regexp.Regexp
	params  map[int]string
	forward string
}

// new regex route
func NewRegexRoute(path string, forward string) *RegexRoute {
	//split the url into sections
	parts := strings.Split(path, "/")

	if len(strings.Split(forward, "/")) != 4 {
		return nil
	}

	//find params that start with ":"
	//replace with regular expressions
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’ 
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part
			parts[i] = expr
			j++
		}
	}

	//recreate the url pattern, with parameters replaced
	//by regular expressions. then compile the regex
	pattern := strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		return nil
	}
	return &RegexRoute{regex: regex, params: params, forward: forward}
}
