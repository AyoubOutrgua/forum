package handlers

import (
	"database/sql"
	"net/http"

	"forum/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	/* if r.Method != "POST" {
		helpers.Errorhandler(w, "statusPage.html", http.StatusMethodNotAllowed)
		return
	} */
	//  check if the user is alrady logged in
	/* if exists, _ := helpers.SessionChecked(w, r); exists {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} */

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		helpers.Errorhandler(w, "login.html", http.StatusBadRequest)
		return
	}

	stmt := `SELECT password FROM users WHERE username = ? OR email = ?`
	row := Db.QueryRow(stmt, username, username)

	var hashPass string
	err := row.Scan(&hashPass)
	if err == sql.ErrNoRows {
		helpers.Errorhandler(w, "login.html", http.StatusBadRequest)
		return
	} else if err != nil {
		helpers.Errorhandler(w, "statusPage.html", http.StatusInternalServerError)

		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password)) != nil {
		helpers.Errorhandler(w, "login.html", http.StatusBadRequest)
		return
	}

	sessionID := uuid.New().String()
	stmt2 := `UPDATE users SET session = ? WHERE username = ? or   email = ?`
	_, err = Db.Exec(stmt2, sessionID, username, username)
	if err != nil {
		helpers.Errorhandler(w, "login.html", http.StatusInternalServerError)
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
