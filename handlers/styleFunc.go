package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func StyleFunc(w http.ResponseWriter, r *http.Request) {
	filePath := ".." + r.URL.Path
	fmt.Println(filePath)
	file, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if file.IsDir() {
		return
	}
	http.ServeFile(w, r, filePath)
}
