package lib

import (
	"fmt"
	"net/http"
	"regexp"
)

/*
Router Object
*/
type Router struct {
	id            string
	path          string
	matchers      []*RequestMatcher
	methods       map[string][]string
	errorHandlers map[int]RequestHandler
	reg           *regexp.Regexp
}

func (rg *Router) init(ctx *Context) {
	rg.methods = map[string][]string{}
	rg.errorHandlers = map[int]RequestHandler{}

	params := ParseParams(rg.path)
	rg.reg = MakePatternRegexp(ctx, rg.path, params)
}

/*
ID Function
*/
func (rg *Router) ID() string {
	return rg.id
}

/*
SubRouter Function
*/
func (rg *Router) SubRouter(ctx *Context, path string) *Router {
	newPath := path
	if len(rg.path) > 0 {
		newPath = rg.path + "/" + newPath
	}
	sub := &Router{path: newPath, id: rg.id + "-" + path}
	sub.init(ctx)
	rg.AddMatcher(sub)
	return sub
}

/*
SetErrorHandler Function
*/
func (rg *Router) SetErrorHandler(code int, handler RequestHandler) {
	rg.errorHandlers[code] = handler
}

/*
AddRoute Function
*/
func (rg *Router) AddRoute(ctx *Context, path string, handler RequestHandler, config RouteConfig, methods ...string) RequestMatcher {
	if config == nil {
		config = RouteConfig{}
	}
	newPath := path
	if len(rg.path) > 0 {
		newPath = rg.path + "/" + newPath
	}
	r := &Route{
		path:     newPath,
		config:   config,
		handlers: []RequestHandler{handler},
	}
	r.init(ctx)
	rg.AddMatcher(r, methods...)
	return r
}

/*
AddMatcher Function
*/
func (rg *Router) AddMatcher(matcher RequestMatcher, methods ...string) {
	rg.matchers = append(rg.matchers, &matcher)
	rg.methods[matcher.ID()] = methods
}

/*
Match Function
*/
func (rg *Router) Match(ctx *Context) RequestMatcher {
	for _, matcher := range rg.matchers {
		match := (*matcher).Match(ctx)
		if match != nil {
			if rg.matchMethod(ctx, (*matcher).ID()) {
				return match
			}
			rg.HandleError(ctx, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
	return nil
}

/*
Handle Function
*/
func (rg *Router) Handle(ctx *Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	match := rg.Match(ctx)
	if match != nil {
		match.Handle(ctx)
	} else {
		rg.HandleError(ctx, http.StatusNotFound, "Page not found")
	}
}

/*
Regex Function
*/
func (rg *Router) Regex() *regexp.Regexp {
	return rg.reg
}

/*
HandleError Function
*/
func (rg *Router) HandleError(ctx *Context, errorCode int, errorMsg string) {
	errorHandler, found := rg.errorHandlers[errorCode]
	ctx.Status(errorCode)
	if found {
		errorHandler(ctx)
	} else {
		ctx.Write([]byte(fmt.Sprintf("Error #%d => %s", errorCode, errorMsg)))
	}
}

func (rg *Router) matchMethod(ctx *Context, id string) bool {
	matches := false
	validMethods, _ := rg.methods[id]
	if len(validMethods) > 0 {
		for _, m := range validMethods {
			if m == ctx.Method {
				return true
			}
		}
	} else {
		matches = true
	}
	return matches
}

/*
NewRouter Function
*/
func NewRouter(ctx *Context) *Router {
	rg := &Router{id: "default-router"}
	rg.init(ctx)
	return rg
}
