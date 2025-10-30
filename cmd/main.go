package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/routing"
)

func main() {
	database.InitDataBase()
	defer database.CloseDataBase()
	routing.Routing()
	fmt.Println("server is runing http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Error !")
	}
}
