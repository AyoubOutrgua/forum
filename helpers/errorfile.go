package helpers

import (
	"html/template"
	"net/http"
	"strconv"
)

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

	w.WriteHeader(er)
	if execErr := tmpl.Execute(w, myMap); execErr != nil {
		http.Error(w, "500 Internal Server Error (exec error)", http.StatusInternalServerError)
	}
}
