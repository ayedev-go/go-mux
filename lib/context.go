package lib

import (
	"net/http"
	"net/url"
	"strings"
)

/*
Context Object
*/
type Context struct {
	Writer     *Writer
	Request    *http.Request
	Path       string
	RequestURI string
	Method     string
	Host       string
	Port       string
	Schema     string
	Matched    bool
	Errors     *Errors
	Params     *Params
	responder  *ContextResponder
	IOWriter   http.ResponseWriter
	finalized  bool
}

func (ctx *Context) init() {
	path := ctx.Request.RequestURI
	if len(path) < 1 {
		path = "/"
	}
	ctx.Path = path
	if i := strings.Index(ctx.Path, "?"); i != -1 {
		ctx.Path = ctx.Path[:i]
	}
	ctx.RequestURI = path
	ctx.Matched = false
	ctx.finalized = false
	ctx.Errors = &Errors{}
	ctx.Params = &Params{}
	ctx.Method = ctx.Request.Method
	ctx.extractHostInfo()
}

func (ctx *Context) extractHostInfo() {
	ctx.Schema = "http"
	sch := ctx.Request.URL.Scheme
	if len(sch) > 0 {
		ctx.Schema = sch
	}
	if ctx.Schema == "http" && ctx.Request.TLS != nil {
		ctx.Schema = "https"
	}

	port := "80"
	host := ctx.Request.Host
	if ctx.Request.URL.IsAbs() {
		host = ctx.Request.URL.Host
	}
	if i := strings.Index(host, ":"); i != -1 {
		port = host[i+1:]
		host = host[:i]
	}
	ctx.Host = host
	ctx.Port = port
}

/*
SetResponder Method
*/
func (ctx *Context) SetResponder(responder ContextResponder) {
	ctx.responder = &responder
}

/*
SetParams Method
*/
func (ctx *Context) SetParams(params *Params) {
	ctx.Params.Copy(*params)
}

/*
MarkFinalized Method
*/
func (ctx *Context) MarkFinalized() {
	ctx.finalized = true
}

/*
Status Method
*/
func (ctx *Context) Status(status int) {
	ctx.Writer.Status(status)
}

/*
Error Method
*/
func (ctx *Context) Error(errorCode int, errMsg string) {
	ctx.Errors.Add(errMsg, errorCode, nil)
	(*(ctx.responder)).HandleError(ctx, errorCode, errMsg)
}

/*
DetailedError Method
*/
func (ctx *Context) DetailedError(errorCode int, errMsg string, obj interface{}) {
	ctx.Errors.Add(errMsg, errorCode, obj)
	(*(ctx.responder)).HandleError(ctx, errorCode, errMsg)
}

/*
Write Method
*/
func (ctx *Context) Write(data []byte) {
	ctx.Writer.Write(data)
}

/*
WriteString Method
*/
func (ctx *Context) WriteString(str string) {
	ctx.Writer.WriteString(str)
}

/*
Start Method
*/
func (ctx *Context) Start() {
	//	do nothing
}

/*
Finalize Method
*/
func (ctx *Context) Finalize() {
	if !ctx.finalized {
		ctx.Writer.PushTo(ctx.IOWriter)
		ctx.finalized = true
	}
}

/*
SetMatched Method
*/
func (ctx *Context) SetMatched(matched bool) {
	ctx.Matched = matched
}

/*
QueryValue Method
*/
func (ctx *Context) QueryValue(key string, fallback ...string) string {
	return ctx.ReadURLValues(ctx.Request.URL.Query(), key, fallback...)
}

/*
FormValue Method
*/
func (ctx *Context) FormValue(key string, fallback ...string) string {
	if ctx.Request.Form == nil {
		ctx.Request.FormValue(key)
	}
	return ctx.ReadURLValues(ctx.Request.Form, key, fallback...)
}

/*
ReadURLValues Method
*/
func (ctx *Context) ReadURLValues(values url.Values, key string, fallback ...string) string {
	if values != nil {
		val := values.Get(key)
		if len(val) > 0 {
			return val
		}
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

/*
NewContext Function
*/
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{Writer: NewWriter(), IOWriter: w, Request: r}
	ctx.init()
	return ctx
}
