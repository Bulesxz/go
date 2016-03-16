package net

import (
	"fmt"
	"github.com/funny/link"
	"gosxz/codec"	
)



type Server struct {
	*link.Server
	messageCallback MessageCallback
	connectCallback ConnectCallback
	closeCallback   CloseCallback
	start          int32
}

func (this *Server) Start(addr string)  {
	srv, err := link.Serve("tcp", "addr", codec.GetJsonIoCodec())
	if err != nil {
		fmt.Println("link.Serve err|", err)
		return 
	}

	for this.start==1{
		session, err := srv.Accept()
		if err != nil {
			fmt.Println("srv.Accept err|", err)
			return 
		}
		conn:= this.connectCallback(session) //此处考虑池化
		go func(){
			for{
				var msg []byte
				err = conn.Session.Receive(msg)
				if err!=nil{
					fmt.Println(" session.Receive err|", err)
					this.closeCallback()
					return 
				}
				this.messageCallback(msg)
			}
		}()
	}
}

func (this *Server) stop(){
	this.start=0;
}
