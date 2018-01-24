package lib

/*
ContextResponder Object
*/
type ContextResponder interface {
	HandleError(*Context, int, string)
}
