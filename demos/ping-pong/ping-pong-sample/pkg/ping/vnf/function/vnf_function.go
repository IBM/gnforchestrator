package function

import (
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/ping/config"
	"os"
	"log"
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"strconv"
	"sync"
)

const (
	PingMessageFormat = "ping version %s: %s" // ping version 'X': 'pong's response'
)

type PingVnfFunction struct {
	port string
	version string
	pingVnfConfig *config.PingConfig
	configLock sync.Mutex
	managementChannel <- chan *config.PingConfig
	stopChan chan struct{}
}

func (f *PingVnfFunction) hello(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method other than GET is not supported for this endpoint", http.StatusNotFound)
		return
	}
	log.Println("hello called")
	fmt.Fprintln(w, "Hello from Ping VNF")
}

func (f *PingVnfFunction) ping(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("method other than GET is not supported for this endpoint")
		http.Error(w, "method other than GET is not supported for this endpoint", http.StatusNotFound)
		return
	}
	log.Println("ping called")
	numStr := strings.TrimPrefix(req.URL.Path, "/ping/")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		log.Println("failed to convert argument to int -" + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.configLock.Lock()
	if f.pingVnfConfig == nil {
		f.configLock.Unlock()
		log.Println("ping config was not set yet")
		http.Error(w, "ping config was not set yet", http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("http://%s:%s/pong", f.pingVnfConfig.PongAddress, f.pingVnfConfig.PongPort)
	log.Printf("pong endpoint was set to http://%s:%s/pong", f.pingVnfConfig.PongAddress, f.pingVnfConfig.PongPort)
	f.configLock.Unlock()
	var builder strings.Builder
	for i := 0; i<num; i++ {
		log.Printf("trying to ping...")
		response, err := http.Get(url)
		if err != nil {
			log.Println("ping failed - "+ err.Error())
			http.Error(w, err.Error(), response.StatusCode)
			return
		}
		body, err := ioutil.ReadAll(response.Body)
		builder.Write(body)
		response.Body.Close()
	}

	fmt.Fprintf(w, PingMessageFormat, f.version, builder.String())
}

func (f *PingVnfFunction) configChangeHandler() {
	for {
		select {
		case <-f.stopChan:
			return
		case pingVnfConfig := <- f.managementChannel:
			f.configLock.Lock()
			f.pingVnfConfig = pingVnfConfig
			f.configLock.Unlock()
			log.Printf("ping vnf config updated successfully with %s", pingVnfConfig)
		}
	}
}

func (f *PingVnfFunction) Start(managementChannel <- chan *config.PingConfig) { // read only channel
	f.managementChannel = managementChannel
	f.port = os.Getenv("VNF_FUNCTIONAL_PORT")
	f.version = os.Getenv("VNF_VERSION")

	if f.port == "" || f.version == "" {
		log.Fatal("the expected arguments are not set in the env vars")
		return
	}
	f.stopChan = make(chan struct{})

	go f.configChangeHandler() // run as go routine, can be stopped by using the Stop function

	log.Printf("VNF functional Server is listening on port %s...", f.port)
	// Http listener
	http.HandleFunc("/hello", f.hello)
	http.HandleFunc("/ping/", f.ping)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s",f.port), nil))
}

func (f *PingVnfFunction) Stop() {
	f.stopChan <- struct{}{}
}

