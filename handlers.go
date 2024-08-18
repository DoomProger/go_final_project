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
			writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
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
					writeJSONError(w, http.StatusBadRequest, err.Error())
					return
				}
				task.Date = nextDate
			}
		}

		store := NewSchedulerStore(db)

		res, err := store.AddTask(task)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{ID: strconv.Itoa(res)})
	}
}

func postTaskDone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		log.Println("postTaskDone query from http srv:", query)
		id := query.Get("id")
		if id == "" {
			writeJSONError(w, http.StatusNotFound, "Identifier not specified")
			return
		}

		store := NewSchedulerStore(db)
		task, err := store.GetTask(id)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Println("Got task", task)

		if task.Repeat == "" {
			err := store.DeleteTask(task.ID)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{}`))
			return
		}

		taskDate, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid 'date' format")
			return
		}

		now := time.Now()
		nextDate, err := NextDate(now, taskDate.Format(dateFormat), task.Repeat)
		if err != nil {
			log.Printf("Can not get next date; error %v\n", err)
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		task.Date = nextDate

		err = store.UpdateTask(task)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	}
}

func UpdateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if task.Title == "" {
			log.Println("'title' is a required field")
			writeJSONError(w, http.StatusBadRequest, "'title' is a required field")
			return
		}

		if task.Date == "" {
			log.Println("'date' is a required field")
			writeJSONError(w, http.StatusBadRequest, "'date' is a required field")
			return
		}

		if task.ID == "" {
			log.Println("'id' is a required field")
			writeJSONError(w, http.StatusBadRequest, "'title' is a required field")
			return
		}

		id, err := strconv.ParseInt(task.ID, 10, 64)
		if err != nil {
			log.Println("'id' must be int")
			writeJSONError(w, http.StatusBadRequest, "'id' must be int")
			return
		}
		if id == 0 {
			writeJSONError(w, http.StatusBadRequest, "'id' should not be: 0")
			return
		}

		now := time.Now()

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
					log.Printf("Can not get next date; error %v\n", err)
					writeJSONError(w, http.StatusBadRequest, err.Error())
					return
				}
				task.Date = nextDate
			}
		}

		store := NewSchedulerStore(db)

		_, err = store.GetTask(task.ID)
		if err != nil {
			log.Println("task not found", task.ID)
			writeJSONError(w, http.StatusBadRequest, "task not found")
			return
		}

		err = store.UpdateTask(task)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}

func getTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			writeJSONError(w, http.StatusNotFound, "Identifier not specified")
			return
		}

		store := NewSchedulerStore(db)
		task, err := store.GetTask(id)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)

	}
}

func getTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		//TODO: search option

		store := NewSchedulerStore(db)
		res, err := store.GetTasks(task)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
