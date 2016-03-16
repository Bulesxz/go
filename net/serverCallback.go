package server

import (
	"fmt"
	"github.com/funny/link"
)
type Connection struct{
	*link.Session
}

type MessageCallback func(conn *Connection, msg []byte)
type ConnectCallback func(conn *Connection)
type CloseCallback func(conn *Connection)

func OnConnection(conn *Connection){
	fmt.Println("OnConnection",conn.Session.Id())
}

func OnMessage(conn *Connection, msg []byte){
	fmt.Println("OnMessage")
}

func OnClose(conn *Connection){
	fmt.Println("OnClose")
	conn.Close()
}