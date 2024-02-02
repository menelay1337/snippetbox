package main

import (
	"flag"
	"net/http"
	"log/slog"
	"os"
)

// application structure for dependency injection
type application struct {
	logger *slog.Logger
}

// config structure
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
	// application initialize
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
	// application initialize 
	app := &application{
		logger: logger,
	}
	// mux router
	mux := http.NewServeMux()
	// file server handler
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// handling static stuff
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// serving URLs
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// log Start of the server
	logger.Info("Starting server", slog.Any("port", cfg.port))

	err := http.ListenAndServe(cfg.port, mux)

	logger.Error(err.Error())
	os.Exit(1)

}
