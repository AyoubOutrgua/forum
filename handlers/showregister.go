package handlers

import (
	"net/http"

	"forum/helpers"
)
// Showregister displays the registration page to the user.
func Showregister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	helpers.Render(w, "register.html", http.StatusOK, map[string]string{"Error": "", "Username": "", "email": ""})
}
