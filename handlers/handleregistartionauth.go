package handlers

import (
	"net/http"
	"strings"

	"forum/database"
	"forum/helpers"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	cookie, errorsession := r.Cookie("session")
	if errorsession == nil && cookie.Value != "" {
		var userExists bool
		err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
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
			"email":    flash.Email,
		}
		helpers.Render(w, "register.html", http.StatusOK, data)
		return
	}
	if r.Method == http.MethodPost {
		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("firstpass")
		firstpassword := r.FormValue("secondpass")
		if firstpassword != password {
			helpers.SetFlash(w, "Passwords do not match", email, username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		if password == "" || email == "" || firstpassword == "" || username == "" {
			helpers.SetFlash(w, "All fields are required", email, username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		errreg := helpers.ValidateInfo(username, email, password)
		if errreg != "" {
			helpers.SetFlash(w, errreg, email, username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		var existsUsername bool
		err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE userName = ?)", username).Scan(&existsUsername)
		if err != nil {
			helpers.Errorhandler(w, "Database  error", http.StatusInternalServerError)
			return
		}
		if existsUsername {
			helpers.SetFlash(w, "Username already taken", email, username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		var existsEmail bool
		err = database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&existsEmail)
		if err != nil {
			helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
			return
		}
		if existsEmail {
			helpers.SetFlash(w, "Email already used", email, username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			helpers.SetFlash(w, "Unexpected error please try again", email, username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		}

		stmt2 := `INSERT INTO users (userName, email, password) VALUES (?, ?, ?);`
		_, err = database.DataBase.Exec(stmt2, username, email, string(hashPassword))
		if err != nil {
			helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
