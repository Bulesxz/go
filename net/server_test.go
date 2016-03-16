package net
import (
	"testing"
)

func Test_Encode(t *testing.T){
	
	serh:=&ServerHandler{}
	ser :=serh.NewServer("127.0.0.1:9000")
	ser.Start()
}