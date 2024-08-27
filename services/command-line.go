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

func (c *commandLine) addTaskCommand() {
	addTaskSubCommand := flag.NewFlagSet("add", flag.ExitOnError)
	addTaskSubCommand.Parse(os.Args[2:])
	addDescription := addTaskSubCommand.String("description", "", "Description of the task")
	c.store.AddTask(&models.Task{
		Description: *addDescription,
	})
}

func (c *commandLine) updateTaskCommand() {
	updateTaskSubCommand := flag.NewFlagSet("update", flag.ExitOnError)
	updateDescription := updateTaskSubCommand.String("description", "", "Description of the task")
	updateId := updateTaskSubCommand.Int("id", 0, "ID of the task")
	updateTaskSubCommand.Parse(os.Args[2:])
	c.store.UpdateTask(*updateId, *updateDescription)
}

func (c *commandLine) deleteTaskCommand() {
	deleteTaskSubCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteId := deleteTaskSubCommand.Int("id", 0, "ID of the task")
	deleteTaskSubCommand.Parse(os.Args[2:])
	c.store.RemoveTask(*deleteId)
}

func (c *commandLine) markAsTaskDoneCommand() {
	markAsTaskDoneSubCommand := flag.NewFlagSet("mark-done", flag.ExitOnError)
	markAsTaskDoneId := markAsTaskDoneSubCommand.Int("id", 0, "ID of the task")
	markAsTaskDoneSubCommand.Parse(os.Args[2:])
	c.store.MarkDone(*markAsTaskDoneId)
}

func (c *commandLine) markAsTaskInProgressCommand() {
	markAsTaskInProgressSubCommand := flag.NewFlagSet("mark-in-progress", flag.ExitOnError)
	markAsTaskInProgressId := markAsTaskInProgressSubCommand.Int("id", 0, "ID of the task")
	markAsTaskInProgressSubCommand.Parse(os.Args[2:])
	c.store.MarkInProgress(*markAsTaskInProgressId)
}

func (c *commandLine) listCommand() {
	listTaskSubCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listTodo := listTaskSubCommand.Bool("todo", false, "List tasks in todo status")
	listInProgress := listTaskSubCommand.Bool("in-progress", false, "List tasks in in-progress status")
	listDone := listTaskSubCommand.Bool("done", false, "List tasks in done status")
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
}

func (c *commandLine) Run() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a subcommand: add, update, delete, mark-done, mark-in-progress, list")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		c.addTaskCommand()
	case "update":
		c.updateTaskCommand()
	case "delete":
		c.deleteTaskCommand()
	case "mark-done":
		c.markAsTaskDoneCommand()
	case "mark-in-progress":
		c.markAsTaskInProgressCommand()
	case "list":
		c.listCommand()

	default:
		fmt.Println("Invalid subcommand. Expected: add, update, delete, mark-done, mark-in-progress, list")
		os.Exit(1)
	}
}
