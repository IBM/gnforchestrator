package main

import (
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/ping/vnf"
)

func main() {
	pingVnf := &vnf.PingVnf{}
	pingVnf.Start()
}
