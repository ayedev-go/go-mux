package lib

import "net/http"

/*
Context Object
*/
type Context struct {
	writer   http.ResponseWriter
	request  *http.Request
	Path     string
	Method   string
	Patterns map[string]string
}

/*
init Function
*/
func (ctx *Context) init() {
	path := ctx.request.RequestURI[1:]
	if len(path) < 1 {
		path = "/"
	}
	ctx.Path = path
	ctx.Method = ctx.request.Method
	ctx.Patterns = map[string]string{}
}

/*
AddPattern Function
*/
func (ctx *Context) AddPattern(name string, reg string) {
	ctx.Patterns[name] = reg
}

/*
Pattern Function
*/
func (ctx *Context) Pattern(name string, fallback string) string {
	pattern, pFound := ctx.Patterns[name]
	if pFound {
		return pattern
	}
	return fallback
}

/*
Write Function
*/
func (ctx *Context) Write(data []byte) {
	ctx.writer.Write(data)
}

/*
NewContext Function
*/
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{writer: w, request: r}
	ctx.init()
	return ctx
}
