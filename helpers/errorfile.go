package helpers

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
)
// Errorhandler renders a custom error page with the provided error message and status code.
func Errorhandler(w http.ResponseWriter, errors string, er int) {
	const filePath = "templates/error.html"

	myMap := map[string]string{
		"errorText":  errors,
		"statusCode": strconv.Itoa(er),
	}
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		http.Error(w, "500 Internal Server Error (parse error)", http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if execErr := tmpl.Execute(&buf, myMap); execErr != nil {
		http.Error(w, "500 Internal Server Error (exec error)", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(er)
	w.Write(buf.Bytes())
}
