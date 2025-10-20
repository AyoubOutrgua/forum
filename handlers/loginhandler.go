package handlers

import (
	"database/sql"
	"net/http"

	"forum/database"
	"forum/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if database.DB == nil {
		helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
	}
	cookie, errSession := r.Cookie("session")
	if errSession == nil && cookie.Value != "" {
		var userExists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	if r.Method == http.MethodGet {
		flash := helpers.GetFlash(w, r)
		data := map[string]string{
			"Error":    flash.Message,
			"Username": flash.Username,
			"email":    "",
		}

		helpers.Render(w, "login.html", http.StatusOK, data)
		return
	}
	if r.Method == http.MethodPost {

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "All fields are required", "Username": username})
			return
		}
		stmt := `SELECT password FROM users WHERE userName = ? OR email = ?`
		row := database.DB.QueryRow(stmt, username, username)

		var hashPass string
		err := row.Scan(&hashPass)
		if err == sql.ErrNoRows {
			helpers.SetFlash(w, "Invalid username or password", "", username)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else if err != nil {
			helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password)) != nil {
			helpers.SetFlash(w, "Invalid username or password", "", username)
			http.Redirect(w, r, "/login", http.StatusSeeOther)

			return
		}

		sessionID := uuid.New().String()
		stmt2 := `UPDATE users SET session = ? WHERE userName = ? OR email = ?`
		_, err = database.DB.Exec(stmt2, sessionID, username, username)
		if err != nil {
			helpers.Errorhandler(w, "Database update error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			HttpOnly: true,
			Path:     "/",
			MaxAge:   3600,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}
