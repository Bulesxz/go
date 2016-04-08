package net
import (
	"encoding/json"
	"sync"
	log "github.com/Bulesxz/go/logger"
	"github.com/Bulesxz/go/pake"
	"github.com/funny/link"
)

type RcClient struct {
	*link.Session
	protocol string
	addr     string
	recvBuf  map[uint64]chan interface{}
	errChan  map[uint64]chan error
	mesPool *sync.Pool
	sync.Mutex
}


var seq uint64=0

func NewSession() *pake.ContextInfo{
	ctx := &pake.ContextInfo{}
	ctx.SetSess("session")
	ctx.SetUserId("125222")
}

func NewPake(mes *pake.Messages,req interface{},pakeid pake.PakeId) (*pake.Messages,[]byte,uint64 ){
	id := atomic.AddUint64(&seq,1)
	mes.Context.SetId(pakeid)
	mes.Context.SetSeq(id)
	b, err:= json.Marshal(req)
	if err!=nil{
		log.Error(err)
	}
	buf := mes.Encode(b)
	return mes,buf,id
}



func NewRcClient(protocol, addr string) *RcClient{
	&RcClient{
		protocol:protocol,
		addr:addr,
		recvBuf:make(map[uint64]interface{}),
		errChan:make(map[uint64]chan error),
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
	return err
}

func (this *RcClient)  Call(mes *pake.Messages,req interface{},pakeid pake.PakeId) (recvData interface{}, err error) {

	closeChan := make(chan bool, 1)
	GloablTimingWheel.Add(timeout, func() {
		select {
		case <-closeChan: //正常关闭
			//fmt.Println("closeChan")
			return
		default: //超时 干掉连接
			log.Debug("timeout")
			this.errChan <- nil
			close(this.recvBuf)
			this.Close()
		}
	})
	
	
	recvBuf:=make(chan interface{})
	this.Lock()
	if _, ok := this.recvBuf[p.GetSession().Seq]; ok {
		log.Error("[chanrpc] repeated seq ", p.GetSession().Seq)
		this.Unlock()
		return nil,"[chanrpc] repeated seq"
	} else {
		this.recvBuf[p.GetSession().Seq] = recvBuf
	}
	this.Unlock()

	msg,err:=json.Marshal(req)
	
	if err!=nil{
		log.Error(err)
		return nil ,err
	}

	sendData:=mes.Encode(msg)
	
	
	err:=this.Send(sendData)
	if err!=nil{
		log.Error("Send err|",err)
	}
	
	recvBuf,ok = <- this.recvBuf[mes.Context.Seq]
	
	if !ok{
		
	}
	
	
}

func (this *RcClient) run() {
	for {
		//fmt.Println("recvBuf。。。",recvBuf)
		var receiveBuf []byte
		err = this.Receive(&receiveBuf)
		if err != nil {
			log.Error("this.Receive err|", err)
			errChan <- err
			session.Close()
			break
		}
		errChan <- nil
		
		mes:=this.mesPool.Get().(*pake.Messages)
		p:=mes.Decode(receiveBuf)
		
		recvBuf:=make(chan interface{})
		
		this.Lock()
		if _, ok := this.recvBuf[p.GetSession().Seq]; ok {
			err = log.Error("[chanrpc] repeated seq, seq=%v", p.GetSession().Seq)
		} else {
			this.recvBuf[p.GetSession().Seq] = recvBuf
		}
		recvBuf<-p

		this.Unlock()
	}
}