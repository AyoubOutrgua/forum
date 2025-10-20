package middleware

import (
	"net/http"

	"forum/database"
)

func RedirectIfLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err == nil && cookie.Value != "" {
			var userExists bool
			err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
			if err == nil && userExists {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
		next(w, r)
	}
}
