package main

import (
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/pong/vnf-manager"
)

func main() {
	pongVnfManager := &vnf_manager.PongVnfManager{}
	pongVnfManager.Start()
}
