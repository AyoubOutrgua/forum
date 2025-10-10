package main

import (
	"fmt"
	"net/http"

	"forum/handlers"

)

func main() {
	handlers.DataBase()
	defer handlers.Db.Close()

	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", handlers.HanldlerShowHome)
	http.HandleFunc("/login", handlers.HanldlerShowLogin)
	http.HandleFunc("/register", handlers.HanldlerShowRegister)
	http.HandleFunc("/loginauth", handlers.LoginHandler)
	http.HandleFunc("/regtistartion", handlers.RegisterHandler)
	fmt.Println("server is runing http://localhost:8085/")
	http.ListenAndServe(":8085", nil)
}
