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
		Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
		return nil
	}
	GetPostCategories(w, posts)
	return posts
}

func GetPostCategories(w http.ResponseWriter, posts []tools.Post) {
	for i := range posts {
		var postCategories []string
		id := posts[i].ID
		postCategoriesQuery := `
			SELECT categoryId
			FROM postCategories
			WHERE postId = ?
			`
		categoriesID, err := database.SelectPostCategories(postCategoriesQuery, id)
		if err != nil {
			Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
			return
		}
		categories := GetAllCategories(w)
		for _, category := range categories {
			for _, catID := range categoriesID {
				if category.ID == catID {
					postCategories = append(postCategories, category.Category)
				}
			}
		}
		posts[i].Categories = postCategories
	}
}
