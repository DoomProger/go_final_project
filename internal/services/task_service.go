package services

import (
	"gofinalproject/config"
	"gofinalproject/internal/nextdate"
	"gofinalproject/internal/repositories"
	"gofinalproject/pkg/models"
	"time"
)

type TaskService struct {
	repo *repositories.SchedulerStore
}

func NewTaskService(repo *repositories.SchedulerStore) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (ts *TaskService) InsertTask(t models.Task) (int, error) {
	return ts.repo.AddTask(t)
}

func (ts *TaskService) ChangeTask(t *models.Task) error {
	return ts.repo.UpdateTask(t)
}

func (ts *TaskService) DeleteTaskById(id string) error {
	return ts.repo.DeleteTask(id)
}

func (ts *TaskService) GetTaskById(id string) (*models.Task, error) {
	return ts.repo.GetTask(id)
}

func (ts *TaskService) TaskDone(id string) error {
	task, err := ts.GetTaskById(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		err := ts.DeleteTaskById(task.ID)
		if err != nil {
			return err
		}
		return nil
	}

	taskDate, err := time.Parse(config.DateFormat, task.Date)
	if err != nil {
		return err
	}

	now := time.Now()

	nextDate, err := nextdate.NextDate(now, taskDate.Format(config.DateFormat), task.Repeat)

	task.Date = nextDate

	err = ts.ChangeTask(task)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) SearchTasks(search string) ([]*models.Task, error) {
	return ts.repo.SearchTasks(search)
}
