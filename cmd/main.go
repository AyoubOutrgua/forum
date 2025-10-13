package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/handlers"
)

func main() {
	database.CreateTables()

	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", handlers.HanldlerShowHome)
	http.HandleFunc("/login", handlers.HanldlerShowLogin)
	http.HandleFunc("/register", handlers.HanldlerShowRegister)
	http.HandleFunc("/createpost", handlers.CreatePostHandler)
	fmt.Println("server is runing http://localhost:8089")
	http.ListenAndServe(":8089", nil)
}
