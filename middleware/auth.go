package middleware

import (
	"database/sql"
	"net/http"
	"time"

	"forum/database"
	"forum/helpers"
)
// Checksession is a middleware that verifies the user's session cookie and its validity before allowing access to protected routes.
func Checksession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var userExists bool
		err = database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err != nil || !userExists {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var expiredTime time.Time
		err = database.DataBase.QueryRow(
			"SELECT dateexpired FROM users WHERE session = ?", cookie.Value,
		).Scan(&expiredTime)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else if err != nil {
			helpers.Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
			return
		}

		if expiredTime.Before(time.Now()) {
			_, _ = database.DataBase.Exec(
				"UPDATE users SET session = NULL, dateexpired = NULL WHERE session = ?", cookie.Value,
			)
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

			return
		}

		next(w, r)
	}
}
// CheckLogin is a middleware that prevents logged-in users from accessing routes meant for unauthenticated users, such as login or registration pages.

func CheckLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err == nil && cookie.Value != "" {
			var userExists bool
			err = database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
			if err == nil && userExists {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
		next(w, r)
	}
}