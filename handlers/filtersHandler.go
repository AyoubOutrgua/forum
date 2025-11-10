package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/database"
	"forum/helpers"
	"forum/tools"
)

func FilterByAuthorHandler(w http.ResponseWriter, r *http.Request) {
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
        LEFT JOIN users u ON p.userId = u.id
        WHERE p.userId = %d
        ORDER BY p.creationDate DESC
    `, userID)

	posts, err := database.SelectAllPosts(q)
	if err != nil {
		helpers.Errorhandler(w, "internal error", http.StatusInternalServerError)
		return
	}
	helpers.GetPostCategories(w, posts)

	reactionStats := helpers.GetAllReactionStats(w)
	userReactions := helpers.GetUserPostReactions(w, userID)
	comments := helpers.GetAllComments(w)
	connectUserName := helpers.GetConnectUserName(w, userID)
	commentReactionStats := helpers.GetAllCommentReactionStats(w)
	userCommentReactions := helpers.GetUserCommentReactions(w, userID)
	categories, _ := database.SelectAllCategories("SELECT id, category FROM categories")

	pageData := tools.PageData{
		Posts:                posts,
		Categories:           categories,
		IsLogin:              tools.IsLogin{LoggedIn: true, UserID: userID},
		ReactionStats:        reactionStats,
		UserReactions:        userReactions,
		Comment:              comments,
		ConnectUserName:      connectUserName,
		CommentReactionStats: commentReactionStats,
		UserCommentReactions: userCommentReactions,
	}

	helpers.RenderPage(w, r, "index.html", pageData)
}

func FilterByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	catStrs := r.URL.Query()["categories"]
	if len(catStrs) == 0 {
		if single := r.URL.Query().Get("category_id"); single != "" {
			catStrs = []string{single}
		}
	}
	if len(catStrs) == 0 {
		helpers.Errorhandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

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
	inList := make([]string, len(ids))
	for i, id := range ids {
		inList[i] = strconv.Itoa(id)
	}
	in := strings.Join(inList, ",")

	q := fmt.Sprintf(`
        SELECT p.id, p.title, p.post AS description, COALESCE(p.imageUrl,'') AS imageUrl, u.userName, p.creationDate
        FROM posts p
        INNER JOIN postCategories pc ON p.id = pc.postId
        LEFT JOIN users u ON p.userId = u.id
        WHERE pc.categoryId IN (%s)
        GROUP BY p.id
        HAVING COUNT(DISTINCT pc.categoryId) = %d
        ORDER BY p.creationDate DESC
    `, in, len(ids))

	posts, err := database.SelectAllPosts(q)
	if err != nil {
		helpers.Errorhandler(w, "internal error", http.StatusInternalServerError)
		return
	}
	helpers.GetPostCategories(w, posts)

	loggedIn := false
	var userID int
	if cookie, err := r.Cookie("session"); err == nil && cookie.Value != "" {
		var exists bool
		if scanErr := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&exists); scanErr == nil && exists {
			loggedIn = true
			_ = database.DataBase.QueryRow("SELECT id FROM users WHERE session = ?", cookie.Value).Scan(&userID)
		}
	}

	reactionStats := helpers.GetAllReactionStats(w)
	userReactions := helpers.GetUserPostReactions(w, userID)
	comments := helpers.GetAllComments(w)
	connectUserName := helpers.GetConnectUserName(w, userID)
	commentReactionStats := helpers.GetAllCommentReactionStats(w)
	userCommentReactions := helpers.GetUserCommentReactions(w, userID)
	categories, _ := database.SelectAllCategories("SELECT id, category FROM categories")

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

	helpers.RenderPage(w, r, "index.html", pageData)
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

	reactionStats := helpers.GetAllReactionStats(w)
	userReactions := helpers.GetUserPostReactions(w, userID)
	comments := helpers.GetAllComments(w)
	connectUserName := helpers.GetConnectUserName(w, userID)
	commentReactionStats := helpers.GetAllCommentReactionStats(w)
	userCommentReactions := helpers.GetUserCommentReactions(w, userID)
	categories, _ := database.SelectAllCategories("SELECT id, category FROM categories")

	pageData := tools.PageData{
		Posts:                posts,
		Categories:           categories,
		IsLogin:              tools.IsLogin{LoggedIn: true, UserID: userID},
		ReactionStats:        reactionStats,
		UserReactions:        userReactions,
		Comment:              comments,
		ConnectUserName:      connectUserName,
		CommentReactionStats: commentReactionStats,
		UserCommentReactions: userCommentReactions,
	}

	helpers.RenderPage(w, r, "index.html", pageData)
}
