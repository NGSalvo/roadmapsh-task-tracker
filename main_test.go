package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Should add a task to the list", func(t *testing.T) {
		task := &Task{1, "Test Task", IN_PROGRESS}
		taskList := NewTaskList()
		taskList.AddTask(task)

		asserts.Equal(len(taskList.Tasks), 1)
	})

	t.Run("Should remove the first element from the list", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(&Task{1, "Test Task", IN_PROGRESS})
		taskList.RemoveTask(1)

		asserts.Equal(len(taskList.Tasks), 0)
	})

	t.Run("Should remove the second element from the list", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(&Task{1, "Test Task", IN_PROGRESS})
		taskList.AddTask(&Task{2, "Test Task", IN_PROGRESS})
		taskList.RemoveTask(2)

		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Equal(taskList.Tasks[0].Id, 1)
	})

	t.Run("Should return -1 when the task is not found", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(&Task{1, "Test Task", IN_PROGRESS})
		v := taskList.RemoveTask(2)

		asserts.Equal(v, -1)
	})

	t.Run("Should print all tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(&Task{1, "Test Task", IN_PROGRESS})
		taskList.AddTask(&Task{2, "Test Task", IN_PROGRESS})

		expected := joinMessage(taskList)
		expected = expected + "\n--------------- Total Tasks: 2 ---------------\n"
		result := outputToString(taskList.PrintAll)

		asserts.Equal(len(taskList.Tasks), 2)
		asserts.Equal(expected, result)
	})

	t.Run("Should print all done tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(&Task{1, "Test Task", DONE})
		taskList.AddTask(&Task{2, "Test Task", IN_PROGRESS})
		taskList.AddTask(&Task{3, "Test Task", DONE})

		expected := joinMessageWithFilter(taskList, DONE) + "\n"
		result := outputToString(taskList.PrintDone)

		asserts.Equal(len(taskList.Tasks), 3)
		asserts.Equal(expected, result)
	})

	t.Run("Should print all in progress tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(&Task{1, "Test Task", IN_PROGRESS})
		taskList.AddTask(&Task{2, "Test Task", IN_PROGRESS})
		taskList.AddTask(&Task{3, "Test Task", DONE})

		expected := joinMessageWithFilter(taskList, IN_PROGRESS) + "\n"
		result := outputToString(taskList.PrintInProgress)

		asserts.Equal(len(taskList.Tasks), 3)
		asserts.Equal(expected, result)
	})

	t.Run("Should print has no tasks", func(t *testing.T) {
		taskList := NewTaskList()

		expected := "No tasks found" + "\n--------------- Total Tasks: 0 ---------------\n"
		result := outputToString(taskList.PrintAll)
		asserts.Equal(len(taskList.Tasks), 0)
		asserts.Equal(expected, result)
	})

	t.Run("Should not print when there are no done tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(&Task{1, "Test Task", IN_PROGRESS})
		taskList.AddTask(&Task{2, "Test Task", IN_PROGRESS})

		result := outputToString(taskList.PrintDone)
		asserts.Equal(len(taskList.Tasks), 2)
		asserts.Empty(result)
	})

	t.Run("Should not print when there are no in progress tasks", func(t *testing.T) {
		taskList := NewTaskList()
		taskList.AddTask(&Task{3, "Test Task", DONE})

		result := outputToString(taskList.PrintInProgress)
		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Empty(result)
	})

}

func joinMessage(tasks *TaskList) string {
	message := []string{}
	for _, task := range tasks.Tasks {
		message = append(message, fmt.Sprintf("ID: %d, Description: %s, Status: %s", task.Id, task.Description, task.Status))
	}
	return strings.Join(message, "\n")
}

func joinMessageWithFilter(tasks *TaskList, filter Status) string {
	message := []string{}
	for _, task := range tasks.Tasks {
		if filter != "" && task.Status != filter {
			continue
		}
		message = append(message, fmt.Sprintf("ID: %d, Description: %s, Status: %s", task.Id, task.Description, task.Status))
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
