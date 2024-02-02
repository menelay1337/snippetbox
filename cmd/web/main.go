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
	logFormat string
}

var cfg config

func main() {
	// assigning flags
	flag.StringVar(&cfg.port,"port", ":3000", "Port of http server.")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.logFormat, "format", "text", "format of logs(json or plain text)")
	// obtaining cmd flags and assigns them to variables'
	flag.Parse()
	// logger initialize
	var logger *slog.Logger
	var loggerOptions = &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}

	if cfg.logFormat == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, loggerOptions))
	} else if cfg.logFormat == "text" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, loggerOptions))
	}
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

	// log Start of the server
	logger.Info("Starting server", slog.Any("port", cfg.port))

	err := http.ListenAndServe(cfg.port, mux)

	logger.Error(err.Error())
	os.Exit(1)

}
