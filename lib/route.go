package lib

import "regexp"

/*
RouteConfig Object
*/
type RouteConfig StringMap

/*
RouteCallback Object
*/
type RouteCallback func(*Route) bool

/*
Route Object
*/
type Route struct {
	ParamMixin
	FormatMixin
	PatternMixin
	MatcherMixin
	MiddlewareMixin

	path          string
	pathMatcher   *PathMatcher
	methodMatcher *MethodMatcher
	handler       RequestHandler
	config        RouteConfig
}

func (r *Route) init() {
	r.ParamMixin.initParams()
	r.FormatMixin.initFormats()
	r.PatternMixin.initPatterns()
	r.MatcherMixin.initMatchers()
	r.MiddlewareMixin.initMiddlewares()

	r.params.Copy(ParseParams(r.path))
	pathMatcher := PathMatch(r.path)
	r.pathMatcher = pathMatcher
	r.reportToPathMatcher()
	r.AddMatcher(r.pathMatcher)

	methodMatcher := MethodMatch()
	r.methodMatcher = methodMatcher
	r.AddMatcher(methodMatcher)
}

/*
CopyFormats Method
*/
func (r *Route) CopyFormats(fsm FormatMixin) *Route {
	r.FormatMixin.CopyFormats(fsm)
	r.reportToPathMatcher()
	return r
}

/*
ClearFormats Method
*/
func (r *Route) ClearFormats() *Route {
	r.FormatMixin.ClearFormats()
	r.reportToPathMatcher()
	return r
}

/*
AddFormat Method
*/
func (r *Route) AddFormat(f ...string) *Route {
	r.FormatMixin.AddFormat(f...)
	r.reportToPathMatcher()
	return r
}

/*
CopyPatterns Method
*/
func (r *Route) CopyPatterns(pm PatternMixin) *Route {
	r.PatternMixin.CopyPatterns(pm)
	r.reportToPathMatcher()
	return r
}

/*
ClearPatterns Method
*/
func (r *Route) ClearPatterns() *Route {
	r.PatternMixin.ClearPatterns()
	r.reportToPathMatcher()
	return r
}

/*
Where Method
*/
func (r *Route) Where(name string, pattern string) *Route {
	r.SetPattern(name, pattern)
	r.reportToPathMatcher()
	return r
}

func (r *Route) reportToPathMatcher() {
	if r.Params().Size() > 0 {
		r.pathMatcher.SetRegexp(MakePatternRegexp(r.path, r.Patterns(), *r.Params(), true, *r.Formats()...))
	} else {
		r.pathMatcher.SetCompare()
	}
}

/*
GetName Method
*/
func (r *Route) GetName() string {
	name, found := r.config["name"]
	if found {
		return name
	}
	return ""
}

/*
Name Method
*/
func (r *Route) Name(name string) *Route {
	r.config["name"] = name
	return r
}

/*
Match Method
*/
func (r *Route) Match(ctx *Context) bool {
	return r.Matches(ctx)
}

/*
Handle Method
*/
func (r *Route) Handle(ctx *Context) {
	r.handler(ctx)
}

/*
SetMethod Method
*/
func (r *Route) SetMethod(methods ...string) *Route {
	r.methodMatcher.Set(methods...)
	return r
}

/*
AddMethod Method
*/
func (r *Route) AddMethod(method string) *Route {
	r.methodMatcher.Add(method)
	return r
}

/*
Methods Method
*/
func (r *Route) Methods() []string {
	return r.methodMatcher.Methods()
}

/*
PathRegexp Method
*/
func (r *Route) PathRegexp() *regexp.Regexp {
	return r.pathMatcher.MatcherItem.GetRegexp()
}

/*
GenerateParams Method
*/
func (r *Route) GenerateParams(path string) *Params {
	params := &Params{}
	params.Copy(*r.Params())
	reg := r.PathRegexp()
	if reg != nil && params.Size() > 0 {
		matches := reg.FindStringSubmatch(path)
		for i, name := range reg.SubexpNames() {
			if i != 0 {
				params.Set(name, matches[i])
			}
		}
	}
	return params
}

/*
NewRoute Function
*/
func NewRoute(path string, handler RequestHandler, methods ...string) *Route {
	route := &Route{
		path:    path,
		config:  RouteConfig{},
		handler: handler,
	}
	route.init()
	if len(methods) > 0 {
		route.SetMethod(methods...)
	}
	return route
}
