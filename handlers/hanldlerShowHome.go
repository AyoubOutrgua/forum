package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/helpers"
)

type PageData struct {
	LoggedIn bool
}

func HanldlerShowHome(w http.ResponseWriter, r *http.Request) {
	loggedIn := false
	if r.URL.Path != "/" {
		helpers.Errorhandler(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, errSession := r.Cookie("session")
	if errSession == nil && cookie.Value != "" {
		var userExists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {
			loggedIn = true
		}
	}
	data := PageData{LoggedIn: loggedIn}
	temp, errPerse := template.ParseFiles("templates/index.html")
	if errPerse != nil {
		helpers.Errorhandler(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	errExec := temp.Execute(w, data)
	if errExec != nil {
		helpers.Errorhandler(w, "Status Internal Server Error!", http.StatusInternalServerError)
		return
	}
}
