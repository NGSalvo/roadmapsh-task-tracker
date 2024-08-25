package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
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

type InMemoryTaskStore struct {
	Tasks []*Task
}

type JsonTaskStore struct {
	Tasks        []*Task
	JsonFileName string
}

const (
	taskString          = "ID: %d, Description: %s, Status: %s, Created at: %s, Updated at: %s\n"
	TODO         Status = Status("To do")
	IN_PROGRESS  Status = Status("In progress")
	DONE         Status = Status("Done")
	noTaskString        = "No tasks found"
)

func main() {
	tasks := NewInMemoryTaskStore()

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

	fmt.Println("\n\n--------- JSON TASKS ---------")
	t := NewJsonTaskStore("tasks.json")
	t.AddTask(&Task{
		1,
		"First Task",
		IN_PROGRESS,
		time.Now(),
		nil,
	})
	t.AddTask(&Task{
		2,
		"Second Task",
		IN_PROGRESS,
		time.Now(),
		nil,
	})
	t.RemoveTask(1)
	t.UpdateTask(2, "Second Task Updated")
	t.PrintAll()
	t.PrintDone()
	t.PrintInProgress()
}

func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{}
}

func (tl *InMemoryTaskStore) AddTask(task *Task) (*Task, error) {
	tl.Tasks = append(tl.Tasks, task)
	return task, nil
}

func (tl *InMemoryTaskStore) RemoveTask(id int) (*Task, error) {
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
		task.printTask()
	}

	fmt.Printf("--------------- Total Tasks: %d ---------------\n", len(tl.Tasks))
}

func (tl *InMemoryTaskStore) PrintDone() {
	tl.printHasNoTasks()

	filterdTasks := []*Task{}

	for _, task := range tl.Tasks {
		if task.Status == DONE {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(noTaskString)
		return
	}

	for _, task := range filterdTasks {
		task.printTask()
	}
}

func (tl *InMemoryTaskStore) PrintInProgress() {
	tl.printHasNoTasks()

	filterdTasks := []*Task{}

	for _, task := range tl.Tasks {
		if task.Status == IN_PROGRESS {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(noTaskString)
		return
	}

	for _, task := range filterdTasks {
		task.printTask()
	}
}

func (tl *InMemoryTaskStore) printHasNoTasks() {
	if len(tl.Tasks) == 0 {
		fmt.Println(noTaskString)
	}
}

func (tl *InMemoryTaskStore) MarkInProgress(id int) {
	for _, task := range tl.Tasks {
		if task.Id == id {
			task.MarkAs(IN_PROGRESS)
		}
	}
}

func (tl *InMemoryTaskStore) MarkDone(id int) {
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

func NewJsonTaskStore(jsonFileName string) *JsonTaskStore {
	return &JsonTaskStore{
		Tasks:        []*Task{},
		JsonFileName: jsonFileName,
	}
}

func (j *JsonTaskStore) saveToFile() error {
	file, err := json.MarshalIndent(j.Tasks, "", " ")

	if err != nil {
		return err
	}

	err = os.WriteFile(j.JsonFileName, file, 0644)

	if err != nil {
		return err
	}
	return nil
}

func (j *JsonTaskStore) loadFromFile() error {
	file, err := os.ReadFile(j.JsonFileName)

	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &j.Tasks)

	if err != nil {
		return err
	}
	return nil
}

func (j *JsonTaskStore) AddTask(task *Task) (*Task, error) {
	j.Tasks = append(j.Tasks, task)

	err := j.saveToFile()

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (j *JsonTaskStore) RemoveTask(id int) (*Task, error) {
	for i, v := range j.Tasks {
		if v.Id == id {
			j.Tasks = append(j.Tasks[:i], j.Tasks[i+1:]...)
			err := j.saveToFile()

			if err != nil {
				return nil, err
			}

			return v, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

func (j *JsonTaskStore) UpdateTask(id int, description string) error {
	for _, v := range j.Tasks {
		if v.Id == id {
			updatedTime := time.Now()
			v.UpdatedAt = &updatedTime
			v.Description = description
			err := j.saveToFile()

			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (j *JsonTaskStore) printHasNoTasks() {
	if len(j.Tasks) == 0 {
		fmt.Println(noTaskString)
	}
}

func (j *JsonTaskStore) PrintAll() error {
	err := j.loadFromFile()

	if err != nil {
		return err
	}

	j.printHasNoTasks()

	for _, task := range j.Tasks {
		task.printTask()
	}
	fmt.Printf("--------------- Total Tasks: %d ---------------\n", len(j.Tasks))
	return nil
}

func (j *JsonTaskStore) PrintDone() error {
	err := j.loadFromFile()

	if err != nil {
		return err
	}

	j.printHasNoTasks()

	filterdTasks := []*Task{}

	for _, task := range j.Tasks {
		if task.Status == DONE {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(noTaskString)
		return nil
	}

	for _, task := range filterdTasks {
		task.printTask()
	}

	return nil
}

func (j *JsonTaskStore) PrintInProgress() error {
	err := j.loadFromFile()

	if err != nil {
		return err
	}

	j.printHasNoTasks()

	filterdTasks := []*Task{}

	for _, task := range j.Tasks {
		if task.Status == IN_PROGRESS {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(noTaskString)
		return nil
	}

	for _, task := range filterdTasks {
		task.printTask()
	}
	return nil
}
