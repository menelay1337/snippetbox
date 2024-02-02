package main 

import (
	"net/http"
	"log/slog"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method		 = r.Method
		uri			 = r.URL.RequestURI()
		// optional non readable stack trace
		//trace		 = string(debug.Stack())
		statusCode   = http.StatusInternalServerError
	)

	app.logger.Error(err.Error(), slog.Any("method", method), slog.Any("uri", uri))
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
