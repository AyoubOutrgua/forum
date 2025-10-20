package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDataBase() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}

	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal("can't read schema", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON")
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
	if DB != nil {
		return DB.Close()
	}
	return nil
}