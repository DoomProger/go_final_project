package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid 'now' date format: %v", err), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid 'date' format: %v", err), http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, date.Format("20060102"), repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// fmt.Println(nextDate)
	// log.Println(nextDate)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}

// func Index(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// do stuff with db here
// 		fmt.Fprintf(w, "Hello world!")
// 	}
// }

func postTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		// var buf bytes.Buffer

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, "'title' is a required field", http.StatusBadRequest)
			return
		}

		now := time.Now()
		nowStr := now.Format(dateFormat)

		if task.Date == "" {
			task.Date = nowStr
		} else {
			taskDate, err := time.Parse(dateFormat, task.Date)
			if err != nil {
				http.Error(w, "Invalid 'date' format", http.StatusBadRequest)
				return
			}

			if taskDate.Before(now) {
				if task.Repeat == "" {
					task.Date = nowStr
				} else {
					nextDate, err := NextDate(now, taskDate.Format(dateFormat), task.Repeat)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					task.Date = nextDate
				}
			}
		}

		store := NewSchedulerStore(db)

		res, err := store.Add(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Task was added", res)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}
}
