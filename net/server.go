package net

import (
	"github.com/funny/binary"
	"fmt"
	"github.com/funny/link"
	"github.com/Bulesxz/go/codec"
	"github.com/Bulesxz/go/pake"
	"reflect"
	"encoding/json"
)
var (
	ctx pake.ContextInfo
	mes *pake.Messages
)


func init(){
	fmt.Println("init")
	pake.Register(1,&pake.MessageLogin{})
	ctx =pake.ContextInfo{}
	ctx.SetSess(nil)
	ctx.SetUserId(9999)
	mes =&pake.Messages{ctx}
}


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
		start:1,
	}
}

func (this *ServerHandler) OnMessage(conn *Connection,msg []byte){
	//fmt.Println("OnMessage",msg)
	if msg  == nil {
		return 
	}
	r:=binary.NewBufferReader(msg[:8]);
	_=r.ReadUint32LE()
	pakeid := r.ReadUint32LE()
	fmt.Println(pake.MessageMap,"pakeid",pakeid)
	if msgI,ok:=pake.MessageMap[pake.PakeId(pakeid)];!ok {
		fmt.Println("not find")
		return 
	}else{
		
		r := reflect.New(msgI.Type())
		if msgI.Kind() == reflect.Ptr {
			r.Elem().Set(reflect.New(msgI.Elem().Type()))
		}
		t:=r.Elem().Interface().(pake.MessageI)
		p:=mes.Decode(msg)
		//fmt.Println(reflect.TypeOf(p))
		json.Unmarshal(p.GetBody(),t.GetReq())
		fmt.Println(reflect.TypeOf(t))
		
		t.Process()
		
		b,_:=json.Marshal(t.GetRsp())
		buf:=mes.Encode(pake.PakeId(pakeid),b)
		//fmt.Println(buf)
		conn.Send(buf)
	}
	
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
		fmt.Println("Accept")
		go func(session *link.Session){
			conn:= this.connectCallback(session) //此处考虑池化
			for{
				var msg []byte
				err = conn.Receive(&msg)
				if err!=nil{
					fmt.Println(" session.Receive err|", err)
					this.closeCallback()
					return 
				}
				//fmt.Println("Receive")
				go  this.messageCallback(conn,msg)
			}
		}(session)
	}
}

func (this *Server) stop(){
	this.start=0;
}
