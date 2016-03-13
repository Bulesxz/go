// pakeage
package pake
import (
//	"encoding/json"
	"github.com/funny/binary"
	"github.com/funny/link"
	"bytes"
	"fmt"
	"sync/atomic"
	"reflect"
)

var (
	MessageMap map[PakeId]reflect.Value = make(map[PakeId]reflect.Value)
)
type  PakeId uint32

type Pake struct{
	pakeLen uint32
	id	PakeId 
	pakeBody []byte
}

func (this *Pake) GetId()PakeId{
	return this.id	
}
func (this *Pake) GetBody() []byte{
	return this.pakeBody
}
func (this *Pake) GetLen() uint32{
	return this.pakeLen
}

type ContextInfo struct{
	sess *link.Session
	seq uint64
	userId	uint64
}
func (this *ContextInfo) SetSess(sess *link.Session){
	this.sess=sess
}
func (this *ContextInfo) SetUserId(userId uint64){
	this.userId=userId
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



type LoginReq struct{
	A int32	`json:"a"`
	B int32	`json:"b"`
	C string `json:"c"`
}

type LoginRsp struct{
	A int32	`json:"a"`
	B int32	`json:"b"`
	C string `json:"c"`
}

type MessageLogin struct{
	Req LoginReq
	Rsp LoginRsp
}

func (this *MessageLogin) Init(c *ContextInfo){
	fmt.Println("init")
}
func (this *MessageLogin) Process() {
	fmt.Println("process",this.GetReq())
}
func (this *MessageLogin) GetReq() interface{}{
	fmt.Println("GetReq")
	return &this.Req
}
func (this *MessageLogin) GetRsp () interface{}{
	fmt.Println("GetRsp")
	return &this.Rsp
}

 

func (this *Messages) marshal(pake *Pake) []byte{
	atomic.AddUint64(&(this.Context.seq),1)
	var buff []byte=make([]byte,4)
	var buffer bytes.Buffer
	
	binary.PutUint32LE(buff,pake.pakeLen)
	buffer.Write(buff)
	
	binary.PutUint32LE(buff,uint32(pake.id))
	buffer.Write(buff)
	
	buffer.Write(pake.pakeBody)
	
	return buffer.Bytes()
}

func (this *Messages) unmarshal(msg []byte) *Pake{
	
	buffer:=bytes.NewBuffer(msg)
	var buff []byte=make([]byte,4)
	pake:=new(Pake)
	n,err:=buffer.Read(buff)
	if err!=nil || n != 4 {
		fmt.Println("len not equal 4 err|",err)
		return nil
	}
	pake.pakeLen=binary.GetUint32LE(buff)
	
	n,err=buffer.Read(buff)
	if err!=nil || n !=4 {
		fmt.Println("len not equal err|",err)
		return nil
	}
	pake.id=PakeId(binary.GetUint32LE(buff))
	pake.pakeBody=buffer.Bytes()
	return pake
}

func (this *Messages) Encode(id PakeId,msg []byte) []byte{
	if msg == nil {
		return nil
	}
	pake:=Pake{}
	pake.id=id
	pake.pakeBody=msg
	pake.pakeLen=uint32(len(msg))
	return this.marshal(&pake)
}

func (this *Messages) Decode(msg []byte) *Pake{
	if msg == nil {
		return nil
	}
	pake :=this.unmarshal(msg)
	return pake
}



func Register(id PakeId,msqI MessageI){
	MessageMap[id]=reflect.ValueOf(msqI)
}