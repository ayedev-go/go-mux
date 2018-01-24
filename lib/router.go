package lib

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

/*
Router Object
*/
type Router struct {
	ParamMixin
	FormatMixin
	PatternMixin
	MatcherMixin
	HandlerMixin
	MiddlewareMixin
	ErrorHandlerMixin

	prefix        string
	pathMatcher   *PathMatcher
	methodMatcher *MethodMatcher
}

func (rg *Router) init() {
	rg.ParamMixin.initParams()
	rg.FormatMixin.initFormats()
	rg.PatternMixin.initPatterns()
	rg.MatcherMixin.initMatchers()
	rg.HandlerMixin.initHandlers()
	rg.MiddlewareMixin.initMiddlewares()
	rg.ErrorHandlerMixin.initErrorHandlers()

	rg.params.Copy(ParseParams(rg.prefix))
	pathMatcher := PathMatch(rg.prefix)
	rg.pathMatcher = pathMatcher
	rg.reportToPathMatcher()
	rg.AddMatcher(rg.pathMatcher)

	methodMatcher := MethodMatch()
	rg.methodMatcher = methodMatcher
	rg.AddMatcher(methodMatcher)
}

/*
CopyPatterns Method
*/
func (rg *Router) CopyPatterns(pm PatternMixin) *Router {
	rg.PatternMixin.CopyPatterns(pm)
	rg.reportToPathMatcher()
	return rg
}

/*
ClearPatterns Method
*/
func (rg *Router) ClearPatterns() *Router {
	rg.PatternMixin.ClearPatterns()
	rg.reportToPathMatcher()
	return rg
}

/*
Where Method
*/
func (rg *Router) Where(name string, pattern string) *Router {
	rg.SetPattern(name, pattern)
	rg.reportToPathMatcher()
	return rg
}

func (rg *Router) reportToPathMatcher() {
	rg.pathMatcher.SetRegexp(MakePatternRegexp(rg.prefix, rg.Patterns(), *rg.Params(), false))
}

/*
CopyFormats Method
*/
func (rg *Router) CopyFormats(fsm FormatMixin) *Router {
	rg.FormatMixin.CopyFormats(fsm)
	rg.reportToPathMatcher()
	return rg
}

/*
ClearFormats Method
*/
func (rg *Router) ClearFormats() *Router {
	rg.FormatMixin.ClearFormats()
	rg.reportToPathMatcher()
	return rg
}

/*
AddFormat Method
*/
func (rg *Router) AddFormat(f ...string) *Router {
	rg.FormatMixin.AddFormat(f...)
	rg.reportToPathMatcher()
	return rg
}

/*
MakePrefix Method
*/
func (rg *Router) MakePrefix(appendPath string) string {
	newPrefix := appendPath
	if len(rg.prefix) > 0 {
		newPrefix = rg.prefix + "/" + newPrefix
	}
	return newPrefix
}

/*
MakePrefixWithStart Method
*/
func (rg *Router) MakePrefixWithStart(appendPath string) string {
	path := rg.MakePrefix(appendPath)
	if len(path) > 0 && path[0:1] != "/" {
		path = "/" + path
	}
	return path
}

/*
SetMethod Method
*/
func (rg *Router) SetMethod(methods ...string) *Router {
	rg.methodMatcher.Set(methods...)
	return rg
}

/*
AddMethod Method
*/
func (rg *Router) AddMethod(method string) *Router {
	rg.methodMatcher.Add(method)
	return rg
}

/*
Methods Method
*/
func (rg *Router) Methods() []string {
	return rg.methodMatcher.Methods()
}

/*
SubRouter Method
*/
func (rg *Router) SubRouter(prefix string) *Router {
	sub := RouterWithPrefix(rg.MakePrefix(prefix), rg)
	rg.AddHandler(sub)
	return sub
}

/*
AddRoute Method
*/
func (rg *Router) AddRoute(path string, handler RequestHandler, methods ...string) *Route {
	r := NewRoute(rg.MakePrefix(path), handler, methods...)
	r.CopyFormats(rg.FormatMixin)
	r.CopyPatterns(rg.PatternMixin)
	r.CopyMiddlewares(rg.MiddlewareMixin)
	rg.AddHandler(r)
	return r
}

