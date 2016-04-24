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
	closeChan	 map[uint64]chan struct{}//通知recvbuf关闭了,不能写入
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
		closeChan:make(map[uint64]chan struct{}),
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

	id := mes.Context.Seq
	closeChan := make(chan struct{}, 1)	
	recvBuf:=make(chan interface{},1)
	
	this.Lock()
	if _, ok := this.recvBuf[id]; ok {
		log.Error("[chanrpc] repeated seq ", id)
		this.Unlock()
		return nil,fmt.Errorf("%s","[chanrpc] repeated seq")
	} else {
		this.recvBuf[id] = recvBuf
		this.closeChan[id] = closeChan
	}
	this.Unlock()
	
	
	msg,err:=json.Marshal(req)
	if err!=nil{
		log.Error(err)
		closeChan <-struct{}{}
		return nil ,err
	}

	sendData:=mes.Encode(msg)
	//fmt.Println("sendData")
	
	
	GloablTimingWheel.Add(timeout, func() {
		select {
		case <-closeChan: //正常关闭  //清扫
			//fmt.Println("closeChan")
			close(recvBuf)
			this.Lock()
			delete(this.closeChan, id)
			delete(this.recvBuf, id)
			this.Unlock()
			
		default: //超时 干掉连接
			log.Warning("timeout")
			fmt.Println("timeout.......")
		
			close(recvBuf)
			
			this.Lock()
			delete(this.recvBuf, id)
			delete(this.closeChan, id)
			this.Unlock()
			
			//this.Close()
		}
	})
	
	err=this.Send(sendData)
	if err!=nil{
		close(closeChan)
		log.Error("Send err|",err)
		return nil,err
	}
	
	
	select {//防止没有err 阻塞
		case err = <- this.errChan:
		if err!=nil {
			close(closeChan) //
			fmt.Println("err:",err)
			log.Error(err)
			return nil,err
		}
		default:
		//	fmt.Println("no err")
	}

	 defer func(){     //必须要先声明defer，否则不能捕获到panic异常
                                if err := recover(); err != nil {
                                        log.Error("panic",err)    //这里的err其实就是panic传入的内容
                                }
         }()
	
	this.Lock()
	receiveBuf,ok:= this.recvBuf[id]
	this.Unlock()
	if !ok{
		fmt.Println("receiveBuf is delete")
		log.Warning("receiveBuf is delete")
		err = fmt.Errorf("receiveBuf is delete")
		return nil,err
	}
	recvData,ok= <-receiveBuf
	if !ok{
		fmt.Println("!ok:",ok)
		log.Error("!ok")
		err = fmt.Errorf("recvBuf is close")
		return nil,err
	}
	close(closeChan)
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
		
		go func(){	
			mes:=this.mesPool.Get().(*pake.Messages)
			p:=mes.Decode(receiveBuf)
			
			defer func(){     //必须要先声明defer，否则不能捕获到panic异常
				if err := recover(); err != nil {
					fmt.Println("panic",err)    //这里的err其实就是panic传入的内容
				}
			}()
			this.Lock();	
			closeChan,ok:=  this.closeChan[p.GetSession().Seq]
			this.Unlock();
			if !ok{
				fmt.Println("closeChan is delete")
			} else {
				this.Lock();
				recvChan,ok:= this.recvBuf[p.GetSession().Seq];
				this.Unlock();
				if !ok{
					fmt.Println("recvChan is delete")
				}else {	
					select{
					case <- closeChan :
						fmt.Println("closeChan is close")
					default:
						recvChan <- p
					}
				}
			}
			this.mesPool.Put(mes)
		}()
			
	}
}
