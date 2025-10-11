package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateUser() {
	dataBase, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal()
	}
	defer dataBase.Close()
	usersTable := `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userName TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE ,
    password TEXT NOT NULL
);`
	_, err = dataBase.Exec(usersTable)
	if err != nil {
		log.Fatal(err)
	}
}
