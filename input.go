package dogo

import (
	"strings"
	"strconv"
	"regexp"
)

type Input struct{
	Context *Context
} 

var (
	acceptsHTMLRegex = regexp.MustCompile(`(text/html|application/xhtml\+xml)(?:,|$)`)
	acceptsXMLRegex  = regexp.MustCompile(`(application/xml|text/xml)(?:,|$)`)
	acceptsJSONRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
)

// Protocol returns request protocol name, such as HTTP/1.1 .
func (input *Input) Protocol() string {
	return input.Context.Request.Proto
}

// URI returns full request url with query string, fragment.
func (input *Input) URI() string {
	return input.Context.Request.RequestURI
}

// URL returns request url path (without query string, fragment).
func (input *Input) URL() string {
	return input.Context.Request.URL.Path
}

// Site returns base site url as scheme://domain type.
func (input *Input) Site() string {
	return input.Scheme() + "://" + input.Domain()
}

// Scheme returns request scheme as "http" or "https".
func (input *Input) Scheme() string {
	if scheme := input.Header("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if input.Context.Request.URL.Scheme != "" {
		return input.Context.Request.URL.Scheme
	}
	if input.Context.Request.TLS == nil {
		return "http"
	}
	return "https"
}

func (input *Input) Domain() string {
	return input.Host()
}
func (input *Input) Host() string {
	if input.Context.Request.Host != "" {
		hostParts := strings.Split(input.Context.Request.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return input.Context.Request.Host
	}
	return "localhost"
}

func (input *Input) Method() string {
	return input.Context.Request.Method
}

func (input *Input) Header(key string) string {
	return input.Context.Request.Header.Get(key)
}

func (input *Input) Is(method string) bool {
	return input.Method() == method
}

func (input *Input) IsGet() bool {
	return input.Is("GET")
}

func (input *Input) IsPost() bool {
	return input.Is("POST")
}

func (input *Input) IsAjax() bool {
	return input.Header("X-Requested-With") == "XMLHttpRequest"
}

func (input *Input) IsHttps() bool {
	return input.Scheme() == "https"
}

func (input *Input) IsWebsocket() bool {
	return input.Header("Upgrade") == "websocket"
}

// IsUpload returns boolean of whether file uploads in this request or not..
func (input *Input) IsUpload() bool {
	return strings.Contains(input.Header("Content-Type"), "multipart/form-data")
}

// AcceptsHTML Checks if request accepts html response
func (input *Input) AcceptsHTML() bool {
	return acceptsHTMLRegex.MatchString(input.Header("Accept"))
}

// AcceptsXML Checks if request accepts xml response
func (input *Input) AcceptsXML() bool {
	return acceptsXMLRegex.MatchString(input.Header("Accept"))
}

// AcceptsJSON Checks if request accepts json response
func (input *Input) AcceptsJSON() bool {
	return acceptsJSONRegex.MatchString(input.Header("Accept"))
}

// IP returns request client ip.
// if in proxy, return first proxy id.
// if error, return 127.0.0.1.
func (input *Input) IP() string {
	ips := input.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		return rip[0]
	}
	ip := strings.Split(input.Context.Request.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (input *Input) Proxy() []string {
	if ips := input.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

// Referer returns http referer header.
func (input *Input) Referer() string {
	return input.Header("Referer")
}

// Refer returns http referer header.
func (input *Input) Refer() string {
	return input.Referer()
}

// SubDomains returns sub domain string.
// if aa.bb.domain.com, returns aa.bb .
func (input *Input) SubDomains() string {
	parts := strings.Split(input.Host(), ".")
	if len(parts) >= 3 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}

func (input *Input) Cookie(key string) string {
	ck, err := input.Context.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

// Port returns request client port.
// when error or empty, return 80.
func (input *Input) Port() int {
	parts := strings.Split(input.Context.Request.Host, ":")
	if len(parts) == 2 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	return 80
}