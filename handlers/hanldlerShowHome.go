package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/helpers"
	"forum/tools"
)

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
		err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {
			loggedIn = true
		}
	}
	data := tools.IsLogin{LoggedIn: loggedIn}

	temp, errPerse := template.ParseFiles("templates/index.html")
	if errPerse != nil {
		helpers.Errorhandler(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	posts := helpers.GetAllPosts(w)
	categories := helpers.GetAllCategories(w)
	var pageData tools.PageData
	pageData.Posts = posts
	pageData.Categories = categories
	pageData.IsLogin = data

	errExec := temp.Execute(w, pageData)
	if errExec != nil {
		http.Error(w, "Status Internal Server Error!!!!!", http.StatusInternalServerError)
	}
}
