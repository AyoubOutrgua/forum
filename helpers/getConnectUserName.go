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
		// http.Error(w, "------------- ERROR USER NAME--------------!", http.StatusNotFound)
		return ""
	}
	return connectUserName
}
