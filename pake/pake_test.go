package pake
import (
	"testing"
	"fmt"
	"encoding/json"
	"reflect"
)

func Test_Encode(t *testing.T){
	req:=LoginReq{1,2,"3"}
	id:=PakeId(1)
	msgI:=&Messages{}
	msgI.Context.seq=1
	b,_:=json.Marshal(req)
	buf:=msgI.Encode(id,b)
	fmt.Println(buf)
	
	pake:=msgI.Decode(buf)
	fmt.Println(pake.id,pake.pakeLen,pake.pakeBody)
	var reqj LoginReq
	json.Unmarshal(pake.pakeBody,&reqj)
	fmt.Println(reqj)
}
func Test_Register(t *testing.T){
	Register(1, &MessageLogin{})
}

func Test_Messeage(t *testing.T){
	Register(1,&MessageLogin{})
	
	login:=LoginReq{1,2,"sxz"}
	ctx :=ContextInfo{}
	ctx.SetSess(nil)
	ctx.SetUserId(9999)
	mes :=&Messages{ctx}
	
	
	id:=PakeId(1)
	b,_:=json.Marshal(login)
	buf:=mes.Encode(id,b)
	fmt.Println(login,mes,buf)
	
	fmt.Println("-----------")
	
	if msgI,ok:=MessageMap[1];!ok {
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
		json.Unmarshal(p.GetBody(),t.GetReq())
		fmt.Println(reflect.TypeOf(t))
		
		t.Process()
	}
}