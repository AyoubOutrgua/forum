package helpers

import (
	"net/http"

	"forum/database"
	"forum/tools"
)

func GetAllPosts(w http.ResponseWriter) []tools.Post {
	postsQuery := `
			SELECT p.id, p.title, p.post, p.imageUrl, u.userName, p.creationDate
			FROM posts AS p
			INNER JOIN users AS u ON u.id = p.userId
			ORDER BY p.creationDate DESC;
			`
	posts, errSelect := database.SelectAllPosts(postsQuery)
	if errSelect != nil {
		http.Error(w, "------------- ERROR --------------!", http.StatusNotFound)
		return nil
	}
	return posts
}
