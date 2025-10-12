package handlers

import (
	"database/sql"
	"net/http"

	"forum/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)



func LoginHandler(w http.ResponseWriter, r *http.Request) {
	cookie, errSession := r.Cookie("session")
	if errSession == nil && cookie.Value != "" {
		var userExists bool
		err := Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method not allowed", 400)
		return
	}

	
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "All fields are required"})
		return
	}

	stmt := `SELECT password FROM users WHERE userName = ? OR email = ?`
	row := Db.QueryRow(stmt, username, username)

	var hashPass string
	err := row.Scan(&hashPass)
	if err == sql.ErrNoRows {
		helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "Invalid username or password"})
		return
	} else if err != nil {
		helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password)) != nil {
		helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "Invalid username or password"})
		return
	}

	sessionID := uuid.New().String()
	stmt2 := `UPDATE users SET session = ? WHERE userName = ? OR email = ?`
	_, err = Db.Exec(stmt2, sessionID, username, username)
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
