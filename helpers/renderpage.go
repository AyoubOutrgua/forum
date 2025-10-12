package helpers

import (
	
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, templateFile string,errors int , data interface{}) {
	w.WriteHeader(errors)
	tmpl, err := template.ParseFiles("templates/" + templateFile)
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		Errorhandler(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}
