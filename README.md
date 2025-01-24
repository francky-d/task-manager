# Task manager
This project if just for learning purpose. 

It's a CLI that simulates task processing through  a json files provided as an argument. 

A task has and `ID`, a `Description`, a `Status` and A `Duration`

## How to run  it on local
You must have Go install on your machine. 

Once `Go` is install, build the program so that you may have an executable by runing : 

`go build main.go` : buid the program for your operating system.


## All the available commands. 


- `./main -help` : to view all the commands proposed by the program

- `./main -list` : to view all the tasks. Add `-n=<number_tasks_to_list>` to list a specific number of task.
    
- `./main -status=<STATUS>` : to list tasks with specific with a specific status
    - 1(not started) , 2(in progress), 3(done). Ex: -status=1

- `./main -process`: to process all the tasks that are not  yet done. 