package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Not found handler

	router.NotFound = http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// static files handler
	fileServer := http.FileServer(http.Dir(cfg.staticDir)) 
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// handling routes
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
