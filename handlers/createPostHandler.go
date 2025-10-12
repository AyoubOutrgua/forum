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
	imagePath := ""
	imageFile, handler, err := r.FormFile("choose-file")
	if err != nil {
		if err == http.ErrMissingFile {
			fmt.Println("No File Upload !!!!!!!!!!!!")
		} else {
			http.Error(w, "Status Bad Request 5", http.StatusBadRequest)
			return
		}
	} else {

		defer imageFile.Close()

		file, errCreate := os.Create("./static/upload/" + handler.Filename)
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
		imagePath = "/static/upload/" + handler.Filename
	}
	fmt.Println("IMAGE ::::::::::: ", imagePath)

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timeNow)
	query := "INSERT INTO posts (title, post, imageUrl, userId, creationDate) VALUES (?, ?, ?, ?, ?)"
	fmt.Println(categories)
	database.ExecuteData(query, title[0], description[0], imagePath, 1, timeNow)

	// query2 := "INSERT INTO users (userName, email, password) VALUES (?, ?, ?)"
	// name := "user1"
	// email := "user1@example.com"
	// pass := "123456"
	// database.ExecuteData(query2, name, email, pass)

	// query3 := "INSERT INTO users (userName, email, password) VALUES (?, ?, ?)"
	// name2 := "user2"
	// email2 := "user2@example.com"
	// pass2 := "123456"
	// database.ExecuteData(query3, name2, email2, pass2)

	http.Redirect(w, r, "/", 303)
}
