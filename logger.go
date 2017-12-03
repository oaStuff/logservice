package logger


import (
	"github.com/op/go-logging"
	"os"
	"runtime"
	"fmt"
	"path/filepath"
	"gopkg.in/natefinch/lumberjack.v2"
	"math/rand"
)

const (
	iNFO = iota
	wARN
	eRROR
	cRITICAL
)

type Logger struct {
	logQueue 	*LinkedList
	log 		*logging.Logger
	logConfig 	LoggerConfig
}

type LoggerConfig struct {
	Enabled			bool
	AllowFileLog	bool
	AllowConsoleLog	bool
	Filename		string
}


type logDetails struct {
	message string
	level int
}

func (logger *Logger) Info(msg string)  {

	if !logger.logConfig.Enabled {
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

	logger.logQueue.Add(details)
}

func (logger *Logger) Warn(msg string)  {

	if !logger.logConfig.Enabled {
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

	logger.logQueue.Add(details)
}

func (logger *Logger) Error(msg string)  {

	if !logger.logConfig.Enabled {
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

	logger.logQueue.Add(details)
}

func (logger *Logger) Critical(msg string)  {

	if !logger.logConfig.Enabled {
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

	logger.logQueue.Add(details)
}

func Neww(config LoggerConfig) *Logger {
	logger := &Logger{}
	logger.logConfig = config
	logger.logConfig.Enabled = logger.logConfig.Enabled && (logger.logConfig.AllowFileLog || logger.logConfig.AllowConsoleLog)
	logger.initilize()

	return logger
}

func (logger *Logger) initilize() {

	if "" == logger.logConfig.Filename {
		b := make([]byte, 5)
		rand.Read(b)
		logger.logConfig.Filename = os.Args[0] + fmt.Sprintf("%x",b)
	}

	ll := &lumberjack.Logger{
		Filename:"logs/" + logger.logConfig.Filename + ".log",
		MaxAge:30,
		MaxBackups:20,
		MaxSize:1,
	}

	var backends []logging.Backend

	logger.logQueue = NewLinkedList(true)
	logger.log = logging.MustGetLogger("openFEP")

	if logger.logConfig.AllowConsoleLog {
		consoleFormat := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{level:.4s}] %{message} %{color:reset}`)
		consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
		consoleBackendFormatter := logging.NewBackendFormatter(consoleBackend, consoleFormat)
		backends = append(backends, consoleBackendFormatter)
	}


	if logger.logConfig.AllowFileLog {
		fileFormat := logging.MustStringFormatter(`%{time:15:04:05.000} [%{level:.4s}] %{message} `)
		fileBackend := logging.NewLogBackend(ll, "", 0)
		fileBackendFormatter := logging.NewBackendFormatter(fileBackend, fileFormat)
		backends = append(backends, fileBackendFormatter)
	}

	if len(backends) > 0 {
		logging.SetBackend(backends...)
	}

	go logger.doLogging()
}

func (logger *Logger) doLogging()  {
	for{
		details, ok := logger.logQueue.Take().(logDetails)
		if ok {
			switch details.level {
			case iNFO:
				logger.log.Info(details.message)
			case wARN:
				logger.log.Warning(details.message)
			case eRROR:
				logger.log.Error(details.message)
			case cRITICAL:
				logger.log.Critical(details.message)
			}
		}
	}
}