/*
ANY Method
*/
func (rg *Router) ANY(path string, handler RequestHandler) *Route {
	return rg.AddRoute(path, handler)
}

/*
GET Method
*/
func (rg *Router) GET(path string, handler RequestHandler) *Route {
	return rg.AddRoute(path, handler, "GET", "HEAD")
}

/*
POST Method
*/
func (rg *Router) POST(path string, handler RequestHandler) *Route {
	return rg.AddRoute(path, handler, "POST")
}

/*
PUT Method
*/
func (rg *Router) PUT(path string, handler RequestHandler) *Route {
	return rg.AddRoute(path, handler, "PUT")
}

/*
PATCH Method
*/
func (rg *Router) PATCH(path string, handler RequestHandler) *Route {
	return rg.AddRoute(path, handler, "PATCH")
}

/*
DELETE Method
*/
func (rg *Router) DELETE(path string, handler RequestHandler) *Route {
	return rg.AddRoute(path, handler, "DELETE")
}

/*
HEAD Method
*/
func (rg *Router) HEAD(path string, handler RequestHandler) *Route {
	return rg.AddRoute(path, handler, "HEAD")
}

/*
OPTIONS Method
*/
func (rg *Router) OPTIONS(path string, handler RequestHandler) *Route {
	return rg.AddRoute(path, handler, "OPTIONS")
}

/*
ServePath Method
*/
func (rg *Router) ServePath(path string, folder string) *Route {
	return rg.ServeDir(path, http.Dir(folder))
}

/*
ServeDir Method
*/
func (rg *Router) ServeDir(path string, dir http.Dir) *Route {
	fs := http.FileServer(dir)
	route := rg.GET(path+"/:resource_path", rg.ServeHandler(http.StripPrefix(rg.MakePrefixWithStart(path), fs))).Where("resource_path", "(.*)")
	route.ClearFormats()
	return route
}

/*
ServeFile Method
*/
func (rg *Router) ServeFile(path string, file string) *Route {
	route := rg.GET(path, func(ctx *Context) {
		http.ServeFile(ctx.IOWriter, ctx.Request, file)
		ctx.MarkFinalized()
	})
	route.ClearFormats()
	return route
}

/*
ServeHandler Method
*/
func (rg *Router) ServeHandler(handler http.Handler) RequestHandler {
	return func(ctx *Context) {
		handler.ServeHTTP(ctx.IOWriter, ctx.Request)
		ctx.MarkFinalized()
	}
}

/*
Match Method
*/
func (rg *Router) Match(ctx *Context) bool {
	return rg.Matches(ctx)
}

/*
URL Method
*/
func (rg *Router) URL(name string, data StringMap) string {
	link := name
	var route *Route
	queries := url.Values{}
	rg.WalkRoutes(func(r *Route) bool {
		if r.GetName() == name {
			route = r
			return false
		}
		return true
	})
	if route != nil {
		link = route.path
		for key, val := range data {
			if route.Params().Has(key) {
				link = strings.Replace(link, ":"+key, val, -1)
			} else {
				queries.Set(key, val)
			}
		}
		if len(queries) > 0 {
			link += "?" + queries.Encode()
		}
	}
	return link
}

/*
Debug Method
*/
func (rg *Router) Debug() []string {
	lines := []string{}
	rg.WalkRoutes(func(r *Route) bool {
		name := r.GetName()
		if len(name) < 1 {
			name = "-noname-"
		}
		methods := r.Methods()
		if len(methods) < 1 {
			methods = append(methods, "ANY")
		}
		lines = append(lines, fmt.Sprintf("[%s] %s (%s)", strings.Join(methods, "|"), r.path, name))
		return true
	})
	lines = append(lines, "")
	lines = append(lines, "=======================================")
	lines = append(lines, "")
	return lines
}

