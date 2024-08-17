package main

import (
	"database/sql"
	"fmt"
	"gofinalproject/tests"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// appPath, err := os.Executable()
	appPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	// dbFile := filepath.Join(filepath.Dir(appPath), tests.DBFile)
	dbFile := filepath.Join(appPath, tests.DBFile)

	err = checkAndCreateDB(dbFile)
	if err != nil {
		log.Fatalf("Error while setting up database: %v", err)
	}

	//---
	db, err := sql.Open(DBDriver, dbFile)
	if err != nil {
		log.Fatalf("Error while opening database: %v", err)
	}
	defer db.Close()
	//---

	port := tests.GetPort("TODO_PORT")

	r := chi.NewRouter()

	http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.HandleFunc("/api/nextdate", nextDateHandler)
	r.Get("/api/nextdate", nextDateHandler)
	r.Post("/api/task", postTask(db))

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), r)
	if err != nil {
		log.Fatal(err)
	}

}
