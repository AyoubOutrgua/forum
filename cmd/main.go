package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/middleware"
)

func main() {
	database.CreateTables()

	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", handlers.HanldlerShowHome)
	http.HandleFunc("/login", handlers.HanldlerShowLogin)
	http.HandleFunc("/register", handlers.HanldlerShowRegister)
	http.HandleFunc("/createpost", middleware.RateLimitPost(handlers.CreatePostHandler))
	fmt.Println("server is runing http://localhost:8089")
	http.ListenAndServe(":8089", nil)
}
