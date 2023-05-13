package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "forum.db")
	count := 0
	query := `drop table reactions`
	db.Exec(query)
	fmt.Println(count)
}
