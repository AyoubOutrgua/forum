package helpers

import (
	"net/http"

	"forum/database"
)

func GetLastPostDates(userID int) []string {
	dateQuery := `
			SELECT creationDate
			FROM posts
			WHERE userId = ?
			ORDER BY creationDate DESC
			LIMIT 5;
			`

	dates, err := database.SelectLastDates(dateQuery, userID)
	if err != nil {
		Errorhandler(nil, "Status Internal Server Error", http.StatusInternalServerError)
		return nil
	}
	return dates
}
