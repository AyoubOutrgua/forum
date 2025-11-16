package handlers

import (
	"net/http"

	"forum/database"
	"forum/helpers"
	"forum/tools"
)

func HanldlerShowHome(w http.ResponseWriter, r *http.Request) {
	loggedIn := false
	var userID int

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

	dataIsLogin := tools.IsLogin{LoggedIn: loggedIn, UserID: userID}

	posts := helpers.GetAllPosts(w)
	categories := helpers.GetAllCategories(w)
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

	var pageData tools.PageData
	pageData.Posts = posts
	pageData.Categories = categories
	pageData.IsLogin = dataIsLogin
	pageData.ReactionStats = reactionStats
	pageData.UserReactions = userReactions
	pageData.Comment = comments
	pageData.ConnectUserName = connectUserName
	pageData.CommentReactionStats = commentReactionStats
	pageData.UserCommentReactions = userCommentReactions
	helpers.Render(w, "index.html", http.StatusOK, pageData)
}
