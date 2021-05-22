package management

import (
	"os"
	"log"
	"fmt"
	"context"
	"net"
	"google.golang.org/grpc"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/protos"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/ping/config"
)

type PingVnfManagement struct {
	managementPort string
	managementChannel chan <- *config.PingConfig
	isConfigured bool
}

func (m *PingVnfManagement) Configure(_ context.Context, request *protos.PingVnfConfigRequest) (*protos.PingVnfConfigResponse, error) {
	log.Printf("configure called with {%s %s}", request.PongAddress, request.PongPort)
	pingConfig := &config.PingConfig {
		PongAddress: request.PongAddress,
		PongPort: request.PongPort,
	}
	m.managementChannel <- pingConfig
	m.isConfigured = true
	vnfConfigResponse := &protos.PingVnfConfigResponse{}
	return vnfConfigResponse, nil
}

func (m *PingVnfManagement) ReadinessCheck(_ context.Context, request *protos.ReadinessRequest) (*protos.ReadinessResponse, error) {
	log.Println("readiness check called")
	readinessResponse := &protos.ReadinessResponse{ Ready:m.isConfigured }
	return readinessResponse, nil
}

func (m *PingVnfManagement) Start(managementChannel chan <- *config.PingConfig) { // channel that is used only for writing objects to it
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
	protos.RegisterPingVnfServer(vnfServer, m)

	vnfServer.Serve(listener)
}
