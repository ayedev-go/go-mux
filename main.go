package main

import (
	"net/http"

	"./lib"
)

func main() {
	http.HandleFunc("/", serverHandler)
	http.ListenAndServe(":8080", nil)
}

func helloView(ctx *lib.Context) {
	ctx.Write([]byte("reached hello"))
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	ctx := lib.NewContext(w, r)
	ctx.AddPattern("id", "([0-9]+)")

	rg := lib.NewRouter()
	rg.AddRoute(ctx, "user/:id", helloView, nil)

	rg.SetErrorHandler(404, func(ctx *lib.Context) {
		ctx.Write([]byte("404! Page not found"))
	})

	rg.Handle(ctx)
}
