package helpers

import (
	"net/http"

	"forum/database"
)

func GetConnectUserName(w http.ResponseWriter, userId int) string {
	userNameQuery := `
			SELECT userName
			FROM users
			WHERE id = ?;
			`
	connectUserName, errSelect := database.SelectUserName(userNameQuery, userId)
	if errSelect != nil {
		// Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
		return ""
	}
	return connectUserName
}
