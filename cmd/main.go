package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/middleware"
	routego "forum/route.go"
)

func main() {
	database.InitDataBase()
	defer database.CloseDataBase()
	handlers.Db = database.DataBase
	middleware.Db = database.DataBase
	routego.Routing()
	if handlers.Db == nil {
		log.Fatal("Database not initialized!")
	}

	fmt.Println("server is runing http://localhost:8085/")
	http.ListenAndServe(":8085", nil)
}
