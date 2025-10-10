package handlers

import (
	"forum/helpers"
	"html/template"
	"net/http"
)

func HanldlerShowLogin(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	temp, errPerse := template.ParseFiles("templates/login.html")
	if errPerse != nil {
		helpers.Errorhandler(w,"Status Not Found!",http.StatusNotFound)
		return
	}
	errExec := temp.Execute(w,nil)
	if errExec != nil {
		helpers.Errorhandler(w, "Status Internal Server Error!",http.StatusInternalServerError)
		return
	}
}