package main

import (
	"database/sql"
	"fmt"
	"gofinalproject/tests"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DBDriver = "sqlite"
)

type SchedulerStore struct {
	db *sql.DB
}

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

	port := tests.GetPort("TODO_PORT")

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", nextDateHandler)

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}

}
