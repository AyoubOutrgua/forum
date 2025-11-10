package routing

import (
	"net/http"
	"time"

	"forum/handlers"
	"forum/middleware"
)

func Routing() {
	manager := middleware.NewRateLimiterManager(10, 1*time.Minute)
	registerLimiter := middleware.NewRateLimiterManager(10, 1*time.Minute)
	http.HandleFunc("/createcomment", middleware.RateLimitComments(handlers.CreateCommentHandler))
	http.HandleFunc("/login", middleware.RateLimitLoginMiddleware(manager, handlers.LoginHandler))
	http.HandleFunc("/register", middleware.RateLimitLoginMiddleware(registerLimiter, handlers.RegisterHandler))
	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", handlers.HanldlerShowHome)
	http.HandleFunc("/createpost", middleware.RateLimitPost(handlers.CreatePostHandler))
	http.HandleFunc("/logout", handlers.LogOutHandler)
	http.HandleFunc("/reaction", handlers.ReactionHandler)
	http.HandleFunc("/comment-reaction", handlers.CommentReactionHandler)
	http.HandleFunc("/filter/author", handlers.FilterByAuthorHandler)
	http.HandleFunc("/filter/category", handlers.FilterByCategoryHandler)
	http.HandleFunc("/filter/liked", handlers.FilterByLikedHandler)
}
