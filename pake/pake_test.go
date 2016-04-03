package pake
import (
	"testing"
	"fmt"
	"encoding/json"
	"reflect"
	bin "encoding/binary"
)

func Test_Encode(t *testing.T){
	req:=LoginReq{111,222,"333"}
	id:=PakeId(1)
	msgI:=&Messages{}
	msgI.Context.Id=1
	msgI.Context.SetUserId("sxz")
	msgI.Context.SetSess("1111111")
	msgI.Context.SetSeq(1)
	b,_:=json.Marshal(req)
	buf:=msgI.Encode(b)
	//fmt.Println(buf)
	
	pake:=msgI.Decode(buf)
	fmt.Println(id,pake.pakeLen,pake.sessionLen,pake.GetSession(),pake.pakeBody)
	var reqj LoginReq
	json.Unmarshal(pake.pakeBody,&reqj)
	fmt.Println("-\n",reqj)
}
func Test_Register(t *testing.T){
	Register(1, &MessageLogin{})
}

func Test_Messeage(t *testing.T){
	fmt.Println("----------------------------------------------------")
	id:=PakeId(8)
	
	Register(id,&MessageLogin{})
	
	login:=LoginReq{111111111,222222222,"33333333"}
	ctx :=ContextInfo{}
	ctx.SetSess("sess")
	ctx.SetSeq(7)
	ctx.SetUserId("sxz")
	ctx.Id=id
	mes :=&Messages{ctx}
	fmt.Println("mes",mes)
	
	fmt.Println("len11111",bin.Size(mes.Context))
	//fmt.Println("len-",uint32(unsafe.Sizeof(ctx)))
	
	b,_:=json.Marshal(login)
	buf:=mes.Encode(b)
	fmt.Println(buf)
	
	fmt.Println("-----------buf",buf)
	
	if msgI,ok:=MessageMap[id];!ok {
		fmt.Println("not find")
		return 
	}else{
		
		r := reflect.New(msgI.Type())
		if msgI.Kind() == reflect.Ptr {
			r.Elem().Set(reflect.New(msgI.Elem().Type()))
		}
		t:=r.Elem().Interface().(MessageI)
		p:=mes.Decode(buf)
		//fmt.Println(reflect.TypeOf(p))
		t.Init(&p.session)
		json.Unmarshal(p.GetBody(),t.GetReq())
		fmt.Println(reflect.TypeOf(t))
		
		t.Process()
	}
}