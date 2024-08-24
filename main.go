package main

import (
	"database/sql"
	"fmt"
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

	db, err := sql.Open(DBDriver, dbFile)
	if err != nil {
		log.Fatalf("Error while opening database: %v", err)
	}
	defer db.Close()

	port := GetPort("TODO_PORT")

	router := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("./web"))
	router.Handle("/*", http.StripPrefix("/", fileServer))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	router.Route("/api", func(r chi.Router) {
		r.Get("/nextdate", nextDateHandler)

		r.Route("/task", func(rt chi.Router) {
			rt.Get("/", getTaskHandler(db))
			rt.Post("/", addTaskHandler(db))
			rt.Post("/done", doneTaskHandler(db))
			rt.Delete("/", deleteTaskHandler(db))
			rt.Put("/", updateTaskHandler(db))

		})

		r.Route("/tasks", func(rts chi.Router) {
			rts.Get("/", getTasksHandler(db))
		})
	})

	log.Println("Run on port:", port)

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router)
	if err != nil {
		log.Fatal(err)
	}

}
