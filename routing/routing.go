package routing

import (
	"net/http"
	"time"

	"forum/handlers"
	"forum/middleware"
)

func Routing() {
	rateLimiterLogin := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterRegister := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterPost := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterComment := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterRefresh := middleware.NewRateLimiterManager(200, time.Minute)
	rateLimiterReaction := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterCommentReaction := middleware.NewRateLimiterManager(20, time.Minute)

	// manager := middleware.NewRateLimiterManager(10, 1*time.Minute)
	// registerLimiter := middleware.NewRateLimiterManager(10, 1*time.Minute)
	http.HandleFunc("/createcomment", middleware.Checksession(middleware.RateLimitMiddleware(rateLimiterComment, handlers.CreateCommentHandler)))
	http.HandleFunc("/loginAuth", middleware.CheckLogin(middleware.RateLimitMiddleware(rateLimiterLogin, handlers.LoginHandler)))
	http.HandleFunc("/login", middleware.CheckLogin(handlers.Showloginhandler))

	http.HandleFunc("/registerAuth", middleware.CheckLogin(middleware.RateLimitMiddleware(rateLimiterRegister, handlers.RegisterHandler)))
	http.HandleFunc("/register", middleware.CheckLogin(handlers.Showregister))
	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.HanldlerShowHome))
	http.HandleFunc("/createpost", middleware.Checksession(middleware.RateLimitMiddleware(rateLimiterPost, handlers.CreatePostHandler)))
	http.HandleFunc("/logout", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.LogOutHandler))
	http.HandleFunc("/reaction", middleware.Checksession(middleware.RateLimitMiddleware(rateLimiterReaction, handlers.ReactionHandler)))
	http.HandleFunc("/comment-reaction", middleware.Checksession(middleware.RateLimitMiddleware(rateLimiterCommentReaction, handlers.CommentReactionHandler)))
	http.HandleFunc("/filter/author", middleware.Checksession(middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.FilterByAuthorHandler)))
	http.HandleFunc("/filter/category", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.FilterByCategoryHandler))
	http.HandleFunc("/filter/liked", middleware.Checksession(middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.FilterByLikedHandler)))
}
