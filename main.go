package main
import (
	"fmt"
	"github.com/Bulesxz/go/net"
	log "github.com/Bulesxz/go/logger"
)




func main(){
	fmt.Println("main")
	log.Init()
	serh:=&net.ServerHandler{}
	ser :=serh.NewServer("127.0.0.1:9000")
	ser.Start()
}