package main

import (
	"fmt"
	"task-tracker/models"
	"task-tracker/stores"
	"time"
)

func main() {
	tasks := stores.NewInMemoryTaskStore()

	tasks.AddTask(&models.Task{
		1,
		"First Task",
		models.IN_PROGRESS,
		time.Now(),
		nil,
	})
	tasks.AddTask(&models.Task{
		2,
		"Second Task",
		models.IN_PROGRESS,
		time.Now(),
		nil,
	})

	tasks.AddTask(&models.Task{
		3,
		"Third Task",
		models.DONE,
		time.Now(),
		nil,
	})
	tasks.RemoveTask(2)

	updateTime := time.Now().AddDate(0, 0, 1)
	tasks.AddTask(&models.Task{
		2,
		"Second Task",
		"In progress",
		time.Now(),
		&updateTime,
	})

	tasks.PrintAll()

	fmt.Println("\n\n--------- JSON TASKS ---------")
	t := stores.NewJsonTaskStore("tasks.json")
	t.AddTask(&models.Task{
		1,
		"First Task",
		models.IN_PROGRESS,
		time.Now(),
		nil,
	})
	t.AddTask(&models.Task{
		2,
		"Second Task",
		models.IN_PROGRESS,
		time.Now(),
		nil,
	})
	t.RemoveTask(1)
	t.UpdateTask(2, "Second Task Updated")
	t.PrintAll()
	t.PrintDone()
	t.PrintInProgress()
	t.MarkAsDone(2)
	err := t.MarkAsDone(3)
	t.PrintAll()

	if err != nil {
		fmt.Println(err)
	}
}
