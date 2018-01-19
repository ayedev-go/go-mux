package lib

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
	Match(*Context) RequestMatcher
	Handle(*Context)
}
