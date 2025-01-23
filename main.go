package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
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
	Duration    int64  `json:"duration"`
	Status      string `json:"status"`
}

func GetMatchingColor(status string) string {
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

func ColorStatus(status string) string {
	return fmt.Sprintf("%s %s %s", GetMatchingColor(status), status, ESCAPE_COLOR_END)
}

func PrintTasks(tasks []Task) {
	PrintHeader()
	for _, task := range tasks {
		PrintLine(strconv.Itoa(task.ID), task.Description, strconv.FormatInt(task.Duration, 10), ColorStatus(task.Status))
		time.Sleep(time.Second)
	}
}

func GetTaskNotYetDone(tasks []Task) []Task {
	tasksNotYeDone := []Task{}
	for _, tasks := range tasks {
		if tasks.Status != DONE {
			tasksNotYeDone = append(tasksNotYeDone, tasks)
		}
	}
	return tasksNotYeDone
}

func NotifyTaskInProcessState(task Task) string {
	return fmt.Sprintf("Task (%s%d%s) is being processed", YELLOW_BOLD, task.ID, ESCAPE_COLOR_END)
}

func NotifyTaskIsDone(task Task) string {
	return fmt.Sprintf("Task (%s%d%s) has been done in %ds", GREEN_BOLD, task.ID, ESCAPE_COLOR_END, task.Duration)
}

func updateTaskStatusInList(task Task, taskList *[]Task, index int, notifications chan string) {
	if task.Status == NOT_STARTED {
		(*taskList)[index].Status = IN_PROGRESS
		notifications <- NotifyTaskInProcessState(task)
		(*taskList)[index].Status = DONE
		notifications <- NotifyTaskIsDone(task)

	}

	if task.Status == IN_PROGRESS {
		time.Sleep(time.Duration(task.Duration) * time.Second)
		(*taskList)[index].Status = DONE
		notifications <- NotifyTaskIsDone(task)
	}
}

func processTask(taskID int, taskList *[]Task, notifications chan string) {
	for index, task := range *taskList {
		if task.ID == taskID {
			updateTaskStatusInList(task, taskList, index, notifications)
		}
	}
}

func main() {
	wg := sync.WaitGroup{}
	notifications := make(chan string)

	var allTasks []Task
	filePath := flag.String("file", "./task.json", "Give the file path")
	flag.Parse()

	if *filePath == "" {
		log.Fatalf("Please provide the file path")
	}

	file, err := os.Open(*filePath)

	if err != nil {
		log.Fatalf("Error while opening %s : %v", *filePath, err)
	}

	json.NewDecoder(file).Decode(&allTasks)

	for _, task := range GetTaskNotYetDone(allTasks) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processTask(task.ID, &allTasks, notifications)
		}()
	}

	go func() {
		wg.Wait()
		close(notifications)
	}()

	for message := range notifications {
		fmt.Println(message)
	}

}
