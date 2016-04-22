package net
import (
	"fmt"
	"sync/atomic"
	"encoding/json"
	"sync"
	log "github.com/Bulesxz/go/logger"
	"github.com/Bulesxz/go/pake"
	"github.com/funny/link"
	"time"
	"github.com/Bulesxz/go/codec"
)

type RcClient struct {
	*link.Session
	protocol string
	addr     string
	recvBuf  map[uint64]chan interface{}
	exit	 map[uint64]chan bool//通知recvbuf关闭了,不能写入
	errChan  chan error
	mesPool *sync.Pool
	sync.Mutex
}


var seq uint64=0

func NewSession() *pake.ContextInfo{
	ctx := &pake.ContextInfo{}
	ctx.SetSess("session")
	ctx.SetUserId("125222")
	return ctx
}

func InitPake(mes *pake.Messages,pakeid pake.PakeId) (uint64 ){
	id := atomic.AddUint64(&seq,1)
	mes.Context.SetId(pakeid)
	mes.Context.SetSeq(id)
	return id
}



func NewRcClient(protocol, addr string) *RcClient{
	return &RcClient{
		protocol:protocol,
		addr:addr,
		recvBuf:make(map[uint64]chan interface{}),
		exit:make(map[uint64]chan bool),
		errChan:make(chan error),
		mesPool:&sync.Pool{
        	New: func() interface{} {
        	    return &pake.Messages{}
        },
    }}
}

func (this *RcClient) ConnetcTimeOut(timeout time.Duration) error {
	session, err := link.ConnectTimeout(this.protocol, this.addr, timeout, codec.GetJsonIoCodec())
	if err != nil {
		//fmt.Println("link.ConnectTimeout err|", err)
		log.Error("link.ConnectTimeout err|", err)
		return err
	}
	this.Session = session
	go this.run()
	return err
}

func (this *RcClient)  Call(mes *pake.Messages,req interface{},timeout time.Duration) (recvData interface{}, err error) {

	closeChan := make(chan bool, 1)	
	
	recvBuf:=make(chan interface{},1)
	exit:=make(chan bool,1)
	
	this.Lock()
	if _, ok := this.recvBuf[mes.Context.Seq]; ok {
		log.Error("[chanrpc] repeated seq ", mes.Context.Seq)
		this.Unlock()
		return nil,fmt.Errorf("%s","[chanrpc] repeated seq")
	} else {
		this.recvBuf[mes.Context.Seq] = recvBuf
		this.exit[mes.Context.Seq] = exit
	}
	this.Unlock()
	
	GloablTimingWheel.Add(timeout, func() {
		select {
		case <-closeChan: //正常关闭
			//fmt.Println("closeChan")
			return
		default: //超时 干掉连接
			log.Debug("timeout")
			fmt.Println("timeout.......")
			this.errChan <- fmt.Errorf("timeout")
			close(exit)
			close(recvBuf)
			//this.Close()
		}
	})
	
	msg,err:=json.Marshal(req)
	
	if err!=nil{
		log.Error(err)
		return nil ,err
	}

	sendData:=mes.Encode(msg)
	//fmt.Println("sendData")
	exit <- false
	err=this.Send(sendData)
	if err!=nil{
		closeChan <-true
		log.Error("Send err|",err)
		return recvData,err
	}
	
	err= <- this.errChan
	if err!=nil{
		fmt.Println("err:",err)
		log.Error(err)	
	}
		
	this.Lock()
	recvData,ok := <- this.recvBuf[mes.Context.Seq]
	if !ok{
		//fmt.Println("!ok:",ok)
		log.Error("!ok")
		err = fmt.Errorf("timeout")
	}
	this.Unlock()
	
	closeChan <-true
	return recvData,err
	
}

func (this *RcClient) run() {
	for {
		var receiveBuf []byte
		err := this.Receive(&receiveBuf)
		//fmt.Println("1",receiveBuf)
		if err != nil {
			log.Error("this.Receive err|", err)
			this.errChan <- err
			this.Close()
			break
		}
		
		this.errChan <- nil
		
		mes:=this.mesPool.Get().(*pake.Messages)
		p:=mes.Decode(receiveBuf)
		
		/*recvBuf:=make(chan interface{})
		this.Lock()
		if _, ok := this.recvBuf[p.GetSession().Seq]; ok {
			err = log.Error("[chanrpc] repeated seq, seq=%v", p.GetSession().Seq)
		} else {
			this.recvBuf[p.GetSession().Seq] = recvBuf
		}*/
		
		defer func(){     //必须要先声明defer，否则不能捕获到panic异常
			if err := recover(); err != nil {
				//fmt.Println(err)    //这里的err其实就是panic传入的内容
			}
		}()
		
		this.Lock()
		_,ok:= <- this.exit[p.GetSession().Seq]
		if ok{ //没有超时，则没有close ,所以可写
			//fmt.Println("p:",p)
			this.recvBuf[p.GetSession().Seq] <- p
		}
		this.Unlock()
	}
}