package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"database/sql"
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

func main() {
	// On se connecte à la base SQLite existante
	db, err := sql.Open("sqlite3", "./file.db") // sqlite3 is a driver (motrjim) , sql open 
	if err != nil {
		panic("Erreur lors de l'ouverture de file.db : " + err.Error())
	}
	defer db.Close() // kanssedo l connexion li mabin go u file.db 

	// Route principale
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, db)
	})

	fmt.Println("Serveur lancé sur : http://localhost:8080/home")
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
		rows, err := db.Query("SELECT text, author, post_id, date, likes, dislikes FROM comments") //This is the line that "brings back" everything we have stored in the DB to use in our program.
		if err != nil {
			http.Error(w, "500 internal server error: lecture échouée: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var comments []Comment
		for rows.Next() { // like we have 3 rows if we are in the first one rows.next = true 
			var c Comment
			rows.Scan(&c.Text, &c.Author, &c.PostID, &c.Date, &c.Likes, &c.Dislikes) //copie les valeurs de la row de la DB dans la struct Go.
			comments = append(comments, c)
		}
		// If I change the order in the SELECT, I must also change the order of the variables in Scan, otherwise I will have data in the wrong columns.
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
