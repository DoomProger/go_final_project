package main

import (
	"database/sql"
	"time"

	// "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	// ID string `json:"id"`
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}

type SchedulerStore struct {
	db *sql.DB
}

func NewSchedulerStore(db *sql.DB) SchedulerStore {
	return SchedulerStore{db: db}
}

func (s SchedulerStore) AddTask(t Task) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return 0, err
	}

	lastInserted, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInserted), nil
}

func (s SchedulerStore) UpdateTask(t Task) error {

	_, err := s.db.Exec(
		"UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID))

	if err != nil {
		return err
	}

	return nil
}

func (s SchedulerStore) DeleteTask(id string) error {
	_, err := s.db.Exec(
		"DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

func (s SchedulerStore) GetTask(id string) (Task, error) {

	var task Task

	row := s.db.QueryRow(
		"SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s SchedulerStore) SearchTasks(search string) (TaskResponse, error) {

	var tasks TaskResponse
	var rows *sql.Rows

	parsedDate, err := time.Parse(dateFormatSearch, search)

	if err == nil {
		rows, err = s.db.Query(
			"SELECT * FROM scheduler WHERE date = :date",
			sql.Named("date", parsedDate.Format(dateFormat)))
		if err != nil {
			return TaskResponse{}, err
		}
		defer rows.Close()
	} else {
		rows, err = s.db.Query(
			"SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
			sql.Named("search", "%"+search+"%"),
			sql.Named("limit", limit50))
		if err != nil {
			return TaskResponse{}, err
		}
		defer rows.Close()

	}

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return TaskResponse{}, err
		}

		tasks.Tasks = append(tasks.Tasks, task)
	}

	if len(tasks.Tasks) == 0 {
		tasks.Tasks = []Task{}
		return tasks, nil
	}

	return tasks, nil
}

func (s SchedulerStore) GetTasks(t Task) (TaskResponse, error) {

	var tasks TaskResponse

	rows, err := s.db.Query(
		"SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit50))
	if err != nil {
		return TaskResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return TaskResponse{}, err
		}

		tasks.Tasks = append(tasks.Tasks, task)
	}

	if len(tasks.Tasks) == 0 {
		tasks.Tasks = []Task{}
		return tasks, nil
	}

	return tasks, nil
}
