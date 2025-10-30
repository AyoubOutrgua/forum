package helpers

import (
	"net/http"

	"forum/database"
)

func GetUserID(cookieID string) int {
	query := `SELECT id FROM users WHERE session = ?`
	userId, err := database.SelectUserID(query, cookieID)
	if err != nil {
		Errorhandler(nil, "Status Internal Server Error", http.StatusInternalServerError)
		return 0
	}
	return userId
}
