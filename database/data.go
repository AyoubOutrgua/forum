package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"forum/tools"

	_ "github.com/mattn/go-sqlite3"
)

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
	database.Exec(`INSERT INTO categories (category)
	VALUES 
	('Technology'),
	('Education'),
	('Sports'),
	('Movies'),
	('Gaming'),
	('Music'),
	('Health'),
	('Food')`)
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

func SelectAllPosts(query string) ([]tools.Post, error) {
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

	var posts []tools.Post
	for rows.Next() {
		var p tools.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ImageUrl, &p.UserName, &p.CreationDate)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}



func SelectAllCategories(query string) ([]tools.Category, error) {
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

	var categories []tools.Category
	for rows.Next() {
		var p tools.Category
		err := rows.Scan(&p.ID, &p.Category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, p)
	}
	return categories, nil
}



func SelectLastIdOfPosts(query string) (int, error) {
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

	var p tools.Post
	for rows.Next() {
		err := rows.Scan(&p.ID)
		if err != nil {
			return 0, err
		}
	}
	return p.ID, nil
}