package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kyzercmd/snipio/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	config         config
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP Network Address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Static Assets Path")
	flag.StringVar(&cfg.dsn, "dsn", "web:1234@tcp(localhost:3306)/Snipio?parseTime=true", "MySQL Data Source Name")

	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDb(cfg.dsn)
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

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		config:         cfg,
		logger:         logger,
		templateCache:  templateCache,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server", slog.String("addr", ":4000"))
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func openDb(dsn string) (*sql.DB, error) {
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
