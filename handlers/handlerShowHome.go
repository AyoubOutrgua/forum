package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/helpers"
	"forum/tools"
)

func HanldlerShowHome(w http.ResponseWriter, r *http.Request) {
	loggedIn := false
	var userID int = 0

	if r.URL.Path != "/" {
		helpers.Errorhandler(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, errSession := r.Cookie("session")
	if errSession == nil && cookie.Value != "" {
		var userExists bool
		err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {
			loggedIn = true
			database.DataBase.QueryRow("SELECT id FROM users WHERE session = ?", cookie.Value).Scan(&userID)
		}
	}

	data := tools.IsLogin{LoggedIn: loggedIn, UserID: userID}

	temp, errPerse := template.ParseFiles("templates/index.html")
	if errPerse != nil {
		helpers.Errorhandler(w, "Status Not Found!", http.StatusNotFound)
		return
	}

	posts := helpers.GetAllPosts(w)
	categories := helpers.GetAllCategories(w)
	reactionStats := helpers.GetAllReactionStats(w)
	userReactions := helpers.GetUserPostReactions(w, userID)
	comments := helpers.GetAllComments(w)
	connectUserName := helpers.GetConnectUserName(w, userID)
	commentReactionStats := helpers.GetAllCommentReactionStats(w)
	userCommentReactions := helpers.GetUserCommentReactions(w, userID)

	var pageData tools.PageData
	pageData.Posts = posts
	pageData.Categories = categories
	pageData.IsLogin = data
	pageData.ReactionStats = reactionStats
	pageData.UserReactions = userReactions
	pageData.Comment = comments
	pageData.ConnectUserName = connectUserName
	pageData.CommentReactionStats = commentReactionStats
	pageData.UserCommentReactions = userCommentReactions

	errExec := temp.Execute(w, pageData)
	if errExec != nil {
		http.Error(w, "Status Internal Server Error", http.StatusInternalServerError)
	}
}
