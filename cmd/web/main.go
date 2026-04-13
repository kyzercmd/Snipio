package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr string
	staticDir string
}

type application struct {
	logger *slog.Logger
}

func main() {
	mux := http.NewServeMux()

	var cfg config
	
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP Network Address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Static Assets Path")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	logger.Info("Starting server", slog.String("addr", ":4000"))

	err := http.ListenAndServe(cfg.addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}