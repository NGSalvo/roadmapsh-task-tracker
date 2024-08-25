package models

import (
	"fmt"
	"time"
)

const (
	taskString          = "ID: %d, Description: %s, Status: %s, Created at: %s, Updated at: %s\n"
	TODO         Status = Status("To do")
	IN_PROGRESS  Status = Status("In progress")
	DONE         Status = Status("Done")
	NoTaskString        = "No tasks found"
)

type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type Status string

type TaskStore interface {
	AddTask(*Task) (*Task, error)
	RemoveTask(int) (*Task, error)
	UpdateTask(int, string) error
}

func (t *Task) PrintTask() {
	if t.UpdatedAt == nil {
		fmt.Printf(taskString, t.Id, t.Description, t.Status, t.CreatedAt.Format(time.DateOnly), "")
		return
	}
	fmt.Printf(taskString, t.Id, t.Description, t.Status, t.CreatedAt.Format(time.DateOnly), t.UpdatedAt.Format("02/01/2006"))
}

func (t *Task) MarkAs(status Status) {
	t.Status = status
}

func (s Status) String() string {
	return string(s)
}
