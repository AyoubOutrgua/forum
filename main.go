package main

import (
	"net/http"
	"text/template"
)

var Comment []string

func main() {
	http.HandleFunc("/home", Home)
	http.ListenAndServe(":8080", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		input := r.FormValue("comment")
		if input == "" {
			http.Error(w, "401 bad request", http.StatusBadRequest)
			return
		}
		Comment = append(Comment, input)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	errExec := tmpl.Execute(w, Comment)
	if errExec != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
}
