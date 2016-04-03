// pakeage
package pake
import (
	"fmt"
	"encoding/json"
	"github.com/funny/binary"
	"bytes"
	"reflect"
	log "github.com/Bulesxz/go/logger"
	bin "encoding/binary"
)

var (
	MessageMap map[PakeId]reflect.Value = make(map[PakeId]reflect.Value)
)
type  PakeId uint32

type Pake struct{
	pakeLen uint32
	sessionLen uint32
	session ContextInfo
	pakeBody []byte
}

func (this *Pake) GetBody() []byte{
	return this.pakeBody
}
func (this *Pake) GetLen() uint32{
	return this.pakeLen
}
func (this *Pake) GetSession() *ContextInfo{
	return &this.session
}

type ContextInfo struct{
	Id	PakeId `json:"Id"`
	Seq uint64 `json:"Seq"`
	UserId	string `json:"UserId"`
	Sess	string `json:"Sess"`
}

func (this *ContextInfo) SetUserId(userId string){
	this.UserId=userId
}
func (this *ContextInfo) SetSeq(seq uint64){
	this.Seq = seq
}
func (this *ContextInfo) SetSess(sess string){
	this.Sess=sess
}
type Messages struct{
	Context ContextInfo
}

type MessageI interface{
	Init(*ContextInfo)
	Process()
	GetReq() interface{}
	GetRsp() interface{}
}


//PakeId 1



func (this *Messages) marshal(pake *Pake) []byte{
	//atomic.AddUint64(&(this.Context.seq),1)
	var buff []byte=make([]byte,4)
	var buffer bytes.Buffer
	
	//pakeLen
	binary.PutUint32LE(buff,pake.pakeLen)
	buffer.Write(buff)
	
	//sessionLen
	binary.PutUint32LE(buff,pake.sessionLen)
	buffer.Write(buff)
	
	//session
	session_buff,err:= json.Marshal(&pake.session)
	if err!=nil{
		log.Error("json.Marshal err",err)
	}
	buffer.Write(session_buff)
	
	//pakeBody
	buffer.Write(pake.pakeBody)
	
	fmt.Println("marshal",pake.session,"session_buff",session_buff,"len",len(session_buff),"body",pake.pakeBody)
	
	return buffer.Bytes()
}

func (this *Messages) unmarshal(msg []byte) *Pake{
	
	buffer:=bytes.NewBuffer(msg)
	
	var buff []byte=make([]byte,4)
	pake:=new(Pake)
	
	//pakeLen
	n,err:=buffer.Read(buff)
	if err!=nil || n != 4 {
		log.Error("len not equal 4 err|",err)
		return nil
	}
	pake.pakeLen=binary.GetUint32LE(buff)
	fmt.Println("pake.pakeLen",pake.pakeLen)
	
	//sessionLen
	n,err=buffer.Read(buff)
	if err!=nil || n !=4 {
		log.Error("len not equal err|",err)
		return nil
	}
	pake.sessionLen=binary.GetUint32LE(buff)
	fmt.Println("pake.sessionLen",pake.sessionLen)
	
	//session
	var session_buff []byte=make([]byte,pake.sessionLen)
	buffer.Read(session_buff)
	fmt.Println("session_buff",session_buff,"\n")
	err =json.Unmarshal(session_buff,&pake.session)
	if err!=nil{
		log.Error("json.Unmarshal",err)
	}
	fmt.Println("pake.session",pake.session)
	
	//pakeBody
	pake.pakeBody=buffer.Bytes()
	
	return pake
}

func (this *Messages) Encode(msg []byte)([]byte){
	if msg == nil {
		return nil
	}
	pake:=Pake{}
	pake.session=this.Context
	//pake.sessionLen= uint32(unsafe.Sizeof(this.Context)) 
	
	pake.sessionLen= 1 //bin.Size(pake.session)
	fmt.Println("len",bin.Size(pake.session))
	return nil
	
	pake.pakeBody=msg
	pake.pakeLen=uint32(len(msg)) + 4 + pake.sessionLen 
	fmt.Println("encode",pake.pakeLen,pake.sessionLen,pake.session,pake.pakeBody)
	return this.marshal(&pake)
}

func (this *Messages) Decode(msg []byte) *Pake{
	if msg == nil {
		return nil
	}
	pake :=this.unmarshal(msg)
	return pake
}

func (this *Messages) Init(sess string){
	this.Context.SetSess(sess)
	this.Context.SetUserId("sxz")
}


func Register(id PakeId,msqI MessageI){
	MessageMap[id]=reflect.ValueOf(msqI)
}