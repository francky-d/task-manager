package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	NOT_STARTED         = "not_started"
	IN_PROGRESS         = "in_progress"
	DONE                = "done"
	ESCAPTE_COLOR_START = "\033["
	ESCAPE_COLOR_END    = "\033[0m"
	RED_BOLD            = ESCAPTE_COLOR_START + "1;31m"
	YELLOW_BOLD         = ESCAPTE_COLOR_START + "1;33m"
	GREEN_BOLD          = ESCAPTE_COLOR_START + "1;32m"
	VIOLET_BOLD         = ESCAPTE_COLOR_START + "1;34m"
	WHITE               = ESCAPTE_COLOR_START + "0m"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Duration    string `json:"duration"`
	Status      string
}

func getMatchingColor(status string) string {
	switch status {
	case NOT_STARTED:
		return RED_BOLD
	case IN_PROGRESS:
		return YELLOW_BOLD
	case DONE:
		return GREEN_BOLD
	default:
		return WHITE
	}
}

func PrintHeader() {
	PrintLine("ID", "Description", "Duration", "Status")
	fmt.Print("\n")
}

func PrintLine(ID, Description, Duration, Status string) {
	fmt.Printf("%s \t\t %-24s \t\t %-8s \t\t %-12s \n", ID, Description, Duration, Status)
}

func colorStatus(status string) string {
	return fmt.Sprintf("%s %s %s", getMatchingColor(status), status, ESCAPE_COLOR_END)
}

func printTasks(tasks []Task) {
	PrintHeader()
	for _, task := range tasks {
		PrintLine(strconv.Itoa(task.ID), task.Description, task.Duration, colorStatus(task.Status))
		time.Sleep(time.Second)
	}
}

func main() {
	var allTask []Task
	filePath := flag.String("file", "./task.json", "Give the file path")
	flag.Parse()

	if *filePath == "" {
		log.Fatalf("Please provide the file path")
	}

	file, err := os.Open(*filePath)

	if err != nil {
		log.Fatalf("Error while opening %s : %v", *filePath, err)
	}

	json.NewDecoder(file).Decode(&allTask)

	printTasks(allTask)

}
