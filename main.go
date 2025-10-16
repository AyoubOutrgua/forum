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

var Comments []Comment

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		panic(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT,
		author TEXT,
		post_id INTEGER,
		date TEXT,
		likes INTEGER,
		dislikes INTEGER
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	db := initDB()
	defer db.Close()

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, db)
	})

	fmt.Println("http://localhost:8080/home")
	http.ListenAndServe(":8080", nil)
}

func Home(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

		date := time.Now().Format("2006-01-02 15:04:05")

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
		rows, err := db.Query("SELECT text, author, post_id, date, likes, dislikes FROM comments")
		if err != nil {
			http.Error(w, "500 internal server error: lecture échouée", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var comments []Comment
		for rows.Next() {
			var c Comment
			rows.Scan(&c.Text, &c.Author, &c.PostID, &c.Date, &c.Likes, &c.Dislikes)
			comments = append(comments, c)
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "500 internal server error: template introuvable", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, map[string]interface{}{"Comments": comments})
	default:
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

/*
now := time.Now()
formatted := now.Format("2006-01-02 15:04:05")
*/
