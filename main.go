package main

import (
	"database/sql"
	"fmt"
	"gofinalproject/tests"
	"log"
	"net/http"
	"os"
	"path/filepath"

	// "github.com/go-chi/chi"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// appPath, err := os.Executable()
	appPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	// dbFile := filepath.Join(filepath.Dir(appPath), tests.DBFile)
	dbFile := filepath.Join(appPath, DBFile)

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

	// http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.HandleFunc("/api/tasks", getTasks(db))
	// r.Handle("/", http.FileServer(http.Dir("./web")))
	r.Get("/",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./web/index.html")
		})
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/web/*", http.StripPrefix("/web", fileServer))

	r.Get("/api/nextdate", nextDateHandler)

	r.Get("/api/task", getTask(db))
	r.Get("/api/tasks", getTasks(db))

	r.Post("/api/task/done", postTaskDone(db))
	r.Post("/api/task", postTask(db))

	r.Put("/api/task", UpdateTask(db))

	log.Println("Run on port:", port)

	// err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), r)
	if err != nil {
		log.Fatal(err)
	}

}
