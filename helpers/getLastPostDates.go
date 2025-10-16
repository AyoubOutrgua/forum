package helpers

import "forum/database"

func GetLastPostDates(userID int) []string {
	dateQuery := `
			SELECT creationDate
			FROM posts
			WHERE userId = ?
			ORDER BY creationDate DESC
			LIMIT 5;
			`

	dates := database.SelectLastDates(dateQuery, userID)
	return dates
}
