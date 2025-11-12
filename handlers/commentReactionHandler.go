package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/database"
	"forum/helpers"
)

func CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookieValue := helpers.GetCookieValue(w, r)
	if cookieValue == "" {
		return
	}

	var userID int
	err := database.DataBase.QueryRow("SELECT id FROM users WHERE session = ?", cookieValue).Scan(&userID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	commentIDStr := r.FormValue("commentId")
	reactionStr := r.FormValue("reaction")

	if commentIDStr == "" || reactionStr == "" {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || commentID <= 0 {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	reaction, err := strconv.Atoi(reactionStr)
	if err != nil || (reaction != 1 && reaction != -1) {
		http.Error(w, "Invalid reaction value", http.StatusBadRequest)
		return
	}

	var existingReaction int
	err = database.DataBase.QueryRow(
		"SELECT reaction FROM commentReactions WHERE userId = ? AND commentId = ?",
		userID, commentID,
	).Scan(&existingReaction)

	switch err {
	case sql.ErrNoRows:
		_, err = database.DataBase.Exec(
			"INSERT INTO commentReactions (userId, commentId, reaction) VALUES(?, ?, ?)",
			userID, commentID, reaction,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

	case nil:
		if existingReaction == reaction {
			_, err = database.DataBase.Exec(
				"DELETE FROM commentReactions WHERE userId = ? AND commentId = ?",
				userID, commentID,
			)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = database.DataBase.Exec(
				"UPDATE commentReactions SET reaction = ? WHERE userId = ? AND commentId = ?",
				reaction, userID, commentID,
			)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

	default:
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
