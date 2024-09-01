package repositories

import (
	"database/sql"
	"gofinalproject/config"
	"gofinalproject/pkg/models"
	"time"
)

type SchedulerStore struct {
	db *sql.DB
}

func NewSchedulerStore(db *sql.DB) *SchedulerStore {
	return &SchedulerStore{db: db}
}

func (s *SchedulerStore) AddTask(t models.Task) (int, error) {
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

func (s *SchedulerStore) UpdateTask(t *models.Task) error {

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

func (s *SchedulerStore) DeleteTask(id string) error {
	_, err := s.db.Exec(
		"DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

func (s *SchedulerStore) GetTask(id string) (*models.Task, error) {

	var task models.Task

	row := s.db.QueryRow(
		"SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return &models.Task{}, err
	}

	return &task, nil
}

func (s *SchedulerStore) SearchTasks(search string) ([]*models.Task, error) {

	var tasks []*models.Task
	var rows *sql.Rows
	var err error

	if search != "" {

		parsedDate, err := time.Parse(config.DateFormatSearch, search)

		if err == nil {

			rows, err = s.db.Query(
				"SELECT * FROM scheduler WHERE date = :date",
				sql.Named("date", parsedDate.Format(config.DateFormat)))
			if err != nil {
				return make([]*models.Task, 0), err
			}
			defer rows.Close()

		} else {

			rows, err = s.db.Query(
				"SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
				sql.Named("search", "%"+search+"%"),
				sql.Named("limit", config.Limit50))
			if err != nil {
				return make([]*models.Task, 0), err
			}
			defer rows.Close()
		}

	} else {

		rows, err = s.db.Query(
			"SELECT * FROM scheduler ORDER BY date LIMIT :limit",
			sql.Named("limit", config.Limit50))
		if err != nil {
			return make([]*models.Task, 0), err
		}
		defer rows.Close()

	}

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return make([]*models.Task, 0), err
		}

		tasks = append(tasks, &task)

	}

	if err = rows.Close(); err != nil {
		return make([]*models.Task, 0), err
	}

	if len(tasks) == 0 {
		return make([]*models.Task, 0), nil
	}

	return tasks, nil
}
