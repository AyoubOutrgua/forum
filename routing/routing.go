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

	http.HandleFunc("/loginAuth", middleware.RateLimitLoginMiddleware(manager, handlers.LoginHandler))
	http.HandleFunc("/register", handlers.Showregister)
	http.HandleFunc("/login", handlers.Showloginhandler)

	http.HandleFunc("/registerAuth",middleware.RateLimitLoginMiddleware(registerLimiter, handlers.RegisterHandler))
	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", handlers.HanldlerShowHome)
	http.HandleFunc("/createpost",middleware.Checksession( middleware.RateLimitPost(handlers.CreatePostHandler)))
	http.HandleFunc("/logout", handlers.LogOutHandler)
}
