package net

import (
	"fmt"
)
import(
	"github.com/funny/link"
	"github.com/go/codec"
	timewheel "github.com/go/time"
	"time"
)


var (
	GloablTimingWheel *timewheel.TimingWheel 
)

func init(){
	GloablTimingWheel = timewheel.NewTimingWheel(time.Second*1, 10)
	//
}


type Client struct{
	*link.Session
	protocol string
	addr     string
	context  interface{}
	recvBuf chan []byte
}

func NewClient(protocol,addr string) *Client{
	return &Client{protocol:protocol,addr:addr}
}
func (this *Client) ConnetcTimeOut(timeout time.Duration) error{
	session,err := link.ConnectTimeout(this.protocol,this.addr,timeout,codec.GetJsonIoCodec())
	if err!=nil {
		fmt.Println("link.ConnectTimeout err|",err)
		return err
	}
	this.Session=session
	this.recvBuf = make(chan []byte, 1)
	go func(session *link.Session,recvBuf  chan<-  []byte){
		for{
			var receiveBuf []byte
			err = this.Receive(&receiveBuf)
			if err!=nil {
				fmt.Println("this.Receive err|",err)
				session.Close()
				break
			}
			recvBuf<-receiveBuf
			fmt.Println("receive:",receiveBuf)
		}
	}(session,this.recvBuf)
	return err
}

func (this *Client) timeout(){
	fmt.Println("client timeout")
	this.Close()
}
func (this *Client) SendTimeOut(timeout time.Duration,msg interface{}) ( []byte,error){
	GloablTimingWheel.Add(timeout,this.timeout)
	err:=this.Send(msg)
	if err!=nil{
		fmt.Println("this.Send err|",err)
		close(this.recvBuf)
		return nil, err
	}
	var receiveBuf []byte
	receiveBuf <- this.recvBuf
	
	close(this.recvBuf)
	return receiveBuf,err
}










