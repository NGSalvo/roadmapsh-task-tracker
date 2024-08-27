package main

import (
	"task-tracker/services"
	"task-tracker/stores"
)

func main() {
	jsonStore := stores.NewJsonTaskStore("tasks.json")
	commandLine := services.NewCommandLine(jsonStore)
	commandLine.Run()
}
