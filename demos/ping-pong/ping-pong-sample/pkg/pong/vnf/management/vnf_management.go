package management

import (
	"os"
	"log"
	"fmt"
	"context"
	"net"
	"google.golang.org/grpc"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/protos"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/pong/config"
)

type PongVnfManagement struct {
	managementPort string
	managementChannel chan <- *config.PongConfig
	isConfigured bool
}

func (m *PongVnfManagement) Configure(_ context.Context, request *protos.PongVnfConfigRequest) (*protos.PongVnfConfigResponse, error) {
	log.Printf("configure called with body {%s}", request.Message)
	pongConfig := &config.PongConfig {
		Message: request.Message,
	}
	m.managementChannel <- pongConfig
	m.isConfigured = true
	vnfConfigResponse := &protos.PongVnfConfigResponse{}
	return vnfConfigResponse, nil
}

func (m *PongVnfManagement) ReadinessCheck(_ context.Context, request *protos.ReadinessRequest) (*protos.ReadinessResponse, error) {
	log.Println("readiness check called")
	readinessResponse := &protos.ReadinessResponse{ Ready:m.isConfigured }
	return readinessResponse, nil
}

func (m *PongVnfManagement) Start(managementChannel chan <- *config.PongConfig) { // channel that is used only for writing objects to it
	m.managementChannel = managementChannel
	m.isConfigured = false
	m.managementPort = os.Getenv("VNF_MANAGEMENT_PORT")
	if m.managementPort == "" {
		log.Fatal("'the expected arguments are not set in the env vars")
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", m.managementPort))
	if err != nil {
		log.Fatalf("Error - %s", err)
	}
	log.Printf("VNF management Server is listening on port %s ...", m.managementPort)

	vnfServer := grpc.NewServer()
	protos.RegisterPongVnfServer(vnfServer, m)

	vnfServer.Serve(listener)
}
