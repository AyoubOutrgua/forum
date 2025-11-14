package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"forum/database"
	"forum/helpers"
)

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

		fmt.Println(cookie.Value)
		var expiredTime time.Time
		fmt.Println("error here 3awd")
		err = database.DataBase.QueryRow(
			"SELECT dateexpired FROM users WHERE session = ?", cookie.Value,
		).Scan(&expiredTime)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
			} else if err != nil {
				helpers.Errorhandler(w, "Unexpected error", http.StatusInternalServerError)
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
