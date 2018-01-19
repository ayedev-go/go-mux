package lib

import "regexp"

/*
RequestHandler Object
*/
type RequestHandler func(*Context)

/*
RequestHandlers Object
*/
type RequestHandlers []RequestHandler

/*
RequestMatcher Object
*/
type RequestMatcher interface {
	ID() string
	Regex() *regexp.Regexp
	Match(*Context) RequestMatcher
	Handle(*Context)
}

/*
Responder Object
*/
type Responder interface {
	HandleError(*Context, int, string)
}
