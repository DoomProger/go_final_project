package handlers

import "gofinalproject/pkg/models"

func isTaskFieldEmpty(task models.Task) (bool, string) {

	ok := true

	switch {
	case task.Title == "":
		return !ok, "'title' is a required field"
	case task.Date == "":
		return !ok, "'date' is a required field"
	case task.ID == "":
		return !ok, "'title' is a required field"

	default:
		return true, ""
	}

}
