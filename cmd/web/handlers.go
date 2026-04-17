package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/kyzercmd/snipio/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/pages/home.tmpl.html",
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/navbar.tmpl.html",
	// }

	// tmpl, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = tmpl.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
}

func (app *application)snippetView(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("id")
	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord){
			http.NotFound(w, r)
		}else {
			app.serverError(w, r, err)
		}
		return
	}
	data := templateData{
		Snippet: snippet,
	}

	files := []string {
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/navbar.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application)snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("FORM TO CREATE SNIPPET"))
}

func (app *application)snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	title := "Flower bed"
	content := "Flowers bright and colorful \n o' so beautiful \n here i lay till the end"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}