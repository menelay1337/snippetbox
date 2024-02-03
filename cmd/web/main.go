package main

import (
	"net/http"
	"database/sql"
	"log/slog"
	"fmt"
	"flag"
	"os"

	_ "github.com/go-sql-driver/mysql"
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
	dsn		  string
}

var cfg config

// Database functions

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	var dsn string;
	fmt.Print("Enter DSN: ")
	fmt.Scan(&dsn)

	// assigning flags
	flag.StringVar(&cfg.port,"port", ":3000", "Port of http server.")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.logFormat, "format", "text", "format of logs(json or plain text)")
	flag.StringVar(&cfg.dsn, "DSN", "user:pass@/snippetbox?parseTime=true", "data source name for entering database")

	if dsn != "" {
		cfg.dsn = dsn
	}

	// obtaining cmd flags and assigns them to variables'
	flag.Parse()

	// application and logger initialize
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

	app := &application{
		logger: logger,
	}

	// open db connection pool

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// start of the server
	logger.Info("Starting server", slog.Any("port", cfg.port))

	err = http.ListenAndServe(cfg.port, app.routes())

	logger.Error(err.Error())
	os.Exit(1)

}
