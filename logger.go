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
var logConfig LoggerConfig

type LoggerConfig struct {
	Enabled			bool
	AllowFileLog	bool
	AllowConsoleLog	bool
}


type logDetails struct {
	message string;
	level int
}

func Info(msg string)  {

	if !logConfig.Enabled {
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

	if !logConfig.Enabled {
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

	if !logConfig.Enabled {
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

	if !logConfig.Enabled {
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

func ConfigLogger(config LoggerConfig)  {
	logConfig = config
	logConfig.Enabled = logConfig.Enabled && (logConfig.AllowFileLog || logConfig.AllowConsoleLog)
	initilize()
}

func initilize() {

	ll := &lumberjack.Logger{
		Filename:"logs/" + os.Args[0] + ".log",
		MaxAge:30,
		MaxBackups:20,
		MaxSize:1,
	}

	//var index = 0
	var backends []logging.Backend

	logQueue = NewLinkedList(true)
	log = logging.MustGetLogger("openFEP")

	if logConfig.AllowConsoleLog {
		consoleFormat := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{level:.4s}] %{message} %{color:reset}`)
		consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
		consoleBackendFormatter := logging.NewBackendFormatter(consoleBackend, consoleFormat)
		backends = append(backends, consoleBackendFormatter)
	}


	if logConfig.AllowFileLog {
		fileFormat := logging.MustStringFormatter(`%{time:15:04:05.000} [%{level:.4s}] %{message} `)
		fileBackend := logging.NewLogBackend(ll, "", 0)
		fileBackendFormatter := logging.NewBackendFormatter(fileBackend, fileFormat)
		backends = append(backends, fileBackendFormatter)
	}

	if len(backends) > 0 {
		logging.SetBackend(backends...)
	}

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
