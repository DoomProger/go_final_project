package repositories

import (
	"database/sql"
	"gofinalproject/config"
	"gofinalproject/pkg/models"
	"time"
)

type SchedulerStore struct {
	DB *sql.DB
}

func NewSchedulerStore(db *sql.DB) *SchedulerStore {
	return &SchedulerStore{DB: db}
}

func (s *SchedulerStore) AddTask(t models.Task) (int, error) {
	res, err := s.DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
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

func (s *SchedulerStore) UpdateTask(t *models.Task) error {

	_, err := s.DB.Exec(
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

func (s *SchedulerStore) DeleteTask(id string) error {
	_, err := s.DB.Exec(
		"DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

func (s *SchedulerStore) GetTask(id string) (*models.Task, error) {

	var task models.Task

	row := s.DB.QueryRow(
		"SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return &models.Task{}, err
	}

	return &task, nil
}

func (s *SchedulerStore) SearchTasks(search string) (*models.TaskResponse, error) {

	var tasks models.TaskResponse
	var rows *sql.Rows

	parsedDate, err := time.Parse(config.DateFormatSearch, search)

	if err == nil {
		rows, err = s.DB.Query(
			"SELECT * FROM scheduler WHERE date = :date",
			sql.Named("date", parsedDate.Format(config.DateFormat)))
		if err != nil {
			return &models.TaskResponse{}, err
		}
		defer rows.Close()
	} else {
		rows, err = s.DB.Query(
			"SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
			sql.Named("search", "%"+search+"%"),
			sql.Named("limit", config.Limit50))
		if err != nil {
			return &models.TaskResponse{}, err
		}
		defer rows.Close()
	}

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return &models.TaskResponse{}, err
		}

		tasks.Tasks = append(tasks.Tasks, task)
	}

	if len(tasks.Tasks) == 0 {
		tasks.Tasks = []models.Task{}
		return &tasks, nil
	}

	return &tasks, nil
}

func (s *SchedulerStore) GetTasks() (*models.TaskResponse, error) {

	var tasks models.TaskResponse

	rows, err := s.DB.Query(
		"SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", config.Limit50))
	if err != nil {
		return &models.TaskResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return &models.TaskResponse{}, err
		}

		tasks.Tasks = append(tasks.Tasks, task)
	}

	if len(tasks.Tasks) == 0 {
		tasks.Tasks = []models.Task{}
		return &tasks, nil
	}

	return &tasks, nil
}
