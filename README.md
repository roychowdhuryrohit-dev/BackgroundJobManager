# Background Job Manager

This is an example service written in Golang and powered by [Gin](https://github.com/gin-gonic/gin) web framework, demonstrating a technique to stop, resume and terminate long-running background jobs by client API request.
This will ensure that the resources like compute/memory/storage/time are used efficiently and do not go into processing tasks that have already been stopped (and then to roll back the work done post the stop-action). 

 
## Features

 - Background workers are implemented by *goroutines* which are very lightweight and do not consume much resources. These threads are different from OS threads and are multiplexed by the Go scheduler.
 
 - Uses sync.Map which scales [better](https://weekly-geekly.github.io/articles/338718/index.html) than lock-based technique to achieve thread safety, especially in a large multi-core **">4 "** system *[in certain scenarios](https://golang.org/pkg/sync/#Map)*.
 
 - Each worker during "pause", consumes very little resource both in CPU time and memory. This is because :
	 - Go scheduler doesn't schedule a goroutine until data is received in channel. Channel operations tell the scheduler to schedule another goroutine, that’s why a         program doesn’t block forever on the same goroutine. This performs better than other techniques like long-polling(loop).
	 - File operations are implemented with the help of *"bufio"* package. This makes sure large files do not congest the memory as they are loaded in small *chunks*.

 -  Server exits *"gracefully"* on OS signals like SIGINT, SIGTERM (since SIGKILL cannot be catched). This happens in 2 steps :
	 - First  http.Server's built-in Shutdown() method gracefully shuts down the server without interrupting any active connections. Shutdown works by first closing all open listeners, then closing all idle connections, and then waiting indefinitely for connections to return to idle and then shut down.
	 - To terminate the background workers, first a *terminate* signal is sent to every running goroutines. Then with the help of sync.WaitGroup's Wait() method, server waits until all goroutines have returned fully. This can ensure workers have done all the necessary cleanups.

## Usage

 - To run locally, type:
 
	 `make`
	 
 - To build in Docker, type:
  
	 `make build_docker` 
	 
 - To run in Docker, type:
 
	 `make run_docker`

Make sure environment variables are set in file *.env* present in parent folder.

## API doc

 - *POST /uploadCSV -F file=@/path/file.csv*

    Response:
	*{"Job ID":"Jon-Doe-baselineCSV-1575836394-8469"}*
	
	Uploads and processes each row of file. 
	
 - *POST /createBulkTeam -F file=@/path/file.csv*

    Response:
	*{"Job ID":"Jon-Doe-bulkTeam-1575836394-8469"}*
	
	Creates each team in in file.
	
 - *POST /exportData -d '{
	"from":"2018-07-01",
	"to":"2018-08-01"
	}'*
	
	Response:
	*{"Job ID":"Jon-Doe-createExport-1575836394-8469"}*
	
	Exports data for each date row.
	
 - GET */updateTask?id=Jon-Doe-bulkTeam-1575829098-7618&action=start*
	
	Resumes background worker of that id.
	
 - GET */updateTask?id=Jon-Doe-bulkTeam-1575829098-7618&action=stop*
   
    Pauses background worker until it gets resumed or terminated.
   
 - GET */updateTask?id=Jon-Doe-bulkTeam-1575829098-7618&action=terminate*

	Terminates background worker.

## Future Work

 - Currently workers do not survive server restarts. This will require saving the state of the worker in a persistent storage or an in-memory database like Redis. This will be easy to implement now as the server shutdowns due to SIGINT, SIGTERM gracefully, where workers can save states and exit. 
 For fault-tolerance, regular backup of state will be required, thus consuming extra resources.