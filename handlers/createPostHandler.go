package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"forum/database"
	"forum/middleware"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	checkerLimit := middleware.RateLimitPost()
	if checkerLimit {
		http.Error(w, "baraka postiti bzaf !!!!!!!!!!!!", http.StatusMethodNotAllowed)
		return
	}
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
	if len(categories) == 0 {
		http.Error(w, "Status Bad Request !!! madrti hta chi category !!!", http.StatusBadRequest)
		return
	}
	categoriesID := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	for _, catsID := range categories {
		if !slices.Contains(categoriesID, catsID) {
			http.Error(w, "Status Bad Request ----- makaynach had categories", http.StatusBadRequest)
			return
		}
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

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	queryInsertPost := "INSERT INTO posts (title, post, imageUrl, userId, creationDate) VALUES (?, ?, ?, ?, ?)"
	userID := 2
	database.ExecuteData(queryInsertPost, title[0], description[0], imagePath, userID, timeNow)

	lastPostID, err := database.SelectLastIdOfPosts("SELECT id FROM posts ORDER BY creationDate DESC LIMIT 1;")
	if err != nil {
		fmt.Println("ERROR : ", err)
		return
	}

	queryInsertCategory := "INSERT INTO postCategories (postId, categoryId) VALUES (?, ?)"

	for _, catID := range categories {
		categoryID, err := strconv.Atoi(catID)
		if err != nil {
			fmt.Println("ERROR ATOI : ", err)
			return
		}
		database.ExecuteData(queryInsertCategory, lastPostID, categoryID)
	}

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
