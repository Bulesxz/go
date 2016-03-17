package main
import (
	"github.com/Bulesxz/go/net"
)

func main(){
	

	serh:=&net.ServerHandler{}
	ser :=serh.NewServer("127.0.0.1:9000")
	ser.Start()
}