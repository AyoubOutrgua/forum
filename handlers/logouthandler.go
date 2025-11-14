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

	cookieValue := helpers.GetCookieValue(w, r)
	if cookieValue == "" {
		return
	}

	_, err := database.DataBase.Exec("UPDATE users SET session = NULL WHERE session = ?", cookieValue)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
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
