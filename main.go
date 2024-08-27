package main

import (
	"task-tracker/services"
	"task-tracker/stores"
)

/*
TODO: refactor as using done/todo/in-progress as boolean flags
# Listing all tasks
task-cli list

# Listing tasks by status
task-cli list done
task-cli list todo
task-cli list in-progress
*/

func main() {
	jsonStore := stores.NewJsonTaskStore("tasks.json")
	commandLine := services.NewCommandLine(jsonStore)
	commandLine.Run()
}
