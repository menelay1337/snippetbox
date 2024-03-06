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
	router.Handler(http.MethodGet, "/", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.home))))
	router.Handler(http.MethodGet, "/snippet/view/:id", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.snippetView))))
	router.Handler(http.MethodGet, "/snippet/create", app.sessionManager.LoadAndSave(noSurf(app.requireAuthentication(http.HandlerFunc(app.snippetCreate)))))
	router.Handler(http.MethodPost, "/snippet/create", app.sessionManager.LoadAndSave(noSurf(app.requireAuthentication(http.HandlerFunc(app.snippetCreatePost)))))
	// authorization
	router.Handler(http.MethodGet, "/user/signup", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userSignup))))
	router.Handler(http.MethodPost, "/user/signup", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userSignupPost))))
	router.Handler(http.MethodGet, "/user/signin", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userLogin))))
	router.Handler(http.MethodPost, "/user/signin", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userLoginPost))))
	router.Handler(http.MethodPost, "/user/logout", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userLogoutPost))))

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
