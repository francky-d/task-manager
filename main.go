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

func (task Task) IsInProgress() bool {
	return task.Status == IN_PROGRESS
}

func (task Task) IsNotStarted() bool {
	return task.Status == NOT_STARTED
}

func (task Task) IsDone() bool {
	return task.Status == DONE
}

func (task Task) IsNotYetDone() bool {
	return task.Status != DONE
}

func (task *Task) UpdateStatusTo(newStatus string) {
	task.Status = newStatus
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
		time.Sleep(time.Nanosecond)
	}
}

func GetTaskNotYetDone(tasks []Task) []Task {
	tasksNotYetDone := []Task{}
	for _, task := range tasks {
		if task.IsNotYetDone() {
			tasksNotYetDone = append(tasksNotYetDone, task)
		}
	}
	return tasksNotYetDone
}

func NotifyTaskIsBeingProcessed(task Task, workerID int) string {
	return fmt.Sprintf("Task (%s%d%s) is being processed by worker %d", YELLOW_BOLD, task.ID, ESCAPE_COLOR_END, workerID)
}

func NotifyTaskIsDone(task Task, workerID int) string {
	return fmt.Sprintf("Task (%s%d%s) has been done in %d(s) by worker %d", GREEN_BOLD, task.ID, ESCAPE_COLOR_END, task.Duration, workerID)
}

func updateTaskStatusInList(task Task, taskList *[]Task, index int, notifications chan string, workerID int) {
	if task.IsNotStarted() {
		task.UpdateStatusTo(IN_PROGRESS)
		notifications <- NotifyTaskIsBeingProcessed(task, workerID)

		task.UpdateStatusTo(DONE)
		notifications <- NotifyTaskIsDone(task, workerID)

	}

	if task.IsInProgress() {
		time.Sleep(time.Duration(task.Duration) * time.Second)
		task.UpdateStatusTo(DONE)
		notifications <- NotifyTaskIsDone(task, workerID)
	}

	(*taskList)[index] = task
}

func processTask(taskID int, taskList *[]Task, notifications chan string, workerID int) {
	for index, task := range *taskList {
		if task.ID == taskID {
			updateTaskStatusInList(task, taskList, index, notifications, workerID)
		}
	}
}

func filterTaskBasedOnStatus(tasks []Task, status string) []Task {
	tasksFiltered := []Task{}
	for _, task := range tasks {
		if task.Status == status {
			tasksFiltered = append(tasksFiltered, task)
		}
	}
	return tasksFiltered

}

func parseArgumentsFromFlags() (*string, *bool, *bool, *int, *int, *int) {
	arguments := map[string]string{
		"file":    "the file of the path",
		"list":    "print the list of all tasks",
		"num":     "use with -list to print a specific number of task",
		"status":  "use to print task with specifi status. values are 1(not started) , 2(in progress), 3(done). Ex: -status=1",
		"process": "use to process task that are not yet done",
		"workers": "use to process task that are not yet done. 5 by default",
	}

	filePath := flag.String("file", "./task.json", arguments["file"])
	shouldList := flag.Bool("list", false, arguments["list"])
	shouldProcessTasks := flag.Bool("process", false, arguments["process"])
	numberOfLineToPrint := flag.Int("num", 0, arguments["num"])
	status := flag.Int("status", 0, arguments["status"])
	workers := flag.Int("workers", 5, arguments["workers"])

	flag.Parse()

	return filePath, shouldList, shouldProcessTasks, numberOfLineToPrint, status, workers
}

func getTasksFromFile(filePath *string) []Task {
	var allTasks []Task

	if *filePath == "" {
		log.Fatalf("Please provide the file path")
	}

	file, err := os.Open(*filePath)

	if err != nil {
		log.Fatalf("Error while opening %s : %v", *filePath, err)
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&allTasks)

	if err != nil {
		log.Fatalf("Error while decoding %s : %v", *filePath, err)
	}

	return allTasks
}

func startProcessingTasks(allTasks *[]Task, workers *int) {
	wg := sync.WaitGroup{}
	notifications := make(chan string)

	tasksNotYeDone := GetTaskNotYetDone(*allTasks)

	if len(tasksNotYeDone) == 0 {
		log.Fatal("No tasks to process")
	}

	taskDispatcher := make(chan int, len(tasksNotYeDone))

	for _, task := range tasksNotYeDone {
		taskDispatcher <- task.ID
	}

	close(taskDispatcher)

	for i := 1; i <= *workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			//ProcessTaskUntilNoneRemain
			for taskID := range taskDispatcher {
				processTask(taskID, allTasks, notifications, workerID)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(notifications)
	}()

	for message := range notifications {
		fmt.Println(message)
	}

	os.Exit(0)
}

func filterTasks(allTasks []Task, status *int) {
	var finalStatus string
	switch *status {
	case 1:
		finalStatus = NOT_STARTED
	case 2:
		finalStatus = IN_PROGRESS
	case 3:
		finalStatus = DONE
	default:
		log.Fatal("status must be 1, 2, or 3 :  1=>not_started, 2=>in_progress, 3=>done")
	}

	tasks := filterTaskBasedOnStatus(allTasks, finalStatus)

	PrintTasks(tasks)
	os.Exit(0)
}

func listTasks(allTasks []Task, numberOfLineToPrint *int) {
	if *numberOfLineToPrint > 0 {
		PrintTasks(allTasks[:*numberOfLineToPrint])
		os.Exit(0)
	}
	PrintTasks(allTasks)
	os.Exit(0)
}

func main() {
	filePath, shouldList, shouldProcessTasks, numberOfLineToPrint, status, workers := parseArgumentsFromFlags()

	allTasks := getTasksFromFile(filePath)

	if *shouldList {
		listTasks(allTasks, numberOfLineToPrint)
	}

	if *status != 0 {
		filterTasks(allTasks, status)
	}

	if *shouldProcessTasks {
		startProcessingTasks(&allTasks, workers)
	}
}
