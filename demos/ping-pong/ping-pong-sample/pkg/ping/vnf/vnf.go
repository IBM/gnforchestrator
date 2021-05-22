package vnf

import (
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/ping/config"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/ping/vnf/management"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/ping/vnf/function"
)

type PingVnf struct {
	managementChannel  chan *config.PingConfig
	managementInterface *management.PingVnfManagement
	functionalInterface *function.PingVnfFunction
}

func (vnf *PingVnf) Start() {
	vnf.managementChannel = make(chan *config.PingConfig, 0)
	// management channel should be shared by management interface and functional interface
	// management interface gets requests from vnf manager and updates config.
	// functional interface has to read the config to use it.

	vnf.managementInterface = &management.PingVnfManagement{}
	vnf.functionalInterface = &function.PingVnfFunction{}

	go vnf.managementInterface.Start(vnf.managementChannel)
	vnf.functionalInterface.Start(vnf.managementChannel)
}
