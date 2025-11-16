package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/database"
	"forum/helpers"
	"forum/tools"
)

func FilterByAuthorHandler(w http.ResponseWriter, r *http.Request) {
	// getting session from cookies
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		helpers.Errorhandler(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// checking if the session exists in db
	var userExists bool
	errSelect := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
	if errSelect != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if !userExists {
		helpers.Errorhandler(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// getting the userId from db
	var userID int
	errSelect = database.DataBase.QueryRow("SELECT id FROM users WHERE session = ?", cookie.Value).Scan(&userID)
	if errSelect == sql.ErrNoRows {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	} else if errSelect != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	q := fmt.Sprintf(`
        SELECT p.id, p.title, p.post AS description, COALESCE(p.imageUrl,'') AS imageUrl, u.userName, p.creationDate
        FROM posts p
        LEFT JOIN users u ON p.userId = u.id
        WHERE p.userId = %d
        ORDER BY p.creationDate DESC
    `, userID)

	posts, err := database.SelectAllPosts(q)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	helpers.GetPostCategories(w, posts)

	RenderPostsPage(w, posts, true, userID)
}

func FilterByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// getting categories from the query
	catStrs := r.URL.Query()["categories"]
	if len(catStrs) == 0 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// checking the validity of categories
	ids := []int{}
	for _, s := range catStrs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
			return
		}
		ids = append(ids, v)
	}
	if len(ids) == 0 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return

	}
	for _, id := range ids {
		if id < 1 || id > 8 {
			helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}
	inList := make([]string, len(ids))
	for i, id := range ids {
		inList[i] = strconv.Itoa(id)
	}
	in := strings.Join(inList, ",")
	q := fmt.Sprintf(`
        SELECT DISTINCT p.id, p.title, p.post AS description, COALESCE(p.imageUrl,'') AS imageUrl, u.userName, p.creationDate
        FROM posts p
        INNER JOIN postCategories pc ON p.id = pc.postId
        LEFT JOIN users u ON p.userId = u.id
        WHERE pc.categoryId IN (%s)
        ORDER BY p.creationDate DESC
    `, in)

	posts, err := database.SelectAllPosts(q)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	helpers.GetPostCategories(w, posts)
	loggedIn := false
	var userID int

	cookie, err := r.Cookie("session")
	if err == nil && cookie.Value != "" {
		err = database.DataBase.QueryRow(
		"SELECT id FROM users WHERE session = ?", cookie.Value,
		).Scan(&userID)

		if err == sql.ErrNoRows {

			helpers.Errorhandler(w, "Status Bad Request", http.StatusBadRequest)
			return
		} else if err != nil {
			helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			loggedIn = true
		}
	}
	RenderPostsPage(w, posts, loggedIn, userID)
}

func FilterByLikedHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		helpers.Errorhandler(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var userExists bool
	if err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists); err != nil || !userExists {
		helpers.Errorhandler(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := database.DataBase.QueryRow("SELECT id FROM users WHERE session = ?", cookie.Value).Scan(&userID); err != nil {
		helpers.Errorhandler(w, "Internal error", http.StatusInternalServerError)
		return
	}

	q := fmt.Sprintf(`
        SELECT p.id, p.title, p.post AS description, COALESCE(p.imageUrl,'') AS imageUrl, u.userName, p.creationDate
        FROM posts p
        INNER JOIN postReactions pr ON p.id = pr.postId
        LEFT JOIN users u ON p.userId = u.id
        WHERE pr.userId = %d AND pr.reaction = 1
        ORDER BY p.creationDate DESC
    `, userID)

	posts, err := database.SelectAllPosts(q)
	if err != nil {
		helpers.Errorhandler(w, "internal error", http.StatusInternalServerError)
		return
	}
	helpers.GetPostCategories(w, posts)
	RenderPostsPage(w, posts, true, userID)
}
//getting post informations and rendering the page
func RenderPostsPage(w http.ResponseWriter, posts []tools.Post, loggedIn bool, userID int) {
	reactionStats, err := helpers.GetAllReactionStats()
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userReactions, err := helpers.GetUserPostReactions(userID)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	comments := helpers.GetAllComments(w)
	connectUserName := helpers.GetConnectUserName(w, userID)
	commentReactionStats, err := helpers.GetAllCommentReactionStats() 
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return 
	}
	userCommentReactions, err := helpers.GetUserCommentReactions(userID) 
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return 
	}
	categories, err := database.SelectAllCategories("SELECT id, category FROM categories")
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	pageData := tools.PageData{
		Posts:                posts,
		Categories:           categories,
		IsLogin:              tools.IsLogin{LoggedIn: loggedIn, UserID: userID},
		ReactionStats:        reactionStats,
		UserReactions:        userReactions,
		Comment:              comments,
		ConnectUserName:      connectUserName,
		CommentReactionStats: commentReactionStats,
		UserCommentReactions: userCommentReactions,
	}

	helpers.Render(w, "index.html", http.StatusOK, pageData)
}
