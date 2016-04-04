package pake

import (
	"fmt"
)

type LoginReq struct {
	A int32  `json:"a"`
	B int32  `json:"b"`
	C string `json:"c"`
}

type LoginRsp struct {
	A int32  `json:"a"`
	B int32  `json:"b"`
	C string `json:"c"`
}

type MessageLogin struct {
	ContextInfo
	Req LoginReq
	Rsp LoginRsp
}

func (this *MessageLogin) Init(c *ContextInfo) {
	this.ContextInfo = *c
	//fmt.Println("init")
}
func (this *MessageLogin) Process() {
	fmt.Println("process", this.GetReq())
	//this.Rsp.A=1
	//this.Rsp.B=1
	//this.Rsp.C="ssssssssss"
	this.Rsp = LoginRsp(this.Req)
}
func (this *MessageLogin) GetReq() interface{} {
	//fmt.Println("GetReq")
	return &this.Req
}
func (this *MessageLogin) GetRsp() interface{} {
	//fmt.Println("GetRsp")
	return &this.Rsp
}

func init() {
	Register(1, &MessageLogin{})
}
