package logger

type LOG_LEVEL int

var LOG_STR = []string{
	"LOG_FATAL","LOG_ERROR","WARNING","LOG_INFO","LOG_DEBUG",
}
const (
    FATAL LOG_LEVEL= iota
    ERROR
    WARNING
    INFO
    DEBUG
)

type LoggerInterface interface{
	LOG_FATAL(format string, params ...interface{})
	LOG_DEBUG(format string, params ...interface{})
}
