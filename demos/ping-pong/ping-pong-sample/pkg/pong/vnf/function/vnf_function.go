package function

import (
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/pong/config"
	"os"
	"log"
	"fmt"
	"net/http"
	"sync"
)

const (
	PongMessageFormat = "pong version %s message %s."
)

type PongVnfFunction struct {
	port string
	version string
	pongVnfConfig *config.PongConfig
	configLock sync.Mutex
	managementChannel <- chan *config.PongConfig
	stopChan chan struct{}
}

func (f *PongVnfFunction) hello(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("method other than GET is not supported for this endpoint")
		http.Error(w, "method other than GET is not supported for this endpoint", http.StatusNotFound)
		return
	}
	log.Println("hello called")
	fmt.Fprintln(w, "Hello from Pong VNF")
}

func (f *PongVnfFunction) pong(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("method other than GET is not supported for this endpoint")
		http.Error(w, "method other than GET is not supported for this endpoint", http.StatusNotFound)
		return
	}
	log.Println("pong called")
	f.configLock.Lock()
	defer f.configLock.Unlock()
	if f.pongVnfConfig == nil {
		log.Println("pong config was not set yet")
		http.Error(w, "pong config was not set yet", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, PongMessageFormat, f.version, f.pongVnfConfig.Message)
}

func (f *PongVnfFunction) configChangeHandler() {
	for {
		select {
		case <-f.stopChan:
			return
		case pongVnfConfig := <- f.managementChannel:
			f.configLock.Lock()
			f.pongVnfConfig = pongVnfConfig
			f.configLock.Unlock()
			log.Printf("pong vnf config updated successfully with body %s", pongVnfConfig)
		}
	}
}

func (f *PongVnfFunction) Start(managementChannel <- chan *config.PongConfig) { // read only channel
	f.managementChannel = managementChannel
	f.port = os.Getenv("VNF_FUNCTIONAL_PORT")
	f.version = os.Getenv("VNF_VERSION")

	if f.port == "" || f.version == "" {
		log.Fatal("the expected arguments are not set in the env vars")
		return
	}
	f.stopChan = make(chan struct{})

	go f.configChangeHandler() // run as go routine, can be stopped by using the Stop function

	log.Printf("VNF functional Server is listening on port %s ...", f.port)
	// Http listener
	http.HandleFunc("/hello", f.hello)
	http.HandleFunc("/pong", f.pong)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s",f.port), nil))
}

func (f *PongVnfFunction) Stop() {
	f.stopChan <- struct{}{}
}

