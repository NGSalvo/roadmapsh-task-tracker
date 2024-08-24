package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	asserts := assert.New(t)

	t.Run("✅ Should create a new task list", func(t *testing.T) {
		taskList := NewTaskList()

		asserts.Equal(len(taskList.Tasks), 0)
	})

	t.Run("✅ Should add a task to the list", func(t *testing.T) {
		task := createTask(1, IN_PROGRESS)
		taskList := NewTaskList()
		taskList.AddTask(task)

		asserts.Equal(len(taskList.Tasks), 1)
	})

	t.Run("✅ Should remove the first element from the list", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, IN_PROGRESS))
		taskList.RemoveTask(1)

		asserts.Equal(len(taskList.Tasks), 0)
	})

	t.Run("✅ Should remove the second element from the list and return the removed task", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, IN_PROGRESS))
		taskList.AddTask(createTask(2, IN_PROGRESS))
		task, err := taskList.RemoveTask(2)

		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Equal(taskList.Tasks[0].Id, 1)
		asserts.Equal(task.Id, 2)
		asserts.Nil(err)
	})

	t.Run("❌ Should return an error when removing a task that does not exist", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, IN_PROGRESS))
		task, err := taskList.RemoveTask(2)

		asserts.Equal(len(taskList.Tasks), 1)
		asserts.EqualError(err, "task with ID 2 not found")
		asserts.Nil(task)
	})

	t.Run("✅ Should update the description of a task", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, TODO))
		err := taskList.UpdateTask(1, "Updated Task")

		updatedAt := time.Now().Format(time.DateOnly)

		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Equal(taskList.Tasks[0].Description, "Updated Task")
		asserts.Equal(taskList.Tasks[0].UpdatedAt.Format(time.DateOnly), string(updatedAt))
		asserts.GreaterOrEqual(taskList.Tasks[0].UpdatedAt.Format(time.DateOnly), string(updatedAt))
		asserts.Nil(err)
	})

	t.Run("✅ Should update the description of the second task", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, IN_PROGRESS))
		taskList.AddTask(createTask(2, IN_PROGRESS))
		err := taskList.UpdateTask(2, "Updated Task")

		updatedAt := time.Now().Format(time.DateOnly)

		asserts.Equal(len(taskList.Tasks), 2)
		asserts.Equal(taskList.Tasks[1].Description, "Updated Task")
		asserts.Equal(taskList.Tasks[1].UpdatedAt.Format(time.DateOnly), string(updatedAt))
		asserts.GreaterOrEqual(taskList.Tasks[1].UpdatedAt.Format(time.DateOnly), string(updatedAt))
		asserts.Nil(err)
	})

	t.Run("❌ Should return an error when trying to update the description of a task that does not exist", func(t *testing.T) {
		taskList := NewTaskList()
		err := taskList.UpdateTask(1, "Updated Task")

		asserts.Nil(taskList.Tasks)
		asserts.EqualError(err, "task with ID 1 not found")
	})

	t.Run("✅ Should print all tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, IN_PROGRESS))
		taskList.AddTask(createTask(2, DONE))

		expected := joinMessage(taskList)
		expected = expected + "\n--------------- Total Tasks: 2 ---------------\n"
		result := outputToString(taskList.PrintAll)

		asserts.Equal(len(taskList.Tasks), 2)
		asserts.Equal(expected, result)
	})

	t.Run("✅ Should print all done tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, DONE))
		taskList.AddTask(createTask(2, IN_PROGRESS))
		taskList.AddTask(createTask(3, DONE))

		expected := joinMessageWithFilter(taskList, DONE) + "\n"
		result := outputToString(taskList.PrintDone)

		asserts.Equal(len(taskList.Tasks), 3)
		asserts.Equal(expected, result)
	})

	t.Run("✅ Should print all in progress tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, IN_PROGRESS))
		taskList.AddTask(createTask(2, IN_PROGRESS))
		taskList.AddTask(createTask(3, DONE))

		expected := joinMessageWithFilter(taskList, IN_PROGRESS) + "\n"
		result := outputToString(taskList.PrintInProgress)

		asserts.Equal(len(taskList.Tasks), 3)
		asserts.Equal(expected, result)
	})

	t.Run("✅ Should print has no tasks", func(t *testing.T) {
		taskList := NewTaskList()

		expected := "No tasks found" + "\n--------------- Total Tasks: 0 ---------------\n"
		result := outputToString(taskList.PrintAll)
		asserts.Equal(len(taskList.Tasks), 0)
		asserts.Equal(expected, result)
	})

	t.Run("❌ Should not print when there are no done tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, IN_PROGRESS))
		taskList.AddTask(createTask(2, IN_PROGRESS))

		result := outputToString(taskList.PrintDone)
		asserts.Equal(len(taskList.Tasks), 2)
		asserts.Empty(result)
	})

	t.Run("❌ Should not print when there are no in progress tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(3, DONE))

		result := outputToString(taskList.PrintInProgress)
		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Empty(result)
	})

	t.Run("✅ Should mark a task as in progress", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, TODO))
		taskList.MarkInProgress(1)

		asserts.Equal(taskList.Tasks[0].Status, IN_PROGRESS)
	})

	t.Run("✅ Should mark a task as done", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(createTask(1, TODO))
		taskList.MarkDone(1)
	})
}

func joinMessage(tasks *TaskList) string {
	message := []string{}
	for _, task := range tasks.Tasks {
		if task.UpdatedAt == nil {
			message = append(message, fmt.Sprintf("ID: %d, Description: %s, Status: %s, Created at: %s, Updated at: %s", task.Id, task.Description, task.Status, task.CreatedAt.Format(time.DateOnly), ""))
			continue
		}
		message = append(message, fmt.Sprintf("ID: %d, Description: %s, Status: %s, Created at: %s, Updated at: %s", task.Id, task.Description, task.Status, task.CreatedAt.Format(time.DateOnly), task.UpdatedAt.Format("02/01/2006")))
	}
	return strings.Join(message, "\n")
}

func joinMessageWithFilter(tasks *TaskList, filter Status) string {
	message := []string{}
	for _, task := range tasks.Tasks {
		if filter != "" && task.Status != filter {
			continue
		}
		if task.UpdatedAt == nil {
			message = append(message, fmt.Sprintf("ID: %d, Description: %s, Status: %s, Created at: %s, Updated at: %s", task.Id, task.Description, task.Status, task.CreatedAt.Format(time.DateOnly), ""))
			continue
		}
		message = append(message, fmt.Sprintf("ID: %d, Description: %s, Status: %s, Created at: %s, Updated at: %s", task.Id, task.Description, task.Status, task.CreatedAt.Format(time.DateOnly), task.UpdatedAt.Format("02/01/2006")))
	}
	return strings.Join(message, "\n")
}

func outputToString(callback func()) string {

	// Create a pipe to capture the output
	r, w, _ := os.Pipe()

	// Save the original stdout
	oldStdout := os.Stdout

	// Assign the write end of the pipe to stdout
	os.Stdout = w

	// Ensure that stdout is restored after the test
	defer func() {
		os.Stdout = oldStdout
		w.Close()
	}()

	// Call the function or code that prints to stdout
	callback()

	// Close the write end of the pipe to signal EOF
	w.Close()

	// Read the output from the read end of the pipe
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func createTask(id int, status Status) *Task {
	return &Task{
		Id:          id,
		Description: fmt.Sprintf("Task %d", id),
		Status:      status,
		CreatedAt:   time.Date(2024, 8, 24, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   nil,
	}
}
