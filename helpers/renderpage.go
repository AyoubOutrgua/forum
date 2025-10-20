package helpers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, templateFile string, Status int, data interface{}) {
	w.WriteHeader(Status)

	tmpl, err := template.ParseFiles("templates/" + templateFile)
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		fmt.Println("error here asahbi")
		return
	}
	err = tmpl.Execute(w, data)
	data = nil

	if err != nil {
		Errorhandler(w, "Template execution error", http.StatusInternalServerError)

		return
	}
}
