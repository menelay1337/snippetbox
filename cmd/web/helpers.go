package main 

import (
	"bytes"
	"net/http"
	"log/slog"
	"fmt"
	"time"
	//"runtime/debug"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page] 
	if !ok {
		err := fmt.Errorf("The template %s doesn't exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData {
		CurrentYear: time.Now().Year(),
	}
}



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
