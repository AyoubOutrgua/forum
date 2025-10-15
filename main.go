package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type Comment struct {
	Text     string
	Author   string
	PostID   int
	Date     string
	Likes    int
	Dislikes int
}

var Comments []Comment

func main() {
	http.HandleFunc("/home", Home)
	fmt.Println("http://localhost:8080/home")
	http.ListenAndServe(":8080", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		text := r.FormValue("text")
		author := r.FormValue("author")
		idPost := r.FormValue("postId")
		like := r.FormValue("like")
		deslike := r.FormValue("deslike")
		if text == "" || author == "" || idPost == "" {
			http.Error(w, "400 bad request: le commentaire est vide", http.StatusBadRequest)
			return
		}

		now := time.Now()
		date := now.Format("2006-01-02 15:04:05")
		
		postId, err := strconv.Atoi(idPost)
		likes, _ := strconv.Atoi(like)
		deslikes, _ := strconv.Atoi(deslike)
		if err != nil {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}
		details := Comment{
			Text:   text,
			Author: author,
			PostID: postId,
			Date: date,
			Likes: likes,
			Dislikes: deslikes,
		}

		Comments = append(Comments, details)

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return

	case http.MethodGet:
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "500 internal server error: template introuvable", http.StatusInternalServerError)
			return
		}

		errExec := tmpl.Execute(w, map[string]interface{}{"Comments": Comments})

		if errExec != nil {
			http.Error(w, "500 internal server error: erreur template", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)

	}
}

/*
now := time.Now()
formatted := now.Format("2006-01-02 15:04:05")
*/
