package main

import (
	"fmt"
	log "github.com/Bulesxz/go/logger"
	"github.com/Bulesxz/go/net"
)

func main() {
	fmt.Println("main")
	//log.Init()
	log.Info("-------")
	serh := &net.ServerHandler{}
	log.Info("serh.NewServer")
	ser := serh.NewServer("127.0.0.1:9000")
	ser.Start()
}
