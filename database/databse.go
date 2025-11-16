package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"forum/tools"

	_ "github.com/mattn/go-sqlite3"
)

var DataBase *sql.DB

func InitDataBase() error {
	var err error
	DataBase, err = sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}

	schema, err := os.ReadFile("./db/schema.sql")
	if err != nil {
		log.Fatal("can't read schema", err)
	}

	_, err = DataBase.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	_, err = DataBase.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}

	_, err = DataBase.Exec(`INSERT INTO categories (category)
	VALUES 
	('Technology'),
	('Education'),
	('Sports'),
	('Movies'),
	('Gaming'),
	('Music'),
	('Health'),
	('Other')`)
	if err != nil {
		return err
	}
	fmt.Println("The database was created successfully.")
	return nil
}

func CloseDataBase() error {
	if DataBase != nil {
		return DataBase.Close()
	}
	return nil
}

func ExecuteData(query string, args ...interface{}) error {
	_, errExuc := DataBase.Exec(query, args...)
	if errExuc != nil {
		return errExuc
	}
	return nil
}

func SelectAllPosts(query string) ([]tools.Post, error) {
	rows, err := DataBase.Query(query)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return 0, err
	}
	return lastID, nil
}

func SelectPostCategories(query string, id int) ([]int, error) {
	rows, err := DataBase.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []int{}
	for rows.Next() {
		var cat int
		err := rows.Scan(&cat)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func SelectUserID(query string, cookieID string) (int, error) {
	var userID int
	err := DataBase.QueryRow(query, cookieID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func SelectAllComments(query string) (map[int][]tools.Comment, error) {
	rows, err := DataBase.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make(map[int][]tools.Comment)
	for rows.Next() {
		var c tools.Comment
		err := rows.Scan(&c.ID, &c.CommentText, &c.PostID, &c.UserID, &c.UserName, &c.CreationDate)
		if err != nil {
			return nil, err
		}
		comments[c.PostID] = append(comments[c.PostID], c)
	}
	return comments, nil
}

func SelectUserName(query string, userId int) (string, error) {
	userName := ""
	err := DataBase.QueryRow(query, userId).Scan(&userName)
	if err != nil {
		return "", err
	}
	return userName, nil
}
