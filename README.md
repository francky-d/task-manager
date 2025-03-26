# Task Manager

A command-line task processing system that simulates concurrent task execution from a JSON configuration file. This project is created for learning purposes, demonstrating Go's concurrency features using goroutines and channels.

## Features
- List all tasks or a specific number of tasks
- Filter tasks by status
- Process tasks concurrently with configurable number of workers
- Real-time task processing notifications with colored output
- JSON-based task configuration

## Task Structure
Each task in the JSON file has the following properties:
- `ID`: Unique identifier for the task
- `Description`: Text describing the task
- `Duration`: Time in seconds the task takes to complete
- `Status`: Current state of the task (not_started, in_progress, done)

## Installation
You must have Go installed on your machine. 

### Build the Program
```bash
go build main.go
```

## Usage

### Basic Commands
```bash
# View all available commands and options
./main -help

# List all tasks
./main -list

# List specific number of tasks
./main -list -num=5

# Filter tasks by status
./main -status=1  # not_started tasks
./main -status=2  # in_progress tasks
./main -status=3  # done tasks

# Process tasks
./main -process           # Process with default 5 workers
./main -process -workers=3  # Process with 3 workers
```

### Command Options
- `-file`: Specify the JSON file path (default: "./task.json")
- `-list`: Print the list of all tasks
- `-num`: Use with -list to print a specific number of tasks
- `-status`: Filter tasks by status (1=not_started, 2=in_progress, 3=done)
- `-process`: Start processing tasks that are not yet done
- `-workers`: Number of concurrent workers (default: 5)

### Task Status Codes
1. `not_started` (Status code: 1)
2. `in_progress` (Status code: 2)
3. `done` (Status code: 3)

## JSON File Format
Create a task.json file in the following format:
```json
[
  {
    "id": 1,
    "description": "Task description",
    "duration": 5,
    "status": "not_started"
  },
  ...
]
```

## Examples

### List First 3 Tasks
```bash
./main -list -num=3
```

### Process Tasks with 2 Workers
```bash
./main -process -workers=2
```

### View Tasks with "Done" Status
```bash
./main -status=3
```

## Output Colors
The program uses colored output to indicate task status:
- Red: Not Started
- Yellow: In Progress
- Green: Done

## Notes
- The program will exit automatically when all tasks are processed
- Task processing is concurrent and the order of completion may vary
- Use Ctrl+C to stop the program during processing