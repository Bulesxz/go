package net

import (
	"github.com/funny/link"
)
type Connection struct{
	*link.Session
}

type MessageCallback func(msg []byte)
type ConnectCallback func(sess *link.Session) *Connection
type CloseCallback func()
