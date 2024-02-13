package main

import "net/http"

func (app *application) routes() http.Handler {
	// standard router
	mux := http.NewServeMux()

	// static files handler
	fileServer := http.FileServer(http.Dir(cfg.staticDir)) 
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// handling routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
