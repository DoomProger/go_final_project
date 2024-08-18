package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorResponse := Response{
		Error: message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}

func postTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Println("Invalid request payload")
			json.NewEncoder(w).Encode(Response{Error: "Invalid request payload"})
			return
		}

		if task.Title == "" {
			log.Println("'title' is a required field")
			writeJSONError(w, http.StatusBadRequest, "'title' is a required field")
			return
		}

		now := time.Now()

		if task.Date == "" {
			task.Date = now.Format(dateFormat)
		}

		taskDate, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid 'date' format")
			return
		}

		if taskDate.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(dateFormat)
			} else {
				nextDate, err := NextDate(now, taskDate.Format(dateFormat), task.Repeat)
				if err != nil {
					writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Can not get next date; error %v", err))
					return
				}
				task.Date = nextDate
			}
		}

		store := NewSchedulerStore(db)

		res, err := store.Add(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{ID: strconv.Itoa(res)})
	}
}

func getTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		//TODO: search option

		store := NewSchedulerStore(db)
		res, err := store.Get(task)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}
