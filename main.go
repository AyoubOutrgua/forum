package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Comment struct {
	Text     string
	Author   string
	PostID   int
	Date     string
	Likes    int
	Dislikes int
}

type Post struct {
	ID       int
	Title    string
	Content  string
	Author   string
	Date     string
	Comments []Comment // slice dial comments per post
}

//var posts []Post

func main() {
	// On se connecte à la base SQLite existante
	dbcomment, err1 := sql.Open("sqlite3", "./comment.db") // sqlite3 is a driver (motrjim) , sql open
	if err1 != nil {
		panic("Erreur lors de l'ouverture de file.db : " + err1.Error())
	}
	defer dbcomment.Close() // kanssedo l connexion li mabin go u file.db
	// On se connecte à la base SQLite existante
	dbpost, err2 := sql.Open("sqlite3", "./post.db") // daba ghi post.db
	if err2 != nil {
		panic("Erreur lors de l'ouverture post.db : " + err2.Error())
	}
	defer dbpost.Close()

	// Route principale
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, dbcomment, dbpost)
	})

	fmt.Println("Serveur lancé sur : http://localhost:8080/home")
	http.ListenAndServe(":8080", nil)
}

func Home(w http.ResponseWriter, r *http.Request, db *sql.DB, Pdb *sql.DB) {
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

		postId, err := strconv.Atoi(idPost)
		likes, _ := strconv.Atoi(like)
		dislikes, _ := strconv.Atoi(deslike)
		if err != nil {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		date := time.Now().Format("2006-01-02 15:04")

		// Insérer un commentaire dans la table "comments" de la DB en toute sécurité
		_, err = db.Exec(`
			INSERT INTO comments (text, author, post_id, date, likes, dislikes) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			text, author, postId, date, likes, dislikes,
		)
		if err != nil {
			http.Error(w, "500 internal server error: insertion échouée", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return

	case http.MethodGet:
		// 1️⃣ Jibou comments men DB
		rows, err := db.Query("SELECT text, author, post_id, date, likes, dislikes FROM comments")
		if err != nil {
			http.Error(w, "500 internal server error: lecture comments échouée: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var comments []Comment
		for rows.Next() {
			var c Comment
			rows.Scan(&c.Text, &c.Author, &c.PostID, &c.Date, &c.Likes, &c.Dislikes)
			comments = append(comments, c)
		}

		// 2️⃣ Jibou posts men DB
		postRows, err := Pdb.Query("SELECT id, title, content, author, date FROM posts")
		if err != nil {
			http.Error(w, "500 internal server error: lecture posts échouée: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer postRows.Close()

		var posts []Post
		for postRows.Next() {
			var p Post
			postRows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Date)

			// 3️⃣ Attach comments li kaynin f had post
			for _, c := range comments {
				if c.PostID == p.ID {
					p.Comments = append(p.Comments, c)
				}
			}

			posts = append(posts, p)
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "500 internal server error: template introuvable", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, map[string]interface{}{"Posts": posts})

	default:
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}
