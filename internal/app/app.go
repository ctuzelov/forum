package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"forum/internal/controller"
	"forum/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

const PORT = ":8080"

func Run() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	DB, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer DB.Close()
	err = repository.CreateTables(DB)
	if err != nil {
		errorLog.Fatal(err)
	}
	handler, err := controller.NewHandler(errorLog, DB)
	if err != nil {
		errorLog.Fatal(err)
	}
	srv := &http.Server{
		Addr:     PORT,
		ErrorLog: errorLog,
		Handler:  controller.Routes(handler),
	}
	log.Printf("listening on http://localhost" + PORT)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
