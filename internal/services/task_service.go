package services

import (
	"gofinalproject/internal/repositories"
	"gofinalproject/pkg/models"
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

func (ts *TaskService) SearchTasks(search string) (*models.TasksResponse, error) {
	return ts.repo.SearchTasks(search)
}

func (ts *TaskService) GetTasks() (*models.TasksResponse, error) {
	return ts.repo.GetTasks()
}
