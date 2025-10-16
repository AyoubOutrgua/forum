package handlers

import (
	"html/template"
	"net/http"

	"forum/helpers"
	"forum/tools"
)

func HanldlerShowHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	temp, errPerse := template.ParseFiles("templates/index.html")
	if errPerse != nil {
		http.Error(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	posts := helpers.GetAllPosts(w)
	categories := helpers.GetAllCategories(w)
	var index tools.Index
	index.Posts = posts
	index.Categories = categories

	errExec := temp.Execute(w, index)
	if errExec != nil {
		http.Error(w, "Status Internal Server Error!!!!!", http.StatusInternalServerError)
		return
	}
}
