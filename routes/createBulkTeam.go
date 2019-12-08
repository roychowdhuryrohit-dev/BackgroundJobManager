package routes

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/roychowdhuryrohit-dev/BackgroundJobManager/config"
	"github.com/roychowdhuryrohit-dev/BackgroundJobManager/signals"
)

func CreateBulkTeam(context *gin.Context) {

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Recovered from ", rec)
			context.JSON(http.StatusInternalServerError, gin.H{})
		}
	}()

	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errorType": "File not accepted"})
		return
	}

	sid := signals.GenSID("Jon-Doe", "bulkTeam")
	basepath, ok := config.ConfigMap.Load(config.TeamCSVPath)
	if !ok {
		log.Panicln("Basepath not found.")
	}
	filepath := path.Join(basepath.(string), sid+".csv")

	if err := context.SaveUploadedFile(file, filepath); err != nil {
		log.Panicln(err.Error())
	}

	sg := &signals.SignalsGroup{
		StopCh:      make(chan struct{}),
		StartCh:     make(chan struct{}),
		TerminateCh: make(chan struct{}),
	}

	signals.Signals.Store(sid, sg)
	//Start background goroutine for processing CSV
	signals.WorkersGroup.Add(1)
	go processBulkTeam(sg, sid, filepath)

	context.JSON(http.StatusOK, gin.H{"Job ID": sid})

}

func processBulkTeam(sg *signals.SignalsGroup, sid string, filepath string) {
	var (
		csvFile *os.File
		err     error
	)
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Recovered from ", rec)
		}
		if csvFile != nil {
			if err := csvFile.Close(); err != nil {
				log.Println(err.Error())
			}
		}
		signals.Clear(sid)
		log.Printf("Worker %s exited.\n", sid)
		signals.WorkersGroup.Done()
	}()

	csvFile, err = os.Open(filepath)
	if err != nil {
		log.Panicln(err.Error())
	}

	reader := csv.NewReader(csvFile)

	reader.LazyQuotes = true
loop:
	for i := 0; ; {
		select {
		case <-sg.StopCh:
			log.Printf("Job %s stopped.\n", sid)
			//Sleep current worker so that it doesn't consume CPU
		innerselect:
			select {
			case <-sg.StopCh:
				log.Printf("Job %s already stopped.\n", sid)
				goto innerselect
			case <-sg.StartCh:
				log.Printf("Job %s resumed.\n", sid)
			case <-sg.TerminateCh:
				return
			}
		case <-sg.TerminateCh:
			return
		case <-sg.StartCh:
			log.Printf("Job %s already resumed.\n", sid)
		default:
			line, err := reader.Read()
			if err == io.EOF {
				break loop
			} else if err != nil {
				log.Panicln(err.Error())
			}
			//Print each team row.
			log.Printf("Job - %s Team members - %+q\n", sid, line)
			//Simulate time heavy task for creating each team in CSV.
			time.Sleep(1 * time.Second)
			//Print log for each created team.
			i++
			log.Printf("Team %d created for job %s \n", i, sid)

		}
	}
}
