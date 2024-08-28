package stores

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"task-tracker/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() {
	os.Remove("test.json")
}

func TestJsonTaskStore(t *testing.T) {
	asserts := assert.New(t)

	t.Run("✅ Should create a new task list", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")

		asserts.Equal(len(taskList.Tasks), 0)
		asserts.FileExists("test.json")
	})

	t.Run("✅ Should add a task to the list", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		errorLoadingFile := taskList.loadFromFile()

		task := createTask2(1)
		createdTask, err := taskList.AddTask(task)

		asserts.FileExists("test.json")
		asserts.Nil(errorLoadingFile)
		asserts.Nil(err)
		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Contains(taskList.Tasks, createdTask)
		asserts.Equal(createdTask.Id, 1)
	})

	t.Run("✅ Should remove the first element from the list", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.RemoveTask(1)

		asserts.Equal(len(taskList.Tasks), 0)
	})

	t.Run("✅ Should remove the second element from the list and return the removed task", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.AddTask(createTask2(2))
		task, err := taskList.RemoveTask(2)

		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Equal(taskList.Tasks[0].Id, 1)
		asserts.Equal(task.Id, 2)
		asserts.Nil(err)
	})

	t.Run("❌ Should return an error when removing a task that does not exist", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		task, err := taskList.RemoveTask(2)

		asserts.Equal(len(taskList.Tasks), 1)
		asserts.EqualError(err, "task with ID 2 not found")
		asserts.Nil(task)
	})

	t.Run("✅ Should update the description of a task", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		err := taskList.UpdateTask(1, "Updated Task")

		updatedAt := time.Now().Format(time.DateOnly)

		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Equal(taskList.Tasks[0].Description, "Updated Task")
		asserts.Equal(taskList.Tasks[0].UpdatedAt.Format(time.DateOnly), string(updatedAt))
		asserts.GreaterOrEqual(taskList.Tasks[0].UpdatedAt.Format(time.DateOnly), string(updatedAt))
		asserts.Nil(err)
	})

	t.Run("✅ Should update the description of the second task", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.AddTask(createTask2(2))
		err := taskList.UpdateTask(2, "Updated Task")

		updatedAt := time.Now().Format(time.DateOnly)

		asserts.Equal(len(taskList.Tasks), 2)
		asserts.Equal(taskList.Tasks[1].Description, "Updated Task")
		asserts.Equal(taskList.Tasks[1].UpdatedAt.Format(time.DateOnly), string(updatedAt))
		asserts.GreaterOrEqual(taskList.Tasks[1].UpdatedAt.Format(time.DateOnly), string(updatedAt))
		asserts.Nil(err)
	})

	t.Run("❌ Should return an error when trying to update the description of a task that does not exist", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		err := taskList.UpdateTask(1, "Updated Task")

		asserts.Empty(taskList.Tasks)
		asserts.EqualError(err, "task with ID 1 not found")
	})

	t.Run("✅ Should print all tasks", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.AddTask(createTask2(2))

		err := taskList.loadFromFile()

		expected := joinMessage2(taskList)
		expected = expected + "\n--------------- Total Tasks: 2 ---------------\n"
		result := outputToString2(taskList.PrintAll)

		asserts.Nil(err)
		asserts.Equal(len(taskList.Tasks), 2)
		asserts.Equal(expected, result)
	})

	t.Run("✅ Should print all done tasks", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.AddTask(createTask2(2))
		taskList.AddTask(createTask2(3))
		taskList.MarkDone(3)

		expected := joinMessageWithFilter2(taskList, models.DONE) + "\n"
		result := outputToString2(taskList.PrintDone)

		asserts.Equal(len(taskList.Tasks), 3)
		asserts.Equal(expected, result)
	})

	t.Run("✅ Should print all in progress tasks", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.AddTask(createTask2(2))
		taskList.AddTask(createTask2(3))

		taskList.MarkInProgress(1)

		expected := joinMessageWithFilter2(taskList, models.IN_PROGRESS) + "\n"
		result := outputToString2(taskList.PrintInProgress)

		asserts.Equal(len(taskList.Tasks), 3)
		asserts.Equal(expected, result)
	})

	t.Run("✅ Should print all todo tasks", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.AddTask(createTask2(2))
		taskList.AddTask(createTask2(3))

		expected := joinMessageWithFilter2(taskList, models.TODO) + "\n"
		result := outputToString2(taskList.PrintTodo)

		asserts.Equal(len(taskList.Tasks), 3)
		asserts.Equal(expected, result)
	})

	t.Run("✅ Should print has no tasks", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")

		expected := "No tasks found" + "\n--------------- Total Tasks: 0 ---------------\n"
		result := outputToString2(taskList.PrintAll)
		asserts.Equal(len(taskList.Tasks), 0)
		asserts.Equal(expected, result)
	})

	t.Run("❌ Should print a message when there are no done tasks", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.AddTask(createTask2(2))

		result := outputToString2(taskList.PrintDone)
		asserts.Equal(len(taskList.Tasks), 2)
		asserts.Equal("No tasks found\n", result)
	})

	t.Run("❌ Should print a message when there are no in progress tasks", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(3))

		result := outputToString2(taskList.PrintInProgress)
		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Equal("No tasks found\n", result)
	})

	t.Run("❌ Should print a message when there are no todo tasks", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(3))
		taskList.MarkDone(1)

		result := outputToString2(taskList.PrintTodo)
		asserts.Equal(len(taskList.Tasks), 1)
		asserts.Equal("No tasks found\n", result)
	})

	t.Run("✅ Should mark a task as in progress", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.MarkInProgress(1)

		asserts.Equal(taskList.Tasks[0].Status, models.IN_PROGRESS)
	})

	t.Run("✅ Should mark a task as done", func(t *testing.T) {
		setup()

		taskList := NewJsonTaskStore("test.json")
		taskList.AddTask(createTask2(1))
		taskList.MarkDone(1)
	})
}

func joinMessage2(tasks *JsonTaskStore) string {
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

func joinMessageWithFilter2(tasks *JsonTaskStore, filter models.Status) string {
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

func outputToString2(callback func() error) string {

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
	err := callback()
	if err != nil {
		return err.Error()
	}

	// Close the write end of the pipe to signal EOF
	w.Close()

	// Read the output from the read end of the pipe
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func createTask2(id int) *models.Task {
	return &models.Task{
		Description: fmt.Sprintf("Task %d", id),
	}
}
