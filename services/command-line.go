package services

import (
	"flag"
	"fmt"
	"os"
	"task-tracker/models"
	"task-tracker/stores"
)

type (
	CommandLine interface {
		Run()
	}

	commandLine struct {
		store models.TaskStore
	}
)

func NewCommandLine(store *stores.JsonTaskStore) CommandLine {
	return &commandLine{
		store: store,
	}
}

func (c *commandLine) selectList(listType string) {
	tasks := c.store

	switch listType {
	case "todo":
		tasks.PrintTodo()
	case "in-progress":
		tasks.PrintInProgress()
	case "done":
		tasks.PrintDone()
	default:
		tasks.PrintAll()
	}
}

func (c *commandLine) Run() {
	addTaskSubCommand := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addTaskSubCommand.String("description", "", "Description of the task")

	updateTaskSubCommand := flag.NewFlagSet("update", flag.ExitOnError)
	updateDescription := updateTaskSubCommand.String("description", "", "Description of the task")
	updateId := updateTaskSubCommand.Int("id", 0, "ID of the task")

	deleteTaskSubCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteId := deleteTaskSubCommand.Int("id", 0, "ID of the task")

	markAsTaskDoneSubCommand := flag.NewFlagSet("mark-done", flag.ExitOnError)
	markAsTaskDoneId := markAsTaskDoneSubCommand.Int("id", 0, "ID of the task")

	markAsTaskInProgressSubCommand := flag.NewFlagSet("mark-in-progress", flag.ExitOnError)
	markAsTaskInProgressId := markAsTaskInProgressSubCommand.Int("id", 0, "ID of the task")

	listTaskSubCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listTodo := listTaskSubCommand.Bool("todo", false, "List tasks in todo status")
	listInProgress := listTaskSubCommand.Bool("in-progress", false, "List tasks in in-progress status")
	listDone := listTaskSubCommand.Bool("done", false, "List tasks in done status")

	if len(os.Args) < 2 {
		fmt.Println("Please provide a subcommand: add, update, delete, mark-done, mark-in-progress, list")
		os.Exit(1)
	}

	tasks := c.store

	switch os.Args[1] {
	case "add":
		addTaskSubCommand.Parse(os.Args[2:])
		tasks.AddTask(&models.Task{
			Description: *addDescription,
		})
	case "update":
		updateTaskSubCommand.Parse(os.Args[2:])
		tasks.UpdateTask(*updateId, *updateDescription)
	case "delete":
		deleteTaskSubCommand.Parse(os.Args[2:])
		tasks.RemoveTask(*deleteId)
	case "mark-done":
		markAsTaskDoneSubCommand.Parse(os.Args[2:])
		tasks.MarkDone(*markAsTaskDoneId)
	case "mark-in-progress":
		markAsTaskInProgressSubCommand.Parse(os.Args[2:])
		tasks.MarkInProgress(*markAsTaskInProgressId)
	case "list":
		listTaskSubCommand.Parse(os.Args[2:])
		if len(os.Args) < 3 {
			c.store.PrintAll()
			return
		}

		if len(os.Args) > 3 {
			fmt.Println("It is not possible to list tasks by more than one status at the same time")
			os.Exit(1)
		}

		listByStatus := os.Args[2]

		if *listTodo {
			listByStatus = "todo"
		} else if *listInProgress {
			listByStatus = "in-progress"
		} else if *listDone {
			listByStatus = "done"
		}
		c.selectList(listByStatus)

	default:
		fmt.Println("Invalid subcommand. Expected: add, update, delete, mark-done, mark-in-progress, list")
		os.Exit(1)
	}
}
