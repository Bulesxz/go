package main

import (
	"fmt"
//	log "github.com/Bulesxz/go/logger"
	"github.com/Bulesxz/go/net"
)

func main() {
	fmt.Println("main")
	serh := &net.ServerHandler{}
	ser := serh.NewServer("127.0.0.1:9000")
	ser.Start()
}
