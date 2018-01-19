package lib

/*
Router Object
*/
type Router struct {
	id            string
	prefix        string
	matchers      []*RequestMatcher
	methods       map[string][]string
	errorHandlers map[int]RequestHandler
}

func (rg *Router) init() {
	rg.methods = map[string][]string{}
	rg.errorHandlers = map[int]RequestHandler{}
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
func (rg *Router) SubRouter(prefix string) *Router {
	sub := &Router{prefix: prefix, id: rg.id + "-" + prefix}
	sub.init()
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
func (rg *Router) AddRoute(ctx *Context, p string, hs RequestHandler, config RouteConfig, methods ...string) RequestMatcher {
	if config == nil {
		config = RouteConfig{}
	}
	r := &Route{
		path:     p,
		config:   config,
		handlers: []RequestHandler{hs},
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
		if rg.matchMethod(ctx, (*matcher).ID()) {
			match := (*matcher).Match(ctx)
			if match != nil {
				return match
			}
		}
	}
	return nil
}

/*
Handle Function
*/
func (rg *Router) Handle(ctx *Context) {
	match := rg.Match(ctx)
	if match != nil {
		match.Handle(ctx)
	} else {
		errorHandler, found := rg.errorHandlers[404]
		if found {
			errorHandler(ctx)
		}
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
func NewRouter() *Router {
	rg := &Router{id: "default-router"}
	rg.init()
	return rg
}
