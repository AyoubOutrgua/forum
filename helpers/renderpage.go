package helpers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/" + templateFile)
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		fmt.Println("error parsing template:", err)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		Errorhandler(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func RenderPage(w http.ResponseWriter, r *http.Request, templateFile string, data interface{}) {
	Render(w, templateFile, data)
}
