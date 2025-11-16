package handlers

import (
	"net/http"
	"strings"

	"forum/database"
	"forum/helpers"

	"golang.org/x/crypto/bcrypt"
)
// RegisterHandler processes user registration by validating input, checking for existing users, hashing the password, and storing the new user in the database.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method not allowed", http.StatusNotFound)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("firstpass"))
	firstpassword := strings.TrimSpace(r.FormValue("secondpass"))
	if firstpassword != password {
		helpers.Render(w, "register.html", http.StatusBadRequest, map[string]string{"Error": "Passwords do not match", "Username": username, "email": email})
		return
	}

	if password == "" || email == "" || firstpassword == "" || username == "" {

		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	errreg := helpers.ValidateInfo(username, email, password)
	if !errreg {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var existsUsername bool
	err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE userName = ?)", username).Scan(&existsUsername)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if existsUsername {
		helpers.Render(w, "register.html", http.StatusConflict, map[string]string{"Error": "Username already taken", "Username": "", "email": email})
		return
	}

	var existsEmail bool
	err = database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&existsEmail)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if existsEmail {
		helpers.Render(w, "register.html", http.StatusConflict, map[string]string{"Error": "Email already used", "Username": username, "email": ""})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		helpers.Errorhandler(w,"Internal server Error",http.StatusInternalServerError)
		return
	}

	stmt2 := `INSERT INTO users (userName, email, password) VALUES (?, ?, ?);`
	_, err = database.DataBase.Exec(stmt2, username, email, string(hashPassword))
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
