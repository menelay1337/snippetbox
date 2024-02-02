package main

import (
	"flag"
	"net/http"
	"log/slog"
	"os"
)

type config struct {
	port	  string
	staticDir string
}

var cfg config

func main() {
	// assigning flags
	flag.StringVar(&cfg.port,"port", ":3000", "Port of http server.")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	// obtaining cmd flags and assigns them to variables'
	flag.Parse()
	// logger insatnce
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// mux router
	mux := http.NewServeMux()
	// file server handler
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// handling static stuff
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// serving URLs
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// flag.Parse() returns pointers, we must derefernce them to obtain their values
	log.Printf("Starting server on port %s.\nPress Ctrl+C to stop the server.", *cfg.port)

	err := http.ListenAndServe(*PORT, mux)
	log.Fatal(err)

}
