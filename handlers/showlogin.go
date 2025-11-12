package handlers

import (
	"net/http"

	"forum/helpers"
)

func Showloginhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w,"bad request",http.StatusBadRequest)
		return
		
	}
	
		helpers.Render(w, "login.html", http.StatusOK, map[string]string{"Error": "", "Username": ""})
		
}
