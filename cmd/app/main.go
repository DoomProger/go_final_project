package main

import (
	"database/sql"
	"fmt"
	"gofinalproject/config"
	"gofinalproject/internal/handlers"
	"gofinalproject/internal/repositories"
	"gofinalproject/internal/services"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

// "gofinalproject/internal/handlers"
// "gofinalproject/internal/repositories"
// "gofinalproject/internal/services"
// userRepo := repositories.NewUserRepository(db)
// userService := services.NewUserService(userRepo)
// userHandler := handlers.NewUserHandler(userService)

func main() {

	appPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(appPath, config.DBFile)

	err = repositories.CheckAndCreateDB(dbFile)
	if err != nil {
		log.Fatalf("Error while setting up database: %v", err)
	}

	db, err := sql.Open(config.DBDriver, config.DBFile)
	if err != nil {
		log.Fatalf("Error while opening database: %v", err)
	}
	defer db.Close()

	repo := repositories.NewSchedulerStore(db)
	taskService := services.NewTaskService(repo)
	taskHandler := handlers.NewTaskHandler(taskService)

	router := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("./web"))
	router.Handle("/*", http.StripPrefix("/", fileServer))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	router.Route("/api", func(r chi.Router) {
		r.Get("/nextdate", taskHandler.NextDate)

		r.Route("/task", func(rt chi.Router) {
			rt.Get("/", taskHandler.GetTask)
			rt.Post("/", taskHandler.AddTask)
			rt.Post("/done", taskHandler.DoneTask)
			rt.Delete("/", taskHandler.DeleteTask)
			rt.Put("/", taskHandler.UpdateTask)

		})

		r.Route("/tasks", func(rts chi.Router) {
			rts.Get("/", taskHandler.GetTasks)
		})

	})

	log.Println("Run on port:", config.Port)

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.Port), router)
	if err != nil {
		log.Fatal(err)
	}

}
