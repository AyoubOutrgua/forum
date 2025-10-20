package database

import (
	"database/sql"
	"log"
	"os"

	"forum/tools"

	_ "github.com/mattn/go-sqlite3"
)

var DataBase *sql.DB

func InitDataBase() {
	var err error
	DataBase, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}

	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal("can't read schema", err)
	}

	_, err = DataBase.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	_, err = DataBase.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("can't enable foreign keys,", err)
	}
}

func InsertUser(dataBase *sql.DB, username, email, password string) error {
	query := `INSERT INTO users (userName, email, password) VALUES (?, ?, ?)`
	_, err := dataBase.Exec(query, username, email, password)
	return err
}

func CloseDataBase() error {
	if DataBase != nil {
		return DataBase.Close()
	}
	return nil
}

func ExecuteData(query string, args ...interface{}) {
	_, errExuc := DataBase.Exec(query, args...)
	if errExuc != nil {
		log.Fatal(errExuc)
	}
}

func SelectAllPosts(query string) ([]tools.Post, error) {
	rows, err := DataBase.Query(query)
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
	rows, err := DataBase.Query(query)
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
	var lastID int
	err := DataBase.QueryRow(query).Scan(&lastID)
	if err != nil {
		panic(err)
	}
	return lastID, nil
}

func SelectPostCategories(query string, id int) []int {
	rows, err := DataBase.Query(query, id)
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
	var dates []string
	rows, err := DataBase.Query(query, id)
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
