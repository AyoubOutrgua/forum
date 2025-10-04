package main

import (
	"fmt"
	"net/http"

	"forum/handlers"
)

func main() {
	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/login", handlers.HanldlerLogin)
	fmt.Println("server is runing http://localhost:8085/login")
	http.ListenAndServe(":8085",nil)
}
