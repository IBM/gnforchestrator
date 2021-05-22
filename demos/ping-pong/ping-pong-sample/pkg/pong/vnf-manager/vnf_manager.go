package vnf_manager

import (
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/protos"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/pong/config"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"sync"
	"strings"
	"time"
)

type PongVnfManager struct {
	vnfManagerPort    string
	vnfAddress        string
	vnfManagementPort string
	debugMode		  common.DebugMode
	CustomMetric	  prometheus.Gauge
	prometheusLock 	  sync.Mutex
}

func (manager *PongVnfManager) hello(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("method other than GET is not supported for this endpoint")
		http.Error(w, "method other than GET is not supported for this endpoint", http.StatusNotFound)
		return
	}
	log.Println("hello called")
	fmt.Fprintln(w, "Hello from Pong VNF Manager")
}

func (manager *PongVnfManager) handleConfigure(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Println("method other than POST is not supported for this endpoint")
		http.Error(w, "method other than POST is not supported for this endpoint", http.StatusNotFound)
		return
	}
	bodyBytes, err := ioutil.ReadAll(req.Body)
	pongConfig := &config.PongConfig{}
	err = json.Unmarshal(bodyBytes, pongConfig)
	if err != nil {
		log.Println("error while parsing the body - " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("configure called with body %s",pongConfig)
	// got post request to set message and pong port
	err = manager.configureVnf(pongConfig)
	if err != nil {
		log.Println("failed to configre vnf - "+err.Error())
		http.Error(w, "failed to configre vnf - "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "finished configuring pong vnf successfully with body %s", pongConfig)
}

func (manager *PongVnfManager) handleHealthz(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("method other than GET is not supported for this endpoint")
		http.Error(w, "method other than GET is not supported for this endpoint", http.StatusNotFound)
		return
	}
	if manager.debugMode.Healthy {
		w.WriteHeader(200)
		fmt.Fprintln(w, "ok")
	} else {
		w.WriteHeader(500)
		fmt.Fprintln(w, "unhealthy")
	}
}

func (manager *PongVnfManager) handleDebugUnhealthy(w http.ResponseWriter, req *http.Request) {
	manager.debugMode.Healthy = false
	w.WriteHeader(200)
	fmt.Fprintln(w, "ok")
}

func (manager *PongVnfManager) handleDebugStatus(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	manager.prometheusLock.Lock()
	metricValue := manager.debugMode.MetricValue
	manager.prometheusLock.Unlock()
	fmt.Fprintf(w, "{Healthy: %s, Metric: %d }", strconv.FormatBool(manager.debugMode.Healthy), metricValue)
}

func (manager *PongVnfManager) handleReadinessCheck(w http.ResponseWriter, req *http.Request) {
	grpcClientConnection, err := grpc.Dial(fmt.Sprintf("%s:%s", manager.vnfAddress, manager.vnfManagementPort), grpc.WithInsecure())
	defer grpcClientConnection.Close()
	if err != nil {
		w.WriteHeader(500)
		log.Println("unable to send readiness check to vnf, error: " + err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	client := protos.NewPongVnfClient(grpcClientConnection)
	readinessRequest := &protos.ReadinessRequest{}
	readinessResponse := &protos.ReadinessResponse{}
	readinessResponse, err = client.ReadinessCheck(context.Background(), readinessRequest)
	if err != nil {
		w.WriteHeader(500)
		log.Println("unable to send readiness check to vnf, error: " + err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if readinessResponse.Ready == false {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Not Ready")
		return
	}
	w.WriteHeader(200)
	fmt.Fprintln(w, "Ready")
}

func (manager *PongVnfManager) handleDebugMetric(w http.ResponseWriter, req *http.Request) {
	log.Println("debug metric called")
	numStr := strings.TrimPrefix(req.URL.Path, "/debug/metric/")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		log.Println("failed to convert argument to int -" + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	manager.prometheusLock.Lock()
	manager.debugMode.MetricValue = num
	manager.prometheusLock.Unlock()
	w.WriteHeader(200)
	fmt.Fprintln(w, "ok")
}

func (manager *PongVnfManager) publishMetrics() {
	// publish continuous metrics value every two seconds
	previousValue := 0
	for {
		manager.prometheusLock.Lock()
		value := manager.debugMode.MetricValue
		manager.prometheusLock.Unlock()
		manager.CustomMetric.Add(float64(value-previousValue))
		previousValue = value
		time.Sleep(time.Second * 2)
	}
}

func (manager *PongVnfManager) configureVnf(pongConfig *config.PongConfig) error {
	grpcClientConnection, err := grpc.Dial(fmt.Sprintf("%s:%s", manager.vnfAddress, manager.vnfManagementPort), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer grpcClientConnection.Close()

	client := protos.NewPongVnfClient(grpcClientConnection)
	pongConfigRequest := &protos.PongVnfConfigRequest{
		Message:  pongConfig.Message,
	}
	_, err = client.Configure(context.Background(), pongConfigRequest)
	return err
}

func (manager *PongVnfManager) Start() {
	manager.vnfManagerPort = os.Getenv("VNFMANAGER_PORT")
	manager.vnfAddress = os.Getenv("VNF_ADDRESS")
	manager.vnfManagementPort = os.Getenv("VNF_MANAGEMENT_PORT")
	manager.debugMode = common.DebugMode { Healthy:true, MetricValue:1 }

	if manager.vnfManagerPort == "" || manager.vnfAddress == "" || manager.vnfManagementPort == "" {
		log.Fatal("the expected arguments are not set in the env vars")
		return
	}

	log.Printf("VNFM Server is listening on port %s...", manager.vnfManagerPort)
	http.HandleFunc("/hello", manager.hello)
	http.HandleFunc("/configure", manager.handleConfigure)
	http.HandleFunc("/healthz", manager.handleHealthz)
	http.HandleFunc("/readiness", manager.handleReadinessCheck)
	http.HandleFunc("/debug/unhealthy", manager.handleDebugUnhealthy)
	http.HandleFunc("/debug/status", manager.handleDebugStatus)

	http.HandleFunc("/debug/metric/", manager.handleDebugMetric)
	// prometheus custom metric
	manager.CustomMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ping_pong_sample",
			Name:      "pongness",
			Help:      "Custom metric that describes Pong VNF",
		})
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(manager.CustomMetric)

	go manager.publishMetrics()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", manager.vnfManagerPort), nil))
}

