package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/database"
	"forum/helpers"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	postIDStr := r.FormValue("postId")
	reactionStr := r.FormValue("reaction")
	if postIDStr == "" || reactionStr == "" {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	reaction, err := strconv.Atoi(reactionStr)
	if err != nil || (reaction != 1 && reaction != -1) {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var postExists int
	errSelect := database.DataBase.QueryRow("SELECT COUNT(*) FROM posts WHERE id = ?", postID).Scan(&postExists)
	if errSelect == sql.ErrNoRows {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	} else if errSelect != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if postExists == 0 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var existingReaction int
	err = database.DataBase.QueryRow("SELECT reaction FROM postReactions WHERE userId = ? AND postId = ?", userID, postID).Scan(&existingReaction)

	switch err {
	case sql.ErrNoRows:
		_, err = database.DataBase.Exec("INSERT INTO postReactions (userId, postId, reaction) VALUES(?, ?, ?)", userID, postID, reaction)
		if err != nil {
			helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case nil:
		if existingReaction == reaction {
			_, err = database.DataBase.Exec("DELETE FROM postReactions WHERE userId = ? AND postId = ?", userID, postID)
			if err != nil {
				helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = database.DataBase.Exec("UPDATE postReactions SET reaction = ? WHERE userId = ? AND postId = ?", reaction, userID, postID)
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
