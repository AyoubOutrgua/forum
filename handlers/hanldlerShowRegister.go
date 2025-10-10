package handlers

import (
	"html/template"
	"net/http"
)

func HanldlerShowRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	} 
	temp, errPerse := template.ParseFiles("templates/register.html")
	if errPerse != nil {
		http.Error(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	errExec := temp.Execute(w, nil)
	if errExec != nil {
		http.Error(w, "Status Internal Server Error!", http.StatusInternalServerError)
		return
	}
}
