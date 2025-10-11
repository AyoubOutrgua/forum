package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTables() {
	database, errOpen := sql.Open("sqlite3", "forum.db")
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
	database, err := sql.Open("sqlite3", "forum.db")
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

func SelectData(query string) *sql.Rows {
	database, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal("can't open/create forum.db ", err)
	}
	defer database.Close()

	rows, err := database.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	return rows
}
