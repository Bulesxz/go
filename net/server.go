package net

import (
	"fmt"
	"github.com/Bulesxz/go/codec"
	log "github.com/Bulesxz/go/logger"
	"github.com/Bulesxz/go/pake"
	"github.com/funny/link"
)

type Task struct {
	messageCallback MessageCallback
	conn            *Connection
	msg             []byte
}
type Server struct {
	addr            string
	server          *link.Server
	messageCallback MessageCallback
	connectCallback ConnectCallback
	closeCallback   CloseCallback
	start           int32
	tasks           chan *Task
}

type ServerHandler struct {
}

func (this *ServerHandler) NewServer(addr string) *Server {
	return &Server{
		addr:            addr,
		messageCallback: this.OnMessage,
		connectCallback: this.OnConnection,
		closeCallback:   this.OnClose,
		start:           1,
		tasks:           make(chan *Task, 1024), ///改进
	}
}

func (this *ServerHandler) OnMessage(conn *Connection, msg []byte) {
	log.Debug("OnMessage")
	//fmt.Println("Onmessage")
	if msg == nil {
		log.Debug("mgs == nil")
		//fmt.Println("Onmessage111111")
		return
	}
	rsp := pake.Deal(msg)
	err := conn.Send(rsp)
	if err != nil {
		log.Error(err)
	}
	fmt.Println("send rsp",rsp)
}

func (this *ServerHandler) OnConnection(sess *link.Session) *Connection {
	log.Debug("OnConnection", sess.Id())
	return &Connection{sess}
}
func (this *ServerHandler) OnClose() {
	log.Debug("OnClose")
}

func (this *Server) Start() {
	log.Debug("server start....")
	srv, err := link.Serve("tcp", this.addr, codec.GetJsonIoCodec())
	this.server = srv
	if err != nil {
		log.Error("link.Serve err|", err)
		return
	}

	go this.Dotask()

	for this.start == 1 {
		session, err := srv.Accept()
		if err != nil {
			log.Error("srv.Accept err|", err)
			return
		}
		//fmt.Println("Accept")
		go func(session *link.Session) {
			conn := this.connectCallback(session) //此处考虑池化
			for {
				var msg []byte
				err = conn.Receive(&msg)
				log.Debug("Receive",msg,err)
				if err != nil {
					log.Debug(" session.Receive err|", err)
					this.closeCallback()
					return
				}
				go this.messageCallback(conn, msg)
				//this.tasks <- &Task{this.messageCallback,conn,msg}
				//		fmt.Println("Receive")
			}
		}(session)
	}
}

func (this *Server) stop() {
	this.start = 0
}

func (this *Server) Dotask() {
	for {
		task := <-this.tasks
		task.messageCallback(task.conn, task.msg)
	}
}
