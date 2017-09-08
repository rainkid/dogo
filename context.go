package dogo

import (
	"net/http"
)

type Context struct { 
	Input 		   *Input
	Output 		   *Output
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func (ctx *Context) Init(rw http.ResponseWriter, r *http.Request) {
	ctx.ResponseWriter = rw
	ctx.Request = r
	ctx.Input = &Input{Context:ctx}
	ctx.Output = &Output{Context:ctx, Started:false}
}

func (ctx *Context) Redirect(status int, localurl string) {
	http.Redirect(ctx.ResponseWriter, ctx.Request, localurl, status)
}