package vnf

import (
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/pong/config"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/pong/vnf/management"
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/pong/vnf/function"
)

type PongVnf struct {
	managementChannel  chan *config.PongConfig
	managementInterface *management.PongVnfManagement
	functionalInterface *function.PongVnfFunction
}

func (vnf *PongVnf) Start() {
	vnf.managementChannel = make(chan *config.PongConfig, 0)
	// management channel should be shared by management interface and functional interface
	// management interface gets requests from vnf manager and updates config.
	// functional interface has to read the config to use it.

	vnf.managementInterface = &management.PongVnfManagement{}
	vnf.functionalInterface = &function.PongVnfFunction{}

	go vnf.managementInterface.Start(vnf.managementChannel)
	vnf.functionalInterface.Start(vnf.managementChannel)
}
