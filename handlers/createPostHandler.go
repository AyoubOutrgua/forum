package handlers

import (
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"forum/database"
	"forum/helpers"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	errParse := r.ParseMultipartForm(10 << 20)
	if errParse != nil {
		helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
		return
	}

	cookie, errSession := r.Cookie("session")
	if errSession != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	cookieID := cookie.Value

	userID := helpers.GetUserID(cookieID)

	title, ok := r.PostForm["title"]
	if !ok {
		helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	description, ok := r.PostForm["description"]
	if !ok {
		helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	if len(title[0]) == 0 || len(description[0]) == 0 || len(title[0]) > 100 || len(description[0]) > 1000 {
		helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	categories, ok := r.PostForm["categories"]

	if !ok {
		helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	if len(categories) == 0 {
		helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	categoriesID := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	for _, catsID := range categories {
		if !slices.Contains(categoriesID, catsID) {
			helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
			return
		}
	}

	imagePath := ""
	imageFile, handler, err := r.FormFile("choose-file")
	if err != nil {
		if err != http.ErrMissingFile {
			helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
			return
		}
	} else {

		defer imageFile.Close()

		if !helpers.IsImage(handler.Filename) {
			helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
			return
		}

		const maxSize = 2 * 1024 * 1024
		if handler.Size > maxSize {
			helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
			return
		}

		file, errCreate := os.Create("./static/upload/" + handler.Filename)
		if errCreate != nil {
			helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
			return
		}
		defer file.Close()

		_, errCopy := io.Copy(file, imageFile)
		if errCopy != nil {
			helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
			return
		}
		imagePath = "/static/upload/" + handler.Filename
	}

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	queryInsertPost := "INSERT INTO posts (title, post, imageUrl, userId, creationDate) VALUES (?, ?, ?, ?, ?)"
	errEx := database.ExecuteData(queryInsertPost, title[0], description[0], imagePath, userID, timeNow)
	if errEx != nil {
		helpers.Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
	}

	lastPostID, err := database.SelectLastIdOfPosts("SELECT id FROM posts ORDER BY creationDate DESC LIMIT 1;")
	if err != nil {
		helpers.Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
		return
	}

	queryInsertCategory := "INSERT INTO postCategories (postId, categoryId) VALUES (?, ?)"

	for _, catID := range categories {
		categoryID, err := strconv.Atoi(catID)
		if err != nil {
			helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
			return
		}
		errExec := database.ExecuteData(queryInsertCategory, lastPostID, categoryID)
		if errExec != nil {
			helpers.Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
