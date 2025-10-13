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

	// postsQuery := `
	// 		SELECT p.id, p.title, p.post, p.imageUrl, u.userName, p.creationDate
	// 		FROM posts AS p
	// 		INNER JOIN users AS u ON u.id = p.userId
	// 		ORDER BY p.creationDate DESC;
	// 		`
	// posts, errSelect := database.SelectAllPosts(postsQuery)
	// if errSelect != nil {
	// 	http.Error(w, "------------- ERROR --------------!", http.StatusNotFound)
	// 	return
	// }

	// categoriesQuery := `
	// 		SELECT c.id, c.category
	// 		FROM categories AS c;
	// 		`
	// categories, errSelect := database.SelectAllCategories(categoriesQuery)
	// if errSelect != nil {
	// 	http.Error(w, "------------- ERROR --------------!", http.StatusNotFound)
	// 	return
	// }
	posts := helpers.GetAllPosts(w)
	categories := helpers.GetAllCategories(w)
	var index tools.Insex
	index.Posts = posts
	index.Categories = categories

	errExec := temp.Execute(w, index)
	if errExec != nil {
		http.Error(w, "Status Internal Server Error!!!!!", http.StatusInternalServerError)
		return
	}
}
