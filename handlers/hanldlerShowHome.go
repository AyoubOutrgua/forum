package handlers

import (
	"html/template"
	"net/http"
)

func HanldlerShowHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Status Not Found!", http.StatusNotFound)
		return
	}
 if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	} 
	temp, errPerse := template.ParseFiles("templates/index.html")
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
