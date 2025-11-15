package handlers

import (
	"forum/helpers"
	"net/http"
	"os"
)

func StyleFunc(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	file, err := os.Stat(filePath)
	if err != nil {
			helpers.Errorhandler(w,"Page not Found",http.StatusNotFound)
		return
	}
	if file.IsDir() {
		helpers.Errorhandler(w,"Access Forbidden",http.StatusForbidden)
		return
	}
	http.ServeFile(w, r, filePath)
}
