package handlers

import (
	"net/http"
	"os"
)

func StyleFunc(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	file, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if file.IsDir() {
		return
	}
	http.ServeFile(w, r, filePath)
}