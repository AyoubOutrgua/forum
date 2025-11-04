package helpers

import (
    "fmt"
    "net/http"

    "forum/database"
    "forum/tools"
)

func GetPostsByUserID(w http.ResponseWriter, userID int) []tools.Post {
    if userID <= 0 {
        return []tools.Post{}
    }

    query := fmt.Sprintf(`
        SELECT p.id, p.title, p.post, p.imageUrl, u.userName, p.creationDate
        FROM posts AS p
        INNER JOIN users AS u ON u.id = p.userId
        WHERE p.userId = %d
        ORDER BY p.creationDate DESC;
    `, userID)

    posts, errSelect := database.SelectAllPosts(query)
    if errSelect != nil {
        Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
        return nil
    }
    GetPostCategories(w, posts)
    return posts
}