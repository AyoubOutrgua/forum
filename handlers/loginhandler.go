package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"forum/database"
	"forum/helpers"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// if database.DataBase == nil {
	// 	helpers.Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if len(username) > 50 || len(username) < 4 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if len(password) > 20 || len(password) < 6 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
	}
	stmt := `SELECT password FROM users WHERE userName = ? OR email = ?`
	row := database.DataBase.QueryRow(stmt, username, username)

	var hashPass string
	err := row.Scan(&hashPass)
	if err == sql.ErrNoRows {
		helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "Invalid username or password", "Username": username})
		return
	} else if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password)) != nil {
		helpers.Render(w, "login.html", http.StatusUnauthorized, map[string]string{"Error": "Invalid username or password", "Username": username})
		return
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		helpers.Errorhandler(w, "internal server Error", http.StatusInternalServerError)
	}
	strsessionID := sessionID.String()
	expireTime := time.Now().Add(1 * time.Hour)
	stmt2 := `UPDATE users SET dateexpired = ? ,session = ? WHERE userName = ? OR email = ?`
	_, err = database.DataBase.Exec(stmt2, expireTime, strsessionID, username, username)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    strsessionID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
