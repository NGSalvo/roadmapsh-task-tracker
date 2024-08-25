package stores

import (
	"fmt"
	"task-tracker/models"
	"time"
)

type InMemoryTaskStore struct {
	Tasks []*models.Task
}

func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{}
}

func (tl *InMemoryTaskStore) AddTask(task *models.Task) (*models.Task, error) {
	tl.Tasks = append(tl.Tasks, task)
	return task, nil
}

func (tl *InMemoryTaskStore) RemoveTask(id int) (*models.Task, error) {
	for i, v := range tl.Tasks {
		if v.Id == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return v, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

func (tl *InMemoryTaskStore) UpdateTask(id int, description string) error {
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

func (tl *InMemoryTaskStore) PrintAll() {
	tl.printHasNoTasks()

	for _, task := range tl.Tasks {
		task.PrintTask()
	}

	fmt.Printf("--------------- Total Tasks: %d ---------------\n", len(tl.Tasks))
}

func (tl *InMemoryTaskStore) PrintDone() {
	tl.printHasNoTasks()

	filterdTasks := []*models.Task{}

	for _, task := range tl.Tasks {
		if task.Status == models.DONE {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(models.NoTaskString)
		return
	}

	for _, task := range filterdTasks {
		task.PrintTask()
	}
}

func (tl *InMemoryTaskStore) PrintInProgress() {
	tl.printHasNoTasks()

	filterdTasks := []*models.Task{}

	for _, task := range tl.Tasks {
		if task.Status == models.IN_PROGRESS {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(models.NoTaskString)
		return
	}

	for _, task := range filterdTasks {
		task.PrintTask()
	}
}

func (tl *InMemoryTaskStore) printHasNoTasks() {
	if len(tl.Tasks) == 0 {
		fmt.Println(models.NoTaskString)
	}
}

func (tl *InMemoryTaskStore) MarkInProgress(id int) {
	for _, task := range tl.Tasks {
		if task.Id == id {
			task.MarkAs(models.IN_PROGRESS)
		}
	}
}

func (tl *InMemoryTaskStore) MarkDone(id int) {
	for _, task := range tl.Tasks {
		if task.Id == id {
			task.MarkAs(models.DONE)
		}
	}
}
