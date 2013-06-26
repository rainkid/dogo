package dogo

import (
	"net/http"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{ResponseWriter: w, Request: r}
}

func (ctxt *Context) GetResponse() http.ResponseWriter {
	return ctxt.ResponseWriter
}

func (ctxt *Context) GetRequest() *http.Request {
	return ctxt.Request
}

func (ctxt *Context) SetHeader(hdr, value string) {
	ctxt.Request.Header.Set(hdr, value)
	return
}

func (ctxt *Context) GetHeader(hdr string) string {
	return ctxt.Request.Header.Get(hdr)
}
