package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID           int
	Title        string
	Description  string
	ImageUrl     string
	UserName     string
	CreationDate string
}

func CreateTables() {
	database, errOpen := sql.Open("sqlite3", "database/forum.db")
	if errOpen != nil {
		log.Fatal("can't open/create forum.db ", errOpen)
	}
	defer database.Close()

	schema, errRead := os.ReadFile("database/schema.sql")
	if errRead != nil {
		log.Fatal("can't read schema ", errRead)
	}
	_, errExuc := database.Exec(string(schema))
	if errExuc != nil {
		log.Fatal(errExuc)
	}
	fmt.Println("Database Create OK")
}

func ExecuteData(query string, args ...interface{}) {
	database, err := sql.Open("sqlite3", "database/forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}
	defer database.Close()

	_, errExuc := database.Exec(query, args...)
	if errExuc != nil {
		log.Fatal(err)
	}
	fmt.Println("Execute OK!")
}

func SelectData(query string) ([]Post, error) {
	database, err := sql.Open("sqlite3", "database/forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}
	defer database.Close()

	rows, err := database.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ImageUrl, &p.UserName, &p.CreationDate)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
