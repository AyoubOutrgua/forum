package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"forum/database"
	"forum/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if database.DataBase == nil {
		helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
	}
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "page not found", http.StatusNotFound)
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "All fields are required", "Username": username})
		return
	}
	stmt := `SELECT password FROM users WHERE userName = ? OR email = ?`
	row := database.DataBase.QueryRow(stmt, username, username)

	var hashPass string
	err := row.Scan(&hashPass)
	if err == sql.ErrNoRows {
		helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "Invalid username or password", "Username": username})
		return
	} else if err != nil {
		helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password)) != nil {
		helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "All fields are required", "Username": username})

		return
	}

	sessionID := uuid.New().String()
	expireTime := time.Now().Add(1* time.Hour)
	stmt2 := `UPDATE users SET dateexpired = ? ,session = ? WHERE userName = ? OR email = ?`
	_, err = database.DataBase.Exec(stmt2, expireTime, sessionID, username, username)
	if err != nil {
		helpers.Errorhandler(w, "Database update error", http.StatusInternalServerError)
		fmt.Println("test", err)
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