/*
WalkRoutes Method
*/
func (rg *Router) WalkRoutes(cb RouteCallback) bool {
	for _, handler := range *rg.Handlers() {
		switch hn := handler.(type) {
		case *Router:
			res := hn.WalkRoutes(cb)
			if res == false {
				return false
			}
		case *Route:
			res := cb(hn)
			if res == false {
				return false
			}
		}
	}
	return true
}

/*
Handle Method
*/
func (rg *Router) Handle(ctx *Context) {
	defer func() {
		if r := recover(); r != nil {
			ctx.Error(http.StatusInternalServerError, fmt.Sprintf("%v", r))
			ctx.Finalize()
		}
	}()

	var route *Route
	activeRouter := rg
	rgSub, route := rg.FindRoute(ctx)
	pipeline := NewPipeline(ctx)
	if rgSub != nil {
		activeRouter = rgSub
		ctx.SetResponder(activeRouter)
	} else {
		for _, handler := range *rg.Handlers() {
			switch hn := handler.(type) {
			case *Router:
				if hn.Match(ctx) {
					ctx.SetResponder(activeRouter)
					activeRouter = hn
					break
				}
			}
		}
	}

	if route == nil {
		pipeline.Copy(activeRouter)
		pipeline.Add(activeRouter.MakeErrorHandler(http.StatusNotFound, "Page not found"))
		ctx.SetMatched(false)
	} else {
		pipeline.Copy(route)
		pipeline.Add(activeRouter.wrapRoute(route))
		ctx.SetMatched(true)
		ctx.SetParams(route.GenerateParams(ctx.Path))
	}

	pipeline.Start(func() {
		ctx.Start()
	}, func() {
		ctx.Finalize()
	})
}

/*
ServerHandler Method
*/
func (rg *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	ctx.SetResponder(rg)
	rg.Handle(ctx)
}

/*
MakeErrorHandler Method
*/
func (rg *Router) MakeErrorHandler(errorCode int, errorMsg string) Middleware {
	return func(ctx *Context, next PipelineCallback) {
		rg.HandleError(ctx, errorCode, errorMsg)
		next()
	}
}

/*
HandleError Method
*/
func (rg *Router) HandleError(ctx *Context, errorCode int, errorMsg string) {
	errorHandler := rg.GetErrorHandler(errorCode)
	ctx.Status(errorCode)
	if errorHandler != nil {
		errorHandler(ctx)
	} else {
		ctx.Write([]byte(fmt.Sprintf("Error #%d => %s", errorCode, errorMsg)))
	}
}

func (rg *Router) wrapRoute(r *Route) Middleware {
	return func(ctx *Context, next PipelineCallback) {
		r.Handle(ctx)
		next()
	}
}

/*
FindRoute Method
*/
func (rg *Router) FindRoute(ctx *Context) (*Router, *Route) {
	for _, handler := range *rg.Handlers() {
		switch hn := handler.(type) {
		case *Router:
			rgSub, rt := hn.FindRoute(ctx)
			if rt != nil {
				return rgSub, rt
			}
		case *Route:
			if rg.Match(ctx) && hn.Match(ctx) {
				return rg, hn
			}
		}
	}
	return nil, nil
}

/*
PlainRouter Function
*/
func PlainRouter() *Router {
	return RouterWithPrefix("", nil)
}

/*
NewRouter Function
*/
func NewRouter() *Router {
	router := PlainRouter()
	router.AddMiddleware(PrettyErrorsMiddleware, LoggerMiddleware)
	return router
}

/*
RouterWithPrefix Function
*/
func RouterWithPrefix(prefix string, parent *Router) *Router {
	rg := &Router{prefix: prefix}
	rg.init()
	if parent != nil {
		rg.CopyParams(*parent.Params())
		rg.CopyFormats(parent.FormatMixin)
		rg.CopyPatterns(parent.PatternMixin)
		rg.CopyMatchers(parent.MatcherMixin)
		rg.CopyMiddlewares(parent.MiddlewareMixin)
		rg.CopyErrorHandlers(parent.ErrorHandlerMixin)
	}
	return rg
}
