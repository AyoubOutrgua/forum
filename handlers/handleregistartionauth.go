package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"forum/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	/* 	if exists, _ := helpers.SessionChecked(w, r); exists {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} */
	if r.Method != "POST" {
		helpers.Errorhandler(w, "statusPage.html", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("fkdkfdkf")
	// get the data
	password := r.FormValue("password")
	email := r.FormValue("email")
	username := r.FormValue("username")
	firstpassword := r.FormValue("firstpassword")
	fmt.Println(password,email,username,firstpassword)
	fmt.Println(username, email)
	var ErrorMessage string
	// regex of email
	emailregex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// passregex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$`
	// all possible error messages
	if password == "" || email == "" || username == "" || firstpassword == "" {
		ErrorMessage = "All inputs are required"
	} else if len(email) >= 50 || len(email) <= 10 {
		ErrorMessage = "Email must be between 5 and 50 characters"
	} else if match, _ := regexp.MatchString(emailregex, email); !match {
		ErrorMessage = "Invalid email format"
	} else if firstpassword != password {
		ErrorMessage = "Passwords do not match"
	} else if len(username) < 3 || len(username) > 15 {
		ErrorMessage = "Username must be at least 3 characters"
	} else if len(password) <= 6 || len(password) > 15 {
		ErrorMessage = "Password must be at least 6 characters"
	}
	// check the username if is already used
	stmt := "SELECT id FROM users WHERE username = ? OR email = ?"
	row := Db.QueryRow(stmt, username, email)
	var id string
	err := row.Scan(&id)
	fmt.Println(id)
	fmt.Println(err)
if ErrorMessage !=  ""{
	helpers.Errorhandler(w,"djd",404)
}
  fmt.Println(ErrorMessage)
	/* if err != sql.ErrNoRows {
		ErrorMessage = "The username or email is already used"
		helpers.Errorhandler(w, "register.html", http.StatusBadRequest)
		fmt.Println("test1")
		return
	}
	// if there is  an error
	if ErrorMessage != "" {
		helpers.Errorhandler(w, "register.html", http.StatusBadRequest)
		fmt.Println("test2")
		return
	} */
	// bcrypte the password
	hashPassword, Err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if Err != nil {
		helpers.Errorhandler(w, "statusPage.html", http.StatusInternalServerError)
		return
	}
	// insert the data in the database
	stmt2 := `INSERT INTO users (username, email, password) VALUES (?, ?, ?);`
	_, err = Db.Exec(stmt2, username, email, string(hashPassword))
	if err != nil {
		helpers.Errorhandler(w, "statusPage.html", http.StatusInternalServerError)
		return
	}
	// create a session yith uuid
	sessionID := uuid.New().String()
	stmt3 := `UPDATE users SET session = ? WHERE username = ?`
	_, err = Db.Exec(stmt3, sessionID, username)
	if err != nil {
		helpers.Errorhandler(w, "statusPage.html", http.StatusInternalServerError)
		return
	}
	// create a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
