package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/roychowdhuryrohit-dev/SampleJobManager/signals"

	"github.com/gin-gonic/gin"
)

func UpdateTask(context *gin.Context) {

	sid := context.Query("id")
	sid = strings.TrimSpace(sid)
	if sid == "" {
		context.JSON(http.StatusBadRequest, gin.H{"errorType": "invalid_id"})
		return
	}

	value, ok := signals.Signals.Load(sid)
	if !ok {
		context.JSON(http.StatusNotFound, gin.H{"errorType": "id_doesn't_exist"})
		return
	}

	action := context.Query("action")
	action = strings.TrimSpace(action)
	if action == "" || action != "stop" && action != "start" && action != "terminate" {
		context.JSON(http.StatusBadRequest, gin.H{"errorType": "invalid_action"})
		return
	}

	sg := value.(*signals.SignalsGroup)

	switch action {
	case "stop":
		select {
		case sg.StopCh <- struct{}{}:
			context.JSON(http.StatusOK, gin.H{"status": "Job " + sid + " stopped."})
		//5 sec timeout for blocking request.
		case <-time.After(5 * time.Second):
			context.JSON(http.StatusInternalServerError, gin.H{"errorType": "worker_unresponsive"})
		}
	case "start":
		select {
		case sg.StartCh <- struct{}{}:
			context.JSON(http.StatusOK, gin.H{"status": "Job " + sid + " started."})
		//5 sec timeout for blocking request.
		case <-time.After(5 * time.Second):
			context.JSON(http.StatusInternalServerError, gin.H{"errorType": "worker_unresponsive"})
		}
	case "terminate":
		select {
		case sg.TerminateCh <- struct{}{}:
			context.JSON(http.StatusOK, gin.H{"status": "Job " + sid + " terminated."})
		//5 sec timeout for blocking request.
		case <-time.After(5 * time.Second):
			context.JSON(http.StatusInternalServerError, gin.H{"errorType": "worker_unresponsive"})
		}
	}

}
