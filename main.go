package main

import (
	"fmt"
	"time"
)

type Task struct {
	Id          int
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type Status string

type TaskList struct {
	Tasks []*Task
}

const (
	taskString          = "ID: %d, Description: %s, Status: %s, Created at: %s, Updated at: %s\n"
	TODO         Status = Status("To do")
	IN_PROGRESS  Status = Status("In progress")
	DONE         Status = Status("Done")
	noTaskString        = "No tasks found"
)

func main() {
	tasks := NewTaskList()

	tasks.AddTask(&Task{
		1,
		"First Task",
		IN_PROGRESS,
		time.Now(),
		nil,
	})
	tasks.AddTask(&Task{
		2,
		"Second Task",
		IN_PROGRESS,
		time.Now(),
		nil,
	})

	tasks.AddTask(&Task{
		3,
		"Third Task",
		DONE,
		time.Now(),
		nil,
	})
	tasks.RemoveTask(2)

	updateTime := time.Now().AddDate(0, 0, 1)
	tasks.AddTask(&Task{
		2,
		"Second Task",
		"In progress",
		time.Now(),
		&updateTime,
	})

	tasks.PrintAll()
}

func NewTaskList() *TaskList {
	return &TaskList{}
}

func (tl *TaskList) AddTask(task *Task) {
	tl.Tasks = append(tl.Tasks, task)
}

func (tl *TaskList) RemoveTask(id int) (*Task, error) {
	for i, v := range tl.Tasks {
		if v.Id == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return v, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

func (tl *TaskList) UpdateTask(id int, description string) error {
	for _, v := range tl.Tasks {
		if v.Id == id {
			updatedTime := time.Now()
			v.UpdatedAt = &updatedTime
			v.Description = description
			return nil
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

func (tl *TaskList) PrintAll() {
	tl.printHasNoTasks()

	for _, task := range tl.Tasks {
		task.printTask()
	}

	fmt.Printf("--------------- Total Tasks: %d ---------------\n", len(tl.Tasks))
}

func (tl *TaskList) PrintDone() {
	tl.printHasNoTasks()

	for _, task := range tl.Tasks {
		if task.Status == DONE {
			task.printTask()
		}
	}

}

func (tl *TaskList) PrintInProgress() {
	tl.printHasNoTasks()

	for _, task := range tl.Tasks {
		if task.Status == IN_PROGRESS {
			task.printTask()
		}
	}
}

func (tl *TaskList) printHasNoTasks() {
	if len(tl.Tasks) == 0 {
		fmt.Println(noTaskString)
	}
}

func (tl *TaskList) MarkInProgress(id int) {
	for _, task := range tl.Tasks {
		if task.Id == id {
			task.MarkAs(IN_PROGRESS)
		}
	}
}

func (tl *TaskList) MarkDone(id int) {
	for _, task := range tl.Tasks {
		if task.Id == id {
			task.MarkAs(DONE)
		}
	}
}

func (t *Task) printTask() {
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
