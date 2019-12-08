package signals

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type SignalsGroup struct {
	StopCh      chan struct{}
	StartCh     chan struct{}
	TerminateCh chan struct{}
}

//var Signals map[string]*SignalsGroup

//Signals Thread safe map to store channels of each worker goroutine.
var Signals sync.Map

var WorkersGroup *sync.WaitGroup

func Init() {
	// Signals = map[string]*SignalsGroup{}
	WorkersGroup = &sync.WaitGroup{}
}

//Clear Close all the channels and remove key entry
func Clear(sid string) {
	value, ok := Signals.Load(sid)
	if ok {
		sg := value.(*SignalsGroup)

		Signals.Delete(sid)
		close(sg.StartCh)
		close(sg.StopCh)
		close(sg.TerminateCh)

	}
}

func CloseAll() {
	Signals.Range(
		func(key interface{}, value interface{}) bool {
			sg := value.(*SignalsGroup)
			//Terminate all workers
			sg.TerminateCh <- struct{}{}
			return true
		},
	)

}

func GenSID(username string, jobType string) string {
	rand.Seed(time.Now().UnixNano())
	return username + "-" + jobType + "-" + strconv.Itoa(int(time.Now().Unix())) + "-" + strconv.Itoa(rand.Intn(10000))
}
