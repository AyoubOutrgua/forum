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
	http.HandleFunc("/createcomment", middleware.RateLimitMiddleware(rateLimiterComment, middleware.Checksession(handlers.CreateCommentHandler)))
	http.HandleFunc("/loginAuth", middleware.RateLimitMiddleware(rateLimiterLogin, middleware.CheckLogin(handlers.LoginHandler)))
	http.HandleFunc("/login", middleware.RateLimitMiddleware(rateLimiterRefresh, middleware.CheckLogin(handlers.Showloginhandler)))
	http.HandleFunc("/registerAuth", middleware.RateLimitMiddleware(rateLimiterRegister, middleware.CheckLogin(handlers.RegisterHandler)))
	http.HandleFunc("/register", middleware.RateLimitMiddleware(rateLimiterRefresh, middleware.CheckLogin(handlers.Showregister)))
	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.HanldlerShowHome))
	http.HandleFunc("/createpost", middleware.RateLimitMiddleware(rateLimiterPost, middleware.Checksession(handlers.CreatePostHandler)))
	http.HandleFunc("/logout", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.LogOutHandler))
	http.HandleFunc("/reaction", middleware.RateLimitMiddleware(rateLimiterReaction,middleware.Checksession(handlers.ReactionHandler)))
	http.HandleFunc("/comment-reaction",  middleware.RateLimitMiddleware(rateLimiterCommentReaction, middleware.Checksession(handlers.CommentReactionHandler)))
	http.HandleFunc("/filter/author",  middleware.RateLimitMiddleware(rateLimiterRefresh, middleware.Checksession(handlers.FilterByAuthorHandler)))
	http.HandleFunc("/filter/category", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.FilterByCategoryHandler))
	http.HandleFunc("/filter/liked", middleware.RateLimitMiddleware(rateLimiterRefresh, middleware.Checksession(handlers.FilterByLikedHandler)))
}
