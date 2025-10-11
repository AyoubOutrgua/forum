package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateUser() {
	dataBase, err := sql.Open("sqlite3", "forum.db")
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
}
