package main

import (
	// standard library package
	"net/http"
	"database/sql"
	"log/slog"
	"html/template"
	"fmt"
	"flag"
	"os"
	"time"

	// third party packages
	_ "github.com/go-sql-driver/mysql"
	"github.com/menelay1337/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

// application structure for dependency injection
type application struct {
	logger		   *slog.Logger
	snippets	   *models.SnippetModel
	users		   *models.UserModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
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
	fmt.Print("DSN format: user:pass@options/dbname?var_pairs\nEnter DSN: ")
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

	
	// open db connection pool

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	

	// Session store initialization
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{ DB : db },
		users: &models.UserModel{ DB : db },
		templateCache: templateCache,
		sessionManager: sessionManager,
	}
	// start of the server
	logger.Info("Starting server", slog.Any("port", cfg.port))

	err = http.ListenAndServe(cfg.port, app.routes())

	logger.Error(err.Error())
	os.Exit(1)

}
