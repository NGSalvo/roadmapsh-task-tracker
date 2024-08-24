package main

import (
	"fmt"
)

type Task struct {
	Id          int
	Description string
	Status      Status
}

type Status string

type TaskList struct {
	Tasks []*Task
}

const (
	taskString          = "ID: %d, Description: %s, Status: %s\n"
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
	})
	tasks.AddTask(&Task{
		2,
		"Second Task",
		IN_PROGRESS,
	})

	tasks.AddTask(&Task{
		3,
		"Third Task",
		DONE,
	})
	tasks.RemoveTask(2)

	tasks.AddTask(&Task{
		2,
		"Second Task",
		"In progress",
	})

	tasks.PrintAll()
}

func NewTaskList() *TaskList {
	return &TaskList{}
}

func (tl *TaskList) AddTask(task *Task) {
	tl.Tasks = append(tl.Tasks, task)
}

func (tl *TaskList) RemoveTask(id int) int {
	for i, v := range tl.Tasks {
		if v.Id == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return id
		}
	}
	return -1
}

func (tl *TaskList) UpdateTask(id int, description string) {
	for _, v := range tl.Tasks {
		if v.Id == id {
			v.Description = description
			return
		}
	}
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
	fmt.Printf(taskString, t.Id, t.Description, t.Status)
}

func (t *Task) MarkAs(status Status) {
	t.Status = status
}

func (s Status) String() string {
	return string(s)
}
