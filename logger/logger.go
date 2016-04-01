package logger

import (
	"os"
	"runtime"
	"time"
	"fmt"
	"github.com/udoyu/utils/simini"
)

type Logger struct{
	buf chan string
	max_size int64
	size int64
	file *os.File
	logdir string
	filename string
	level LOG_LEVEL
}
func (this *Logger) Init(){
	fmt.Println("init")
	
	var ini simini.SimIni 
	err :=ini.LoadFile("conf.ini")
	if err!=0{
		fmt.Println("Logger Init err|",err)
		os.Exit(1)
	}
	v, e := ini.GetIntValWithDefault("log","max_size",50*1024*1024/*50M*/)
	if e != nil{
		fmt.Println("Logger Init GetIntValWithDefault e|",e)
		os.Exit(1)
	}
	this.max_size=int64(v)
	
	v, e = ini.GetIntValWithDefault("log","level",2)
	if e != nil{
		fmt.Println("Logger Init GetIntValWithDefault level e|",e)
		os.Exit(1)
	}
	this.level=LOG_LEVEL(v)

	this.logdir = ini.GetStringValWithDefault("log","logdir","log")
	
	this.filename=fmt.Sprintf("%s/log.log",this.logdir)
	this.openFile()
	this.buf = make(chan string,10240)
	fmt.Println(this)
	go this.logWrite()
}

func (this *Logger) checkLevel(level LOG_LEVEL) bool{
	if level <=this.level {
		return true
	}else { 
		return false
	}
}
func (this *Logger) LOG_DEVEL_FORMAT(level LOG_LEVEL)[]byte{
	now := time.Now()
	year,month,day :=now.Date()
	hour,min,sec :=now.Clock()
	
	_,file,line,_:= runtime.Caller(0)
	return []byte(fmt.Sprintf("%4d-%02d-%02d:%02d:%02d%02d[%s]%s:%d |",year,month,day,hour,min,sec,LOG_STR[level],file,line))
}
func (this *Logger) LOG_FATAL(format string, params ...interface{}){
	if this.checkLevel(FATAL)==false{
		return
	}
	timestr:=this.LOG_DEVEL_FORMAT(FATAL)
	context:=fmt.Sprintf(format,params...)
	var logstr []byte
	logstr=append(logstr,timestr...)
	logstr=append(logstr,[]byte(context)...)
	this.Output(logstr)
}

func (this *Logger) LOG_DEBUG(format string, params ...interface{}){
	if this.checkLevel(DEBUG)==false{
		return
	}
	timestr:=this.LOG_DEVEL_FORMAT(DEBUG)
	context:=fmt.Sprintf(format,params...)
	var logstr []byte
	logstr=append(logstr,timestr...)
	logstr=append(logstr,[]byte(context)...)
	this.Output(logstr)
}
func (this *Logger) LOG_ERROR(format string, params ...interface{}){
	if this.checkLevel(ERROR)==false{
		return
	}
	timestr:=this.LOG_DEVEL_FORMAT(ERROR)
	context:=fmt.Sprintf(format,params...)
	var logstr []byte
	logstr=append(logstr,timestr...)
	logstr=append(logstr,[]byte(context)...)
	this.Output(logstr)
}

func (this *Logger) LOG_WARNING(format string, params ...interface{}){
	if this.checkLevel(WARNING)==false{
		return
	}
	timestr:=this.LOG_DEVEL_FORMAT(WARNING)
	context:=fmt.Sprintf(format,params...)
	var logstr []byte
	logstr=append(logstr,timestr...)
	logstr=append(logstr,[]byte(context)...)
	this.Output(logstr)
}

func (this *Logger) LOG_INFO(format string, params ...interface{}){
	if this.checkLevel(INFO)==false{
		return
	}
	timestr:=this.LOG_DEVEL_FORMAT(INFO)
	context:=fmt.Sprintf(format,params...)
	var logstr []byte
	logstr=append(logstr,timestr...)
	logstr=append(logstr,[]byte(context)...)
	this.Output(logstr)
}

func (this *Logger) Output(logstr []byte){
	this.buf <- string(logstr)
}

func (this *Logger) openFile(){
	var err error
	this.file,err = os.OpenFile(this.filename, os.O_APPEND|os.O_CREATE, 0666)
	if err!=nil{
		fmt.Println("Logger openFile err|",err)
		os.Exit(1)
	}
	fi,err:= this.file.Stat()
	if err!=nil{
		fmt.Println("Logger openFile  Stat err|",err)
		os.Exit(1)
	}
	this.size = fi.Size()
}

func (this *Logger) rename(){
	//fmt.Println("rename-------------")
	now := time.Now()
	year,month,day :=now.Date()
	hour,min,sec :=now.Clock()
	filename:=fmt.Sprintf("%s/%4d-%02d-%02d-%02d-%02d-%d.log",this.logdir,year,month,day,hour,min,sec)
	
	this.file.Close()
	err := os.Rename(this.filename,filename)
	if err!=nil{
		fmt.Println("Logger rename err|",err)
		os.Exit(1)
	}
	//
}
func (this *Logger)changefile(){
	this.rename()
	this.openFile()
	this.size = 0
}

func (this *Logger) logWrite() {
	for {
		select {
		case str := <-this.buf:
			//fmt.Println(string(logstr))
			if (this.size+ int64(len(str))) >= this.max_size{
				fmt.Println("output",this.size+ int64(len(str)))
				this.changefile()
			}
			fmt.Fprintln(this.file, str)
			this.size += int64(len(str))
			//this.file.Sync()
		}
	}
	//fmt.Println(string(logstr))
}