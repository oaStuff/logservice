package logger


import (
"github.com/op/go-logging"
"os"
"runtime"
"fmt"
"path/filepath"
"gopkg.in/natefinch/lumberjack.v2"
)

const (
	iNFO = iota
	wARN
	eRROR
	cRITICAL
)

var logQueue *LinkedList
var log *logging.Logger
var loggerEnabled = false

type logDetails struct {
	message string;
	level int
}

func Info(msg string)  {

	if !loggerEnabled {
		return
	}

	details := logDetails{level:iNFO}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, filename := filepath.Split(file)
		details.message = fmt.Sprintf(" ▶  %s ( %s:%d )",msg, filename, line)
	}else {
		details.message = msg
	}

	logQueue.Add(details)
}

func Warn(msg string)  {

	if !loggerEnabled {
		return
	}

	details := logDetails{level:wARN}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, filename := filepath.Split(file)
		details.message = fmt.Sprintf(" ▶  %s ( %s:%d )",msg, filename, line)
	}else {
		details.message = msg
	}

	logQueue.Add(details)
}

func Error(msg string)  {

	if !loggerEnabled {
		return
	}

	details := logDetails{level:eRROR}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, filename := filepath.Split(file)
		details.message = fmt.Sprintf(" ▶  %s ( %s:%d )",msg, filename, line)
	}else {
		details.message = msg
	}

	logQueue.Add(details)
}

func Critical(msg string)  {

	if !loggerEnabled {
		return
	}

	details := logDetails{level:cRITICAL}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, filename := filepath.Split(file)
		details.message = fmt.Sprintf(" ▶  %s ( %s:%d )",msg, filename, line)
	}else {
		details.message = msg
	}

	logQueue.Add(details)
}

func EnableLogging(enabled bool)  {
	loggerEnabled = enabled
}

func init() {

	ll := &lumberjack.Logger{
		Filename:"logs/" + os.Args[0] + ".log",
		MaxAge:30,
		MaxBackups:20,
		MaxSize:1,
	}

	logQueue = NewLinkedList(true)
	log = logging.MustGetLogger("openFEP")
	consoleFormat := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{level:.4s}] %{message} %{color:reset}`)
	fileFormat := logging.MustStringFormatter(`%{time:15:04:05.000} [%{level:.4s}] %{message} `)
	consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	fileBackend := logging.NewLogBackend(ll,"",0)
	consoleBackendFormatter := logging.NewBackendFormatter(consoleBackend, consoleFormat)
	fileBackendFormatter := logging.NewBackendFormatter(fileBackend,fileFormat)
	logging.SetBackend(consoleBackendFormatter,fileBackendFormatter)
	go doLogging()
}

func doLogging()  {
	for{
		details, ok := logQueue.Take().(logDetails)
		if ok {
			switch details.level {
			case iNFO:
				log.Info(details.message)
			case wARN:
				log.Warning(details.message)
			case eRROR:
				log.Error(details.message)
			case cRITICAL:
				log.Critical(details.message)
			}
		}
	}
}
