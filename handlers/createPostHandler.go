package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"forum/database"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("CREATE POST HANDLER")
	errParse := r.ParseMultipartForm(10 << 20)
	if errParse != nil {
		http.Error(w, "Status Bad Request 1", http.StatusBadRequest)
		return
	}

	title, ok := r.PostForm["title"]
	if !ok {
		http.Error(w, "Status Bad Request 2", http.StatusBadRequest)
		return
	}
	description, ok := r.PostForm["description"]
	if !ok {
		http.Error(w, "Status Bad Request 3", http.StatusBadRequest)
		return
	}
	categories, ok := r.PostForm["categories"]
	if !ok {
		http.Error(w, "Status Bad Request 4", http.StatusBadRequest)
		return
	}
	imageFile, handler, err := r.FormFile("choose-file")
	if err != nil {
		http.Error(w, "Status Bad Request 5", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	file, errCreate := os.Create("./upload/" + handler.Filename)
	if errCreate != nil {
		http.Error(w, "Status Bad Request 6", http.StatusBadRequest)
		return
	}
	defer file.Close()

	_, errCopy := io.Copy(file, imageFile)
	if errCopy != nil {
		http.Error(w, "Status Bad Request 7", http.StatusBadRequest)
		return
	}
	timeNow := time.Now().Format("2006-01-02")
	query := "INSERT INTO posts (title,post,userId,creationDate) VALUES (?, ?, ?, ?)"
	// fmt.Println("------------------------\n", query, "\n***************************")
	fmt.Println(categories)
	database.ExecuteData(query, title[0], description[0], 2, timeNow)

	// query2 := "INSERT INTO users (userName, email, password) VALUES (?, ?, ?)"
	// name := "hamid"
	// email := "hamid@example.com"
	// pass := "123456"
	// database.ExecuteData(query2, name, email, pass)
}
