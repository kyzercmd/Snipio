package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	files := []string{
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/navbar.tmpl.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}


}

func (app *application)snippetView(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("id")
	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "SNIPPET VIEW PAGE FOR ID: ", id)

}

func (app *application)snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("FORM TO CREATE SNIPPET"))
}

func (app *application)snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save new snippet"))
}