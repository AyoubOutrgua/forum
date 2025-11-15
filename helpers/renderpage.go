package helpers

import (
	"bytes"
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, templateFile string, status int, data interface{}) {
	tmpl, err := template.ParseFiles("templates/" + templateFile)
	if err != nil {
		Errorhandler(w, "Template parsing error", http.StatusInternalServerError)

		return
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}
