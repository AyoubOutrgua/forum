package handlers

import (
	"database/sql"
	"net/http"
	"time"

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
		var expiredTime time.Time
		err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {

			loggedIn = true
			err = database.DataBase.QueryRow("SELECT id, dateexpired FROM users WHERE session = ?", cookie.Value).Scan(&userID, &expiredTime)
			if err == nil {
				if expiredTime.After(time.Now()) {
					loggedIn = true
				} else {
					_, err = database.DataBase.Exec(
						"UPDATE users SET session = NULL, dateexpired = NULL WHERE session = ?", cookie.Value)
						if err != nil{
							helpers.Errorhandler(w,"internal server error",http.StatusInternalServerError)
						}

					expiredCookie := &http.Cookie{
						Name:     "session",
						Value:    "",
						Path:     "/",
						MaxAge:   -1,
						Expires:  time.Now().Add(-1 * time.Hour),
						HttpOnly: true,
					}
					http.SetCookie(w, expiredCookie)
					loggedIn = false
				}
			} else if err == sql.ErrNoRows {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		}
	}

	
	dataIsLogin := tools.IsLogin{LoggedIn: loggedIn, UserID: userID}

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
	pageData.IsLogin = dataIsLogin
	pageData.ReactionStats = reactionStats
	pageData.UserReactions = userReactions
	pageData.Comment = comments
	pageData.ConnectUserName = connectUserName
	pageData.CommentReactionStats = commentReactionStats
	pageData.UserCommentReactions = userCommentReactions
	helpers.Render(w, "index.html", http.StatusOK, pageData)
}
