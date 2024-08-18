package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
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

func (s SchedulerStore) Add(t Task) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	lastInserted, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	return int(lastInserted), nil
}

func (s SchedulerStore) Get(t Task) (TaskResponse, error) {

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
