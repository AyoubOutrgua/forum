package handlers

import (
	"database/sql"
	"forum/database"
	"net/http"
	"strconv"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
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
	postIDStr := r.FormValue("postId")
	reactionStr := r.FormValue("reaction")
	if postIDStr == "" || reactionStr == "" {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	
	reaction, err := strconv.Atoi(reactionStr)
	if err != nil || (reaction != 1 && reaction != -1) {
		http.Error(w, "Invalid reaction value", http.StatusBadRequest)
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
	var existingReaction int
	err = database.DataBase.QueryRow("SELECT reaction FROM postReactions WHERE userId = ? AND postId = ?", userID, postID).Scan(&existingReaction)
	
	switch err {
	case sql.ErrNoRows:
		_, err = database.DataBase.Exec("INSERT INTO postReactions (userId, postId, reaction) VALUES(?, ?, ?)", userID, postID, reaction)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	case nil:
		if existingReaction == reaction {
			_, err = database.DataBase.Exec("DELETE FROM postReactions WHERE userId = ? AND postId = ?", userID, postID)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = database.DataBase.Exec("UPDATE postReactions SET reaction = ? WHERE userId = ? AND postId = ?", reaction, userID, postID)
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