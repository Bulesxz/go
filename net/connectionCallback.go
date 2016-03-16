package net

import (
	"fmt"
	"github.com/funny/link"
)
type Connection struct{
	*link.Session
}

type MessageCallback func(msg []byte)
type ConnectCallback func(sess *link.Session)
type CloseCallback func()

func OnConnection(sess *link.Session) *Connection{
	fmt.Println("OnConnection",sess.Id())
	return &Connection{sess}
}

func (this *Connection) OnMessage(msg []byte){
	fmt.Println("OnMessage")
}

func (this *Connection) OnClose(){
	fmt.Println("OnClose")
	this.Close()
}