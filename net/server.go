package net

import (
	"fmt"
	"github.com/funny/link"
	"github.com/Bulesxz/go/codec"	
)



type Server struct {
	addr string
	server *link.Server
	messageCallback MessageCallback
	connectCallback ConnectCallback
	closeCallback   CloseCallback
	start          int32
}

type ServerHandler struct{
	
}
func (this *ServerHandler)NewServer(addr string) *Server {
	return &Server{
		addr:            addr,
		messageCallback: this.OnMessage,
		connectCallback: this.OnConnection,
		closeCallback:   this.OnClose,
	}
}

func (this *ServerHandler) OnMessage(msg []byte){
	fmt.Println("OnMessage")
}

func (this *ServerHandler) OnConnection(sess *link.Session) *Connection{
	fmt.Println("OnConnection",sess.Id())
	return &Connection{sess}
}
func (this *ServerHandler) OnClose(){
	fmt.Println("OnClose")
}

func (this *Server) Start()  {
	srv, err := link.Serve("tcp", this.addr, codec.GetJsonIoCodec())
	this.server=srv
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
		go func(){
			conn:= this.connectCallback(session) //此处考虑池化
			for{
				var msg []byte
				err = conn.Receive(msg)
				if err!=nil{
					fmt.Println(" session.Receive err|", err)
					this.closeCallback()
					return 
				}
				go this.messageCallback(msg)
			}
		}()
	}
}

func (this *Server) stop(){
	this.start=0;
}
