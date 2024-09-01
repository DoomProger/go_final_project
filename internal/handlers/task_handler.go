package handlers

import (
	"encoding/json"
	"fmt"
	"gofinalproject/config"
	"gofinalproject/internal/nextdate"
	"gofinalproject/internal/services"
	"gofinalproject/pkg/models"
	"net/http"
	"strconv"
	"time"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func (th *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

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
		task.Date = now.Format(config.DateFormat)
	}

	normNow := nextdate.NormalizeTime(now)
	taskDate, err := time.Parse(config.DateFormat, task.Date)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid 'date' format")
		return
	}

	if taskDate.Before(normNow) {
		if task.Repeat == "" {
			task.Date = now.Format(config.DateFormat)
		} else {
			nextDate, err := nextdate.NextDate(now, taskDate.Format(config.DateFormat), task.Repeat)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}
			task.Date = nextDate
		}
	}

	res, err := th.taskService.InsertTask(task)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	restask, err := th.taskService.GetTaskById(strconv.Itoa(res))

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{ID: restask.ID})
}

func (th *TaskHandler) NextDate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	now, err := time.Parse(config.DateFormat, nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid 'now' date format: %v", err), http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(nextDate))
}

func (th *TaskHandler) DoneTask(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	if id == "" {
		writeJSONError(w, http.StatusNotFound, "Identifier not specified")
		return
	}

	err := th.taskService.TaskDone(id)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}

func (th *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")
	if id == "" {
		writeJSONError(w, http.StatusNotFound, "Identifier not specified")
		return
	}

	task, err := th.taskService.GetTaskById(id)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = th.taskService.DeleteTaskById(task.ID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}

func (th *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task *models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ok, msg := isTaskFieldEmpty(*task)

	if !ok {
		writeJSONError(w, http.StatusBadRequest, msg)
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

	taskDate, err := time.Parse(config.DateFormat, task.Date)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid 'date' format")
		return
	}

	if taskDate.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(config.DateFormat)
		} else {
			nextDate, err := nextdate.NextDate(now, taskDate.Format(config.DateFormat), task.Repeat)

			if err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}
			task.Date = nextDate
		}
	}

	_, err = th.taskService.GetTaskById(task.ID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "task not found")
		return
	}

	err = th.taskService.ChangeTask(task)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}

func (th *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	id := query.Get("id")
	if id == "" {
		writeJSONError(w, http.StatusNotFound, "Identifier not specified")
		return
	}

	task, err := th.taskService.GetTaskById(id)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)

}

func (th *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")

	tasks, err := th.taskService.SearchTasks(search)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
