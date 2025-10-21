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

	http.HandleFunc("/login", middleware.RateLimitLoginMiddleware(manager, handlers.LoginHandler))
	http.HandleFunc("/register", middleware.RateLimitLoginMiddleware(registerLimiter, handlers.RegisterHandler))
	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", handlers.HanldlerShowHome)
	http.HandleFunc("/createpost", middleware.RateLimitPost(handlers.CreatePostHandler))
	http.HandleFunc("/logout", handlers.LogOutHandler)
}
