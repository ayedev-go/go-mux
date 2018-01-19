package main

import (
	"fmt"
	"net/http"

	"./lib"
)

func main() {
	http.HandleFunc("/", serverHandler)
	http.ListenAndServe(":8081", nil)
}

func homeView(ctx *lib.Context) {
	ctx.Write([]byte("Homepage"))
}

func profileView(ctx *lib.Context) {
	ctx.Write([]byte("User Profile"))
}

func errorView(ctx *lib.Context) {
	ctx.Error(http.StatusInternalServerError, "Internal Server Error")
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	ctx := lib.NewContext(w, r)
	ctx.AddPattern("id", "([0-9]+)")

	rg := lib.NewRouter(ctx)
	rg.AddRoute(ctx, "/", homeView, nil)
	rg.AddRoute(ctx, "error", errorView, nil)
	rg.SetErrorHandler(404, func(ctx *lib.Context) {
		ctx.Write([]byte("404! Page not found"))
	})

	sub := rg.SubRouter(ctx, "user")
	sub.AddRoute(ctx, ":id", profileView, nil)

	fmt.Println(ctx.Method, ctx.Path)
	ctx.SetResponder(rg)
	rg.Handle(ctx)
}
