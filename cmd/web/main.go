package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

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
	config        config
	logger        *slog.Logger
	templateCache map[string]*template.Template
	snippets      *models.SnippetModel
	formDecoder   *form.Decoder
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP Network Address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Static Assets Path")
	flag.StringVar(&cfg.dsn, "dsn", "web:1234@/Snipio?parseTime=true", "MySQL Data Source Name")

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

	app := &application{
		config:        cfg,
		logger:        logger,
		templateCache: templateCache,
		snippets:      &models.SnippetModel{DB: db},
		formDecoder:   formDecoder,
	}

	logger.Info("Starting server", slog.String("addr", ":4000"))

	err = http.ListenAndServe(cfg.addr, app.routes())
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
