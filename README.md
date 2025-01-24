# Task manager
This project if just for learning purpose. 

It's a CLI that simulates task processing through  a json files provided as an argument. 

A task has and `ID`, a `Description`, a `Status` and A `Duration`

## How run  it on local
You must have Go install on your machine. 

Once `Go` is install, build the program so that you may have an executable by runing : 

All the availables command. 
`go build main.go` : buid the program for your operating system.

`./main -help` : to view all the commands proposed by the program

 `./main -list` : to view all the tasks. Add `-n=<number_tasks_to_list>` to list a specific number of task.

 `./main -status=<STATUS>` : to list tasks with specific with a specific status

 `./main -process`: to process all the tasks that are not  yet done. 