package lib

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
)

/*
RouteConfig Object
*/
type RouteConfig map[string]string

/*
Route Object
*/
type Route struct {
	id       string
	path     string
	config   RouteConfig
	handlers RequestHandlers
	params   RouteParams
	reg      *regexp.Regexp
}

func (r *Route) init(ctx *Context) {
	rng := rand.Reader
	num := int64(19891991)
	bigInt := big.NewInt(num)
	randomNumber, err := rand.Int(rng, bigInt)
	if err == nil {
		r.id = fmt.Sprintf("%s", randomNumber)
	}

	r.params = ParseParams(r.path)
	r.reg = MakePatternRegexp(ctx, r.path, r.params, "json", "xml")
}

/*
ID Function
*/
func (r *Route) ID() string {
	return r.id
}

/*
Match Function
*/
func (r *Route) Match(ctx *Context) RequestMatcher {
	matched := false
	if len(r.params) > 0 {
		if r.reg != nil {
			matched = r.reg.MatchString(ctx.Path)
		}
	} else {
		matched = ctx.Path == r.path
	}
	if matched {
		return r
	}
	return nil
}

/*
Handle Function
*/
func (r *Route) Handle(ctx *Context) {
	for _, handler := range r.handlers {
		handler(ctx)
	}
}

/*
Regex Function
*/
func (r *Route) Regex() *regexp.Regexp {
	return r.reg
}

/*
HasParam Function
*/
func (r *Route) HasParam(key string) bool {
	_, pError := r.params.Param(key)
	return pError == nil
}

/*
ParamValue Function
*/
func (r *Route) ParamValue(key string, fallback ...string) string {
	param := r.params.FindParam(key)
	if param != nil {
		return param.Value()
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}
