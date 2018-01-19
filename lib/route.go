package lib

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
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

	r.params = RouteParams{}
	rg, rgError := regexp.Compile(":([a-z0-9_]+)")
	if rgError == nil {
		matches := rg.FindAllString(r.path, -1)
		for _, match := range matches {
			r.params.Add(match[1:], match)
		}
	}

	rgText := r.path
	for _, param := range r.params {
		rgText = strings.Replace(rgText, param.dummy, ctx.Pattern(param.key, "([a-zA-Z0-9-_]+)"), -1)
	}
	rgText += "(.(json|xml))?"
	rg2, rgError2 := regexp.Compile("^" + rgText + "$")
	if rgError2 == nil {
		r.reg = rg2
	}
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
		fmt.Println(param)
		return param.Value()
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}
