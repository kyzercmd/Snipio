package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kyzercmd/snipio/internal/models"
)

type config struct {
	addr string
	staticDir string
	dsn string
}

type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	var cfg config
	
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP Network Address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Static Assets Path")
	flag.StringVar(&cfg.dsn,"dsn", "web:1234@/Snipio?parseTime=true", "MySQL Data Source Name")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDb(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("Starting server", slog.String("addr", ":4000"))

	err = http.ListenAndServe(cfg.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDb (dsn string) (*sql.DB, error){
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