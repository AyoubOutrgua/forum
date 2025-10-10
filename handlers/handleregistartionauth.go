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
	if r.Method != "POST" {
		helpers.Errorhandler(w, "statusPage.html", http.StatusMethodNotAllowed)
		return
	}
	password := r.FormValue("secondpassword")
	email := r.FormValue("email")
	username := r.FormValue("username")
	firstpassword := r.FormValue("firstpassword")
	fmt.Println(password,email,username,firstpassword,password)

	

	hashPassword, Err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if Err != nil {
		helpers.Errorhandler(w, "statusPage.html", http.StatusInternalServerError)
		return
	}
	stmt2 := `INSERT INTO users (username, email, password) VALUES (?, ?, ?);`
	_, err := Db.Exec(stmt2, username, email, string(hashPassword))
	if err != nil {
		helpers.Errorhandler(w, "statusPage.html", http.StatusInternalServerError)
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
