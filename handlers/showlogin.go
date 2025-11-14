package handlers

import (
	"net/http"

	"forum/helpers"
)

func Showloginhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return

	}

	helpers.Render(w, "login.html", http.StatusOK, map[string]string{"Error": "", "Username": ""})
}
