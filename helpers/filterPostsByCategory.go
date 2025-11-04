package helpers

import (
    "fmt"
    "net/http"
    "strings"

    "forum/database"
    "forum/tools"
)

func GetPostsByCategoryIDs(w http.ResponseWriter, catIDs []int) []tools.Post {
    if len(catIDs) == 0 {
        return GetAllPosts(w)
    }

    parts := make([]string, 0, len(catIDs))
    for _, id := range catIDs {
        parts = append(parts, fmt.Sprintf("%d", id))
    }
    in := strings.Join(parts, ",")

    query := fmt.Sprintf(`
        SELECT p.id, p.title, p.post, p.imageUrl, u.userName, p.creationDate
        FROM posts AS p
        INNER JOIN users AS u ON u.id = p.userId
        INNER JOIN postCategories AS pc ON pc.postId = p.id
        WHERE pc.categoryId IN (%s)
        GROUP BY p.id
        ORDER BY p.creationDate DESC;
    `, in)

    posts, errSelect := database.SelectAllPosts(query)
    if errSelect != nil {
        Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
        return nil
    }
    GetPostCategories(w, posts)
    return posts
}