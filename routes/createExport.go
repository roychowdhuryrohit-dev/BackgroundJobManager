package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/roychowdhuryrohit-dev/BackgroundJobManager/signals"
)

func CreateExport(context *gin.Context) {

	var parsedBody struct {
		FromDate string `json:"from" binding:"required"`
		ToDate   string `json:"to" binding:"required"`
	}

	if err := context.ShouldBindJSON(&parsedBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errorType": "invalid_json"})
		return
	}

	if fDate, err := time.Parse("2006-01-02", parsedBody.FromDate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errorType": "invalid_from_date"})
		return
	} else if tDate, err := time.Parse("2006-01-02", parsedBody.ToDate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errorType": "invalid_to_date"})
		return
	} else if fDate.After(tDate) {
		context.JSON(http.StatusBadRequest, gin.H{"errorType": "invalid_date_range"})
		return
	}

	sid := signals.GenSID("Jon-Doe", "createExport")

	sg := &signals.SignalsGroup{
		StopCh:      make(chan struct{}),
		StartCh:     make(chan struct{}),
		TerminateCh: make(chan struct{}),
	}

	signals.Signals.Store(sid, sg)
	//Start background goroutine for processing CSV
	signals.WorkersGroup.Add(1)
	go exportData(sg, sid, parsedBody.FromDate, parsedBody.ToDate)

	context.JSON(http.StatusOK, gin.H{"Job ID": sid})

}

func exportData(sg *signals.SignalsGroup, sid string, fromDateStr string, toDateStr string) {

	defer func() {
		signals.Clear(sid)
		log.Printf("Worker %s exited.\n", sid)
		signals.WorkersGroup.Done()
	}()

	fromDate, _ := time.Parse("2006-01-02", fromDateStr)
	toDate, _ := time.Parse("2006-01-02", toDateStr)

loop:
	for iDate := fromDate; ; {
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
			fDate := iDate.Format("2006-01-02")
			//Print i`th date.
			log.Printf("Job - %s Date - %s\n", sid, fDate)
			//Simulate time heavy task for exporting each date.
			time.Sleep(1 * time.Second)
			//Print log for each date export.
			log.Printf("Date %s exported for job %s \n", fDate, sid)
			iDate = iDate.AddDate(0, 0, 1)
			if iDate.YearDay() == toDate.YearDay() && iDate.Year() == toDate.Year() {
				break loop
			}
		}
	}
}
