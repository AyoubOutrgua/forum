package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"forum/database"
	"forum/helpers"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookieValue := helpers.GetCookieValue(w, r)
	if cookieValue == "" {
		return
	}

	userID, errSelect := database.SelectUserID("SELECT id FROM users WHERE session = ?", cookieValue)
	if errSelect == sql.ErrNoRows {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	} else if errSelect != nil {
		helpers.Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
		return
	}
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	errParse := r.ParseForm()
	if errParse != nil {
		helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	commentText := r.FormValue("comment")
	postIDStr := r.FormValue("postId")

	if commentText == "" || postIDStr == "" {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if len(commentText) > 200 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO comments (comment, postId, userId, creationDate) 
              VALUES (?, ?, ?, ?)`

	creationDate := time.Now().Format("2006-01-02 15:04:05")
	errExec := database.ExecuteData(query, commentText, postID, userID, creationDate)
	if errExec != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
