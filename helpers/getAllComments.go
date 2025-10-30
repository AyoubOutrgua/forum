package helpers

import (
	"net/http"

	"forum/database"
	"forum/tools"
)

func GetAllComments(w http.ResponseWriter) map[int][]tools.Comment {
	commentsQuery := `
			SELECT c.id, c.comment, c.postId, c.userId, u.userName, c.creationDate
			FROM comments AS c
			INNER JOIN users AS u ON c.userId = u.id
			ORDER BY c.postId, c.creationDate DESC;
			`
	comments, errSelect := database.SelectAllComments(commentsQuery)
	if errSelect != nil {
		Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
		return nil
	}
	return comments
}
