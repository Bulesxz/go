package pake

import (
	"encoding/json"
	log "github.com/Bulesxz/go/logger"
	"reflect"
)

func Xcode(msg []byte) []byte {
	var mes Messages
	p := mes.Decode(msg)
	if p == nil {
		log.Error("mes.Decode p==nil")
		return nil
	}
	if msgI, ok := MessageMap[PakeId(p.GetSession().Id)]; !ok {
		log.Error("not find")
		return nil
	} else {
		r := reflect.New(msgI.Type())
		if msgI.Kind() == reflect.Ptr {
			r.Elem().Set(reflect.New(msgI.Elem().Type()))
		}
		m := r.Elem().Interface().(MessageI)
		m.Init(p.GetSession())
		//fmt.Println(m)
		err := json.Unmarshal(p.GetBody(), m.GetReq())
		if err != nil {
			log.Error(err)
			return nil
		}
		//fmt.Println(reflect.TypeOf(t))
		//fmt.Println("-nnnn", reflect.TypeOf(m))
		m.Process()

		buf, err := json.Marshal(m.GetRsp())
		if err != nil {
			log.Error(err)
			return nil
		}
		mes.Context = *(p.GetSession())
		return mes.Encode(buf)
	}
}
