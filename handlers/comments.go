package handlers

import (
	"net/http"
	"strconv"
	"time"

	"forum/database"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var userID int
	err = database.DataBase.QueryRow("SELECT id FROM users WHERE session = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	postIDStr := r.FormValue("postID")
	commentText := r.FormValue("comment")

	if postIDStr == "" || commentText == "" {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var postExists int
	err = database.DataBase.QueryRow("SELECT COUNT(*) FROM posts WHERE id = ?", postID).Scan(&postExists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if postExists == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	creationDate := time.Now().Format("2006-01-02 15:04:05")
	_, err = database.DataBase.Exec("INSERT INTO comments (comment, postId, userId, creationDate) VALUES(?, ?, ?, ?)",
		commentText, postID, userID, creationDate)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}