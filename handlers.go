package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid 'now' date format: %v", err), http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}

func addTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if task.Title == "" {
			writeJSONError(w, http.StatusBadRequest, "'title' is a required field")
			return
		}

		now := time.Now()

		if task.Date == "" {
			task.Date = now.Format(dateFormat)
		}

		normNow := normalizeTime(now)
		taskDate, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid 'date' format")
			return
		}

		if taskDate.Before(normNow) {
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

		restask, err := store.GetTask(strconv.Itoa(res))
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{ID: restask.ID})
	}
}

func doneTaskHandler(db *sql.DB) http.HandlerFunc {
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}

func deleteTaskHandler(db *sql.DB) http.HandlerFunc {
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

		err = store.DeleteTask(task.ID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}

func updateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if task.Title == "" {
			writeJSONError(w, http.StatusBadRequest, "'title' is a required field")
			return
		}

		if task.Date == "" {
			writeJSONError(w, http.StatusBadRequest, "'date' is a required field")
			return
		}

		if task.ID == "" {
			writeJSONError(w, http.StatusBadRequest, "'title' is a required field")
			return
		}

		id, err := strconv.ParseInt(task.ID, 10, 64)
		if err != nil {
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
					writeJSONError(w, http.StatusBadRequest, err.Error())
					return
				}
				task.Date = nextDate
			}
		}

		store := NewSchedulerStore(db)

		_, err = store.GetTask(task.ID)
		if err != nil {
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

func getTaskHandler(db *sql.DB) http.HandlerFunc {
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

func getTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		query := r.URL.Query()
		search := query.Get("search")

		if search != "" {
			store := NewSchedulerStore(db)
			res, err := store.SearchTasks(search)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				// w.WriteHeader(http.StatusOK)
				// json.NewEncoder(w).Encode(res)
				// return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
			return
		}

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
