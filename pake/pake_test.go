package pake

import (
	bin "encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test_Encode(t *testing.T) {
	req := LoginReq{111, 222, "333"}
	id := LoginId
	msgI := &Messages{}
	msgI.Context.Id = id
	msgI.Context.SetUserId("sxz")
	msgI.Context.SetSess("1111111")
	msgI.Context.SetSeq(1)
	fmt.Println("Test_Encode lenn", bin.Size(msgI.Context))

	b, _ := json.Marshal(req)
	buf := msgI.Encode(b)
	//fmt.Println(buf)

	pake := msgI.Decode(buf)
	fmt.Println(id, pake.pakeLen, pake.sessionLen, pake.GetSession(), pake.pakeBody)
	var reqj LoginReq
	json.Unmarshal(pake.pakeBody, &reqj)
	fmt.Println("-\n", reqj)
}
func Test_Register(t *testing.T) {
	var a LoginReq = LoginReq{1, 1, "1"}
	Register(1, &MessageLogin{Req: a})
	if msgI, ok := MessageMap[1]; !ok {
		fmt.Println("not find")
		return
	} else {
		r := reflect.New(msgI.Type())
		if msgI.Kind() == reflect.Ptr {
			r.Elem().Set(reflect.New(msgI.Elem().Type()))
		}
		m := r.Elem().Interface().(MessageI)
		fmt.Println("-nnnn", reflect.TypeOf(m))
		m.Process()
	}
}

func Test_Messeage(t *testing.T) {
	fmt.Println("----------------------------------------------------")
	id := PakeId(1)

	Register(id, &MessageLogin{})

	login := LoginReq{111111111, 222222222, "33333333"}
	ctx := ContextInfo{}
	ctx.SetSess("sess")
	ctx.SetSeq(7)
	ctx.SetUserId("sxz")
	ctx.Id = id
	mes := &Messages{ctx}
	fmt.Println("mes", mes)

	fmt.Println("len11111", bin.Size(mes.Context))
	//fmt.Println("len-",uint32(unsafe.Sizeof(ctx)))

	b, _ := json.Marshal(login)
	buf := mes.Encode(b)
	fmt.Println(buf)

	fmt.Println("-----------buf", buf)

	mesr := Xcode(buf)
	if mesr == nil {
		fmt.Println("mesr")
	}
	fmt.Println("mesr",mesr)
}
