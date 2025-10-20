package handlers

import (
	"net/http"
	"time"

	"forum/database"
	"forum/helpers"
)

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err = database.DB.Exec("UPDATE users SET session = NULL WHERE session = ?", cookie.Value)
	if err != nil {
		helpers.Errorhandler(w, "Database error while logging out", http.StatusInternalServerError)
		return
	}

	expiredCookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, expiredCookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
