package logger

type LOG_LEVEL int

var LOG_STR = []string{
	"LOG_FATAL","LOG_ERROR","LOG_WARNING","LOG_INFO","LOG_DEBUG",
}
const (
    FATAL LOG_LEVEL= iota //0
    ERROR //1
    WARNING//2
    INFO//3
    DEBUG//4
)

var (
	glog *Logger = &Logger{}
)
func init() {
	glog.Init()
}


func Fatal(params ...interface{}) {
	glog.LOG_FATAL(params...)
}
func Error(params ...interface{}) {
	glog.LOG_ERROR(params...)
}
func Warning(params ...interface{}) {
	glog.LOG_WARNING(params...)
}
func Info(params ...interface{}) {
	glog.LOG_INFO(params...)
}
func Debug(params ...interface{}) {
	glog.LOG_DEBUG(params...)
}

/*
type LoggerInterface interface{
	LOG_FATAL(format string, params ...interface{})
	LOG_DEBUG(format string, params ...interface{})
}*/
