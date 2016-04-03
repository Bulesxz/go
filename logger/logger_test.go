package logger
import (
//	"fmt"
//	"time"
	"testing"
)
/*
func Test_LOG_FATAL(t *testing.T){
	var log *Logger=&Logger{}
	log.init()
	for i:=0;i<1;i++{
		log.LOG_FATAL("%d",i)
		time.Sleep(10*time.Microsecond)
		log.LOG_DEBUG("ssss%d %d",1,2)
		time.Sleep(10*time.Microsecond)
		log.LOG_INFO("ssss%d %d",1,2)
		time.Sleep(10*time.Microsecond)
		log.LOG_ERROR("ssss%d %d",1,2)
		time.Sleep(10*time.Microsecond)
		log.LOG_WARNING("ssss%d %d",1,2)
	}
}
*/
func Test_LOG_DEBUG(t *testing.T){
	//fmt.Println("---------------****")
	Debug("ssss|",1,2)
	Error("ssss|",1,2)
}