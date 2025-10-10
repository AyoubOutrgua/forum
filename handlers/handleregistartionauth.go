package handlers

import (
	"fmt"
	"net/http"

	//"regexp"

	"forum/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	cookie, errorsession := r.Cookie("session")
	if errorsession == nil && cookie.Value != "" {

		var userExists bool
		err := Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	if r.Method != "POST" {
		helpers.Errorhandler(w, "statusPage.html", http.StatusMethodNotAllowed)
		return
	}
	password := r.FormValue("secondpassword")
	email := r.FormValue("email")
	username := r.FormValue("username")
	firstpassword := r.FormValue("firstpassword")
	fmt.Println(password, email, username, firstpassword, password)
	if firstpassword != password {
		helpers.Errorhandler(w, "imatched pass", 400)
		return
	}
	if password == "" || email == "" || firstpassword == "" || username == "" {
		helpers.Errorhandler(w, "incorrect information", 400)
		return
	}
	errreg := helpers.ValidateInfo(username, email, password)
	if errreg != "" {
		helpers.Errorhandler(w, errreg, 400)
		return
	}
	var existsUsername bool
	err := Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&existsUsername)
	if err != nil {
		helpers.Errorhandler(w, "database error", http.StatusInternalServerError)
		return
	}
	if existsUsername {
		helpers.Errorhandler(w, "username already taken", 400)
		return
	}

	var existsEmail bool
	err = Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&existsEmail)
	if err != nil {
		helpers.Errorhandler(w, "database error", http.StatusInternalServerError)
		return
	}
	if existsEmail {
		helpers.Errorhandler(w, "email already used", 400)
		return
	}

	hashPassword, Err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if Err != nil {
		helpers.Errorhandler(w, "error generate hash pass", http.StatusInternalServerError)
		return
	}
	stmt2 := `INSERT INTO users (username, email, password) VALUES (?, ?, ?);`
	_, err = Db.Exec(stmt2, username, email, string(hashPassword))
	if err != nil {
		helpers.Errorhandler(w, "error db exce", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	sessionID := uuid.New().String()
	stmt3 := `UPDATE users SET session = ? WHERE username = ?`
	_, err = Db.Exec(stmt3, sessionID, username)
	if err != nil {
		helpers.Errorhandler(w, "statusPage.html", http.StatusInternalServerError)
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
