package lib

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

/*
ParseParams Function
*/
func ParseParams(path string) Params {
	params := Params{}
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
func MakePatternRegexp(path string, patterns *Patterns, params Params, addEnd bool, formats ...string) *regexp.Regexp {
	rgText := path
	if len(rgText) > 0 && rgText[0:1] != "/" {
		rgText = "/" + rgText
	}
	if len(rgText) == 0 {
		rgText = "/(.*)"
	}
	for _, param := range params {
		rgText = strings.Replace(rgText, param.placeholder, "(?P<"+param.key+">"+patterns.Get(param.key, "([a-zA-Z0-9-_]+)")+")", -1)
	}
	if len(path) > 0 && path != "/" && len(formats) > 0 {
		rgText += "(.(" + strings.Join(formats, "|") + "))?"
	} else if path != "/" {
		rgText += "(/)?"
	}
	if addEnd {
		rgText += "$"
	}
	rg2, rgError2 := regexp.Compile("^" + rgText)
	if rgError2 == nil {
		return rg2
	}
	return nil
}

/*
RandomString Function
*/
func RandomString() string {
	rng := rand.Reader
	num := int64(19891991)
	bigInt := big.NewInt(num)
	randomNumber, err := rand.Int(rng, bigInt)
	if err == nil {
		return fmt.Sprintf("%s", randomNumber)
	}
	return ""
}

/*
InStringSlice Function
*/
func InStringSlice(list []string, val string) bool {
	found := false
	for _, item := range list {
		if item == val {
			found = true
			break
		}
	}
	return found
}
