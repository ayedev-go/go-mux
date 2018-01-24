package lib

/*
Error Object
*/
type Error struct {
	Code    int
	Message string
	Object  interface{}
}

/*
Errors Object
*/
type Errors []Error

/*
Add Method
*/
func (es *Errors) Add(msg string, code int, obj interface{}) *Error {
	err := &Error{Message: msg, Code: code, Object: obj}
	*es = append(*es, *err)
	return err
}

/*
Size Method
*/
func (es *Errors) Size() int {
	return len(*es)
}

/*
Empty Method
*/
func (es *Errors) Empty() bool {
	return es.Size() < 1
}

/*
ErrorHandlers Object
*/
type ErrorHandlers map[int]RequestHandler

/*
Add Method
*/
func (ehs *ErrorHandlers) Add(code int, h RequestHandler) {
	(*ehs)[code] = h
}

/*
Copy Method
*/
func (ehs *ErrorHandlers) Copy(ehs2 ErrorHandlers) {
	for code, h := range ehs2 {
		ehs.Add(code, h)
	}
}

/*
ErrorHandlerMixin Object
*/
type ErrorHandlerMixin struct {
	errorHandlers *ErrorHandlers
}

func (ehm *ErrorHandlerMixin) initErrorHandlers() {
	ehm.errorHandlers = &ErrorHandlers{}
}

/*
SetErrorHandler Method
*/
func (ehm *ErrorHandlerMixin) SetErrorHandler(code int, h RequestHandler) *ErrorHandlerMixin {
	ehm.errorHandlers.Add(code, h)
	return ehm
}

/*
GetErrorHandler Method
*/
func (ehm *ErrorHandlerMixin) GetErrorHandler(code int) RequestHandler {
	handler, found := (*ehm.errorHandlers)[code]
	if found {
		return handler
	}
	return nil
}

/*
ErrorHandlers Method
*/
func (ehm *ErrorHandlerMixin) ErrorHandlers() *ErrorHandlers {
	return ehm.errorHandlers
}

/*
CopyErrorHandlers Method
*/
func (ehm *ErrorHandlerMixin) CopyErrorHandlers(ehm2 ErrorHandlerMixin) {
	ehm.errorHandlers.Copy(*ehm2.ErrorHandlers())
}
