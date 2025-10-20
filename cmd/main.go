package main

import (
	"fmt"
	"net/http"
	"forum/database"
	routego "forum/route.go"
)

func main() {
	database.InitDataBase()
	defer database.CloseDataBase()
	routego.Routing()
	fmt.Println("server is runing http://localhost:8085")
	
	http.ListenAndServe(":8085", nil)
}
