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

	var lastID int
	err = database.QueryRow(query).Scan(&lastID)
	if err != nil {
		panic(err)
	}
	return lastID, nil
}

func SelectPostCategories(query string, id int) []int {
	database, err := sql.Open("sqlite3", "database/forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}
	defer database.Close()

	rows, err := database.Query(query, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	categories := []int{}
	for rows.Next() {
		var cat int
		err := rows.Scan(&cat)
		if err != nil {
			panic(err)
		}
		categories = append(categories, cat)
	}
	return categories
}

func SelectLastDates(query string, id int) []string {
	database, err := sql.Open("sqlite3", "database/forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}
	defer database.Close()

	var dates []string
	rows, err := database.Query(query, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var date string
		err := rows.Scan(&date)
		if err != nil {
			panic(err)
		}
		dates = append(dates, date)
	}
	return dates
}
