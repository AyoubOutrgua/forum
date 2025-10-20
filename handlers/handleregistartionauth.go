package handlers

import (
	"net/http"
	"strings"

	"forum/database"
	"forum/helpers"

	"golang.org/x/crypto/bcrypt"
)

// Db *sql.DB

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	cookie, errorsession := r.Cookie("session")
	if errorsession == nil && cookie.Value != "" {
		var userExists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method not allowed", 400)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("firstpass")
	firstpassword := r.FormValue("secondpass")
	if firstpassword != password {
		helpers.Render(w, "register.html", http.StatusBadRequest, map[string]string{"Error": "Passwords do not match"})
		return
	}

	if password == "" || email == "" || firstpassword == "" || username == "" {
		helpers.Render(w, "register.html", http.StatusBadRequest, map[string]string{"Error": "All fields are required"})
		return
	}

	errreg := helpers.ValidateInfo(username, email, password)
	if errreg != "" {
		helpers.Render(w, "register.html", http.StatusBadRequest, map[string]string{"Error": errreg})
		return
	}

	var existsUsername bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE userName = ?)", username).Scan(&existsUsername)
	if err != nil {
		helpers.Errorhandler(w, "Database  error", http.StatusInternalServerError)
		return
	}
	if existsUsername {
		helpers.Render(w, "register.html", http.StatusBadRequest, map[string]string{"Error": "Username already taken"})
		return
	}

	var existsEmail bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&existsEmail)
	if err != nil {
		helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
		return
	}
	if existsEmail {
		helpers.Render(w, "register.html", http.StatusBadRequest, map[string]string{"Error": "Email already used"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		helpers.Render(w, "register.html", http.StatusBadRequest, map[string]string{"Error": "inxpected error please try again"})
		return
	}

	stmt2 := `INSERT INTO users (userName, email, password) VALUES (?, ?, ?);`
	_, err = database.DB.Exec(stmt2, username, email, string(hashPassword))
	if err != nil {
		helpers.Errorhandler(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
