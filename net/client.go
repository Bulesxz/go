package net

import(
	"github.com/funny/link"
	"github.com/Bulesxz/go/codec"
	timewheel "github.com/Bulesxz/go/time"
	log "github.com/Bulesxz/go/logger"
	"time"
)


var (
	GloablTimingWheel *timewheel.TimingWheel 
)

func init(){
	GloablTimingWheel = timewheel.NewTimingWheel(time.Millisecond*100, 50)
	//
}


type Client struct{
	*link.Session
	protocol string
	addr     string
	context  interface{}
	recvBuf chan []byte
	errChan chan error
}

func NewClient(protocol,addr string) *Client{
	return &Client{protocol:protocol,addr:addr}
}
func (this *Client) ConnetcTimeOut(timeout time.Duration) error{
	session,err := link.ConnectTimeout(this.protocol,this.addr,timeout,codec.GetJsonIoCodec())
	if err!=nil {
		log.Error("link.ConnectTimeout err|",err)
		return err
	}
	this.Session=session
	this.recvBuf = make(chan []byte, 1)
	this.errChan = make(chan error, 1)
	go func(session *link.Session,recvBuf  chan<-  []byte,errChan chan<- error){
		for{
			//fmt.Println("recvBuf。。。",recvBuf)
			var receiveBuf []byte
			err = this.Receive(&receiveBuf)
			if err!=nil {
				log.Error("this.Receive err|",err)
				errChan<-err
				session.Close()
				break
			}
			errChan<-nil
			recvBuf<-receiveBuf
			//fmt.Println("recvBuf。。。",recvBuf)
			//fmt.Println("receive:",receiveBuf)
		}
	}(session,this.recvBuf,this.errChan)
	return err
}

/*
func (this *Client) timeout(){
	fmt.Println("client timeout")
	
	//this.Close()//有问题，这样只要到了timeout时间，连接都会被关掉
	select {
		case <-this.errChan : //正常关闭
			return
		default://超时 干掉连接
			close(this.recvBuf)
			this.Close()
	}
	
}*/
func (this *Client) SendTimeOut(timeout time.Duration,msg interface{}) ( []byte,error){
	recvBuf := this.recvBuf
	errChan := this.errChan
	
	var closeChan chan bool
	closeChan = make(chan bool, 1)
	GloablTimingWheel.Add(timeout,func(){
		select {
		case <-closeChan : //正常关闭
			//fmt.Println("closeChan")
			return
		default://超时 干掉连接
			log.Info("timeout")
			this.errChan<-nil
			close(this.recvBuf)
			this.Close()
		}
	})
	
	err:=this.Send(msg)
	if err!=nil{
		log.Error("this.Send err|",err)
		closeChan<-true //关掉timeout
		return nil, err
	}
	
	err = <-errChan
	if err!=nil{
		close(closeChan)
		return nil ,err
	}
	var receiveBuf []byte
	//fmt.Println("recvBuf+++",recvBuf)
	
	ok := false
	receiveBuf ,ok= <-recvBuf//阻塞等待 ，直到超时
	if !ok { //关闭
		log.Info("!ok 关闭 this.recvBuf")
	}
	closeChan<-true //关掉timeout
	//fmt.Println("recvBuf+++",recvBuf)
	return receiveBuf,err
}










