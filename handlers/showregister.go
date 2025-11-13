package handlers

import (
	"net/http"

	"forum/helpers"
)

func Showregister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w,"bad request",http.StatusBadRequest)
		return
	}
	catStrs := r.URL.Query()
	if len(catStrs) != 0 {
		helpers.Errorhandler(w, "Bad request", http.StatusBadRequest)
		return
	}
	helpers.Render(w, "register.html", http.StatusOK, map[string]string{"Error": "", "Username": "", "email": ""})
}
