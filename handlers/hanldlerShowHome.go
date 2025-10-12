package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
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
	query := `
			SELECT p.id, p.title, p.post, u.userName, p.creationDate
			FROM posts AS p
			INNER JOIN users AS u ON u.id = p.userId
			ORDER BY p.creationDate DESC;
			`
	posts, errSelect := database.SelectData(query)
	if errSelect != nil {
		http.Error(w, "------------- ERROR --------------!", http.StatusNotFound)
		return
	}

	errExec := temp.Execute(w, posts)
	if errExec != nil {
		http.Error(w, "Status Internal Server Error!!!!!", http.StatusInternalServerError)
		return
	}
}
