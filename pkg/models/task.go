package models

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func convertModelsTasksToTasks(modelsTasks []*Task) []Task {
	var tasks []Task
	for _, mt := range modelsTasks {
		task := Task{
			ID:      mt.ID,
			Date:    mt.Date,
			Title:   mt.Title,
			Comment: mt.Comment,
			Repeat:  mt.Repeat,
		}
		tasks = append(tasks, task)
	}
	return tasks
}
