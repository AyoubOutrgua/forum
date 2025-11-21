package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/database"
	"forum/helpers"
)

// CommentReactionHandler handles like/dislike reactions for comments: it checks the request,
// verifies the user session, validates the form data, and then inserts, updates, or deletes
// the reaction depending on the user's previous action.

func CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	errParse := r.ParseForm()
	if errParse != nil {
		helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
		return
	}

	commentIDStr := r.FormValue("commentId")
	reactionStr := r.FormValue("reaction")

	if commentIDStr == "" || reactionStr == "" {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || commentID <= 0 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var commentExists int
	selectError := database.DataBase.QueryRow("SELECT COUNT(*) FROM comments WHERE id = ?", commentID).Scan(&commentExists)
	if selectError == sql.ErrNoRows {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	} else if selectError != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if commentExists == 0 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	reaction, err := strconv.Atoi(reactionStr)
	if err != nil || (reaction != 1 && reaction != -1) {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
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
			helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	case nil:
		if existingReaction == reaction {
			_, err = database.DataBase.Exec(
				"DELETE FROM commentReactions WHERE userId = ? AND commentId = ?",
				userID, commentID,
			)
			if err != nil {
				helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = database.DataBase.Exec(
				"UPDATE commentReactions SET reaction = ? WHERE userId = ? AND commentId = ?",
				reaction, userID, commentID,
			)
			if err != nil {
				helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

	default:
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
