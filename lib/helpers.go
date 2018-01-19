package lib

import (
	"regexp"
	"strings"
)

/*
ParseParams Function
*/
func ParseParams(path string) RouteParams {
	params := RouteParams{}
	reg, rgError := regexp.Compile(":([a-z0-9_]+)")
	if rgError == nil {
		matches := reg.FindAllString(path, -1)
		for _, match := range matches {
			params.Add(match[1:], match)
		}
	}
	return params
}

/*
MakePatternRegexp Function
*/
func MakePatternRegexp(ctx *Context, path string, params RouteParams, formats ...string) *regexp.Regexp {
	rgText := path
	for _, param := range params {
		rgText = strings.Replace(rgText, param.dummy, ctx.Pattern(param.key, "([a-zA-Z0-9-_]+)"), -1)
	}
	if len(path) > 0 && path != "/" && len(formats) > 0 {
		rgText += "(.(" + strings.Join(formats, "|") + "))?"
	} else if path != "/" {
		rgText += "(/)?"
	}
	rg2, rgError2 := regexp.Compile("^" + rgText + "$")
	if rgError2 == nil {
		return rg2
	}
	return nil
}
