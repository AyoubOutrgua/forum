package handlers

import (
	"forum/helpers"
	"html/template"
	"net/http"
)

func HanldlerShowHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		helpers.Errorhandler(w, "Status Not Found!", http.StatusNotFound)
		return
	}
 if r.Method != http.MethodGet {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	} 
	temp, errPerse := template.ParseFiles("templates/index.html")
	if errPerse != nil {
		helpers.Errorhandler(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	errExec := temp.Execute(w, nil)
	if errExec != nil {
		helpers.Errorhandler(w, "Status Internal Server Error!", http.StatusInternalServerError)
		return
	}
}
