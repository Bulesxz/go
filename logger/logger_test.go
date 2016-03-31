package logger
import (
	"testing"
)
func Test_LOG_FATAL(t *testing.T){
	var log *Logger=&Logger{}
	log.Init()
	for i:=0;i<1000000;i++{
		log.LOG_FATAL("%d",i)
		//log.LOG_DEBUG("ssss%d %d",1,2)
	}
}