package main

import (
	"context"
	"log"
	"net/http"
	"os"
	ossignal "os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/roychowdhuryrohit-dev/SampleJobManager/config"
	"github.com/roychowdhuryrohit-dev/SampleJobManager/routes"
	"github.com/roychowdhuryrohit-dev/SampleJobManager/signals"
)

func main() {

	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20

	signals.Init()
	config.Config()

	router.POST("/uploadCSV", routes.CreateUpload)
	router.POST("/exportData", routes.CreateExport)
	router.POST("/createBulkTeam", routes.CreateBulkTeam)
	router.GET("/updateTask", routes.UpdateTask)

	host, _ := config.ConfigMap.Load(config.HostAddr)

	server := &http.Server{
		Addr:    host.(string),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(err.Error())
		}
	}()
	//Gracefully shut down server on OS interrupt.
	quit := make(chan os.Signal, 1)
	ossignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server.")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Println("Server shutdown error - " + err.Error())
	}

	//Signal all workers to terminate
	signals.CloseAll()
	//Wait for remaining workers to terminate.
	signals.WorkersGroup.Wait()

	log.Println("All workers terminated.")

	log.Println("---Server exited gracefully---")

}
