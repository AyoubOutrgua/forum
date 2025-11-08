package helpers

import (
	"net/http"

	"forum/database"
)

func GetLastCommentsDates(userID int) []string {
	dateQuery := `
			SELECT creationDate
			FROM comments
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
