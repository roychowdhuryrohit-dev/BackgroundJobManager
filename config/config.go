package config

import (
	"log"
	"os"
	"sync"
)

const (
	HostAddr        = "HOST"
	TeamCSVPath     = "TEAM_CSV"
	BaselineCSVPath = "BASELINE_CSV"
)

// var ConfigMap map[string]string
var ConfigMap sync.Map

func Config() {

	if envar, ok := os.LookupEnv(HostAddr); ok {
		ConfigMap.Store(HostAddr, envar)
	} else {	
		log.Panicf("%s environment variable is not set.\n", HostAddr)
	}
	if envar, ok := os.LookupEnv(TeamCSVPath); ok {	
		ConfigMap.Store(TeamCSVPath, envar)
	} else {
		log.Panicf("%s environment variable is not set.\n", TeamCSVPath)
	}
	if envar, ok := os.LookupEnv(BaselineCSVPath); ok {	
		ConfigMap.Store(BaselineCSVPath, envar)
	} else {
		log.Panicf("%s environment variable is not set.\n", BaselineCSVPath)
	}
}
