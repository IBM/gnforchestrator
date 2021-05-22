package main

import (
	"github.com/IBM/gnforchestrator/demos/ping-pong/ping-pong-sample/pkg/pong/vnf"
)

func main() {
	pongVnf := &vnf.PongVnf{}
	pongVnf.Start()
}
