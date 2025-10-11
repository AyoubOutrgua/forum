package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var dataBase *sql.DB

func InitDataBase() {
	dataBase, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}
	defer dataBase.Close()
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal("can't read schema", err)
	}
	_, err = dataBase.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}
	_, err = dataBase.Exec("PRAGMA foreign_keys = ON")
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
	if dataBase != nil {
		return dataBase.Close()
	}
	return nil
}
