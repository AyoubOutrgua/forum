package helpers

import "forum/database"

func GetUserID(cookieID string) int {
	query := `SELECT id FROM users WHERE session = ?`
	return database.SelectUserID(query, cookieID)
}
