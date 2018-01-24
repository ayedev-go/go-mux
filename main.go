package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"./lib"
)

func homeView(ctx *lib.Context) {
	ctx.WriteString("Homepage")
}

func profileView(ctx *lib.Context) {
	ctx.WriteString("User Profile")
}

func errorView(ctx *lib.Context) {
	ctx.Error(http.StatusInternalServerError, "Internal Server Error")
}

func main() {
	router := lib.NewRouter()
	router.Where("id", "([0-9]+)").AddFormat("json", "xml")

	router.GET("/", homeView).Name("home")
	router.GET("error", errorView).Name("error")

	subRouter := router.SubRouter("user")
	subRouter.GET(":id", profileView).Name("user_profile")

	fmt.Println(strings.Join(router.Debug(), "\n"))
	mux := http.NewServeMux()
	mux.Handle("/", router)
	err := http.ListenAndServe(":8080", mux)
	//err := http.ListenAndServeTLS(":10443", "server.crt", "server.key", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
