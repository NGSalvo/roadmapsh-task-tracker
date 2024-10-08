package stores

import (
	"encoding/json"
	"fmt"
	"os"
	"task-tracker/models"
	"time"
)

type JsonTaskStore struct {
	Tasks        []*models.Task
	JsonFileName string
}

func NewJsonTaskStore(jsonFileName string) *JsonTaskStore {
	fileExistAndCreate(jsonFileName, &[]models.Task{})

	return &JsonTaskStore{
		Tasks:        []*models.Task{},
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

func (j *JsonTaskStore) AddTask(task *models.Task) (*models.Task, error) {
	id, err := j.assignId()

	if err != nil {
		return nil, err
	}

	task.Id = id
	task.CreatedAt = time.Now()
	task.Status = models.TODO

	j.Tasks = append(j.Tasks, task)

	err = j.saveToFile()

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (j *JsonTaskStore) RemoveTask(id int) (*models.Task, error) {
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
	err := j.loadFromFile()

	if err != nil {
		return err
	}

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
		fmt.Println("No tasks found")
	}
}

func (j *JsonTaskStore) PrintAll() error {
	err := j.loadFromFile()

	if err != nil {
		return err
	}

	j.printHasNoTasks()

	for _, task := range j.Tasks {
		task.PrintTask()
	}
	fmt.Printf("--------------- Total Tasks: %d ---------------\n", len(j.Tasks))
	return nil
}

func (j *JsonTaskStore) PrintTodo() error {
	err := j.loadFromFile()

	if err != nil {
		return err
	}

	j.printHasNoTasks()

	filterdTasks := []*models.Task{}

	for _, task := range j.Tasks {
		if task.Status == models.TODO {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(models.NoTaskString)
		return nil
	}

	for _, task := range filterdTasks {
		task.PrintTask()
	}

	return nil

}

func (j *JsonTaskStore) PrintDone() error {
	err := j.loadFromFile()

	if err != nil {
		return err
	}

	j.printHasNoTasks()

	filterdTasks := []*models.Task{}

	for _, task := range j.Tasks {
		if task.Status == models.DONE {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(models.NoTaskString)
		return nil
	}

	for _, task := range filterdTasks {
		task.PrintTask()
	}

	return nil
}

func (j *JsonTaskStore) PrintInProgress() error {
	err := j.loadFromFile()

	if err != nil {
		return err
	}

	j.printHasNoTasks()

	filterdTasks := []*models.Task{}

	for _, task := range j.Tasks {
		if task.Status == models.IN_PROGRESS {
			filterdTasks = append(filterdTasks, task)
		}
	}

	if len(filterdTasks) == 0 {
		fmt.Println(models.NoTaskString)
		return nil
	}

	for _, task := range filterdTasks {
		task.PrintTask()
	}
	return nil
}

func (j *JsonTaskStore) MarkInProgress(id int) error {
	for _, task := range j.Tasks {
		if task.Id == id {
			task.MarkAs(models.IN_PROGRESS)

			err := j.saveToFile()
			if err != nil {
				return err
			}

			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (j *JsonTaskStore) MarkDone(id int) error {
	for _, task := range j.Tasks {
		if task.Id == id {
			task.MarkAs(models.DONE)

			err := j.saveToFile()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (j *JsonTaskStore) assignId() (int, error) {
	err := j.loadFromFile()

	if err != nil {
		return 0, err
	}

	if len(j.Tasks) == 0 {
		return 1, nil
	}

	currentMax := 1
	for _, task := range j.Tasks {
		if task.Id > currentMax {
			currentMax = task.Id
		}
	}
	return currentMax + 1, nil
}

func fileExistAndCreate(jsonFileName string, model any) {
	if model == nil {
		model = &[]models.Task{}
	}
	if _, err := os.Stat(jsonFileName); os.IsNotExist(err) {
		file, err := os.Create(jsonFileName)

		if err != nil {
			panic("Error creating file: " + err.Error())
		}

		content, err := json.MarshalIndent(model, "", " ")

		if err != nil {
			panic("Error creating file: " + err.Error())
		}

		err = os.WriteFile(jsonFileName, content, 0644)

		if err != nil {
			panic("Error creating file: " + err.Error())
		}

		file.Close()
	}
}
