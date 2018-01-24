package lib

import (
	"fmt"
	"strings"
)

/*
LoggerMiddleware Function
*/
func LoggerMiddleware(ctx *Context, next PipelineCallback) {
	next()
	fmt.Printf("%s: %s [%v]\n", ctx.Method, ctx.Path, ctx.Matched)
}

/*
PrettyErrorsMiddleware Function
*/
func PrettyErrorsMiddleware(ctx *Context, next PipelineCallback) {
	next()
	if !ctx.Errors.Empty() {
		msgs := []string{}
		for _, err := range *ctx.Errors {
			msgs = append(msgs, fmt.Sprintf("[%v] => %s", err.Code, err.Message))
		}
		fmt.Printf("Errors recovered:\n%s\n\n", strings.Join(msgs, "\n"))
	}
}
