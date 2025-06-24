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
You must have  `docker` and `docker compose` on your machine

### Run with Docker Compose
Build and start the app using Docker Compose:
```bash
docker compose up -d
```

#### Build and Run Inside the Container
After starting the service, you can build with the command :
```bash
docker compose exec task-manager go build main.go
```

This will create an executable inside the container `task-manager`

## Usage

### Basic Commands (Docker)
```bash
# View all available commands and options
docker compose exec task-manager ./main -help

# List all tasks
docker compose exec task-manager ./main -list

# List specific number of tasks
docker compose exec task-manager ./main -list -num=5

# Filter tasks by status
docker compose exec task-manager ./main -status=1  # not_started tasks
docker compose exec task-manager ./main -status=2  # in_progress tasks
docker compose exec task-manager ./main -status=3  # done tasks

# Process tasks
docker compose exec task-manager ./main -process           # Process with default 5 workers
docker compose exec task-manager ./main -process -workers=3  # Process with 3 workers
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
docker compose exec task-manager ./main -list -num=3
```

### Process Tasks with 2 Workers 
```bash
docker compose exec task-manager ./main -process -workers=2
```

### View Tasks with "Done" Status 
```bash
docker compose exec task-manager ./main -status=3
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