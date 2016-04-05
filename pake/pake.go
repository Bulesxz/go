// pakeage
package pake

import (
	"bytes"
	"encoding/json"
//	"fmt"
	log "github.com/Bulesxz/go/logger"
	"github.com/funny/binary"
	"reflect"
)

var (
	MessageMap map[PakeId]reflect.Value = make(map[PakeId]reflect.Value)
)

type PakeId uint32

type Pake struct {
	pakeLen    uint32
	sessionLen uint32
	session    ContextInfo
	pakeBody   []byte
}

func (this *Pake) GetBody() []byte {
	return this.pakeBody
}
func (this *Pake) GetLen() uint32 {
	return this.pakeLen
}
func (this *Pake) GetSession() *ContextInfo {
	return &this.session
}

type ContextInfo struct {
	Id     PakeId `json:"Id"`
	Seq    uint64 `json:"Seq"`
	UserId string `json:"UserId"`
	Sess   string `json:"Sess"`
}

func (this *ContextInfo) SetUserId(userId string) {
	this.UserId = userId
}
func (this *ContextInfo) SetSeq(seq uint64) {
	this.Seq = seq
}
func (this *ContextInfo) SetSess(sess string) {
	this.Sess = sess
}
func (this *ContextInfo) SetId(id PakeId) {
	this.Id = id
}

type Messages struct {
	Context ContextInfo
}

type MessageI interface {
	Init(*ContextInfo)
	Process()
	GetReq() interface{}
	GetRsp() interface{}
}

//PakeId 1

func (this *Messages) marshal(pake *Pake) []byte {
	//atomic.AddUint64(&(this.Context.seq),1)
	var buff []byte = make([]byte, 4)
	var buffer bytes.Buffer

	//session
	session_buff, err := json.Marshal(&pake.session)
	if err != nil {
		log.Error("json.Marshal err", err)
	}

	//pakeLen
	pake.pakeLen = uint32(4 + len(session_buff) + len(pake.pakeBody))
	binary.PutUint32LE(buff, pake.pakeLen)
	buffer.Write(buff)

	//sessionLen
	pake.sessionLen = uint32(len(session_buff))
	binary.PutUint32LE(buff, pake.sessionLen)
	buffer.Write(buff)

	//session
	buffer.Write(session_buff)

	//pakeBody
	buffer.Write(pake.pakeBody)

	return buffer.Bytes()
}

func (this *Messages) unmarshal(msg []byte) *Pake {

	buffer := bytes.NewBuffer(msg)

	var buff []byte = make([]byte, 4)
	pake := new(Pake)

	//pakeLen
	n, err := buffer.Read(buff)
	if err != nil || n != 4 {
		log.Error("len not equal 4 err|", err)
		return nil
	}
	pake.pakeLen = binary.GetUint32LE(buff)

	//sessionLen
	n, err = buffer.Read(buff)
	if err != nil || n != 4 {
		log.Error("len not equal err|", err)
		return nil
	}
	pake.sessionLen = binary.GetUint32LE(buff)

	//session
	var session_buff []byte = make([]byte, pake.sessionLen)
	buffer.Read(session_buff)
	err = json.Unmarshal(session_buff, &pake.session)
	if err != nil {
		log.Error("json.Unmarshal", err)
	}

	//pakeBody
	pake.pakeBody = buffer.Bytes()

	return pake
}

func (this *Messages) Encode(msg []byte) []byte {
	if msg == nil {
		return nil
	}
	pake := Pake{}
	pake.session = this.Context
	pake.pakeBody = msg

	return this.marshal(&pake)
}

func (this *Messages) Decode(msg []byte) *Pake {
	if msg == nil {
		return nil
	}
	pake := this.unmarshal(msg)
	return pake
}

func (this *Messages) Init(sess string) {
	//this.Context.SetSess(sess)
	//this.Context.SetUserId("sxz")
}

func Register(id PakeId, msqI MessageI) {
	log.Debug("id", id, "Register", reflect.TypeOf(msqI))
	MessageMap[id] = reflect.ValueOf(msqI)
}
