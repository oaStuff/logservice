package asyncLogger


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
	logQueue 		*LinkedList
	log 			*logging.Logger
	logConfig 		LoggerConfig
	callStackSkip	int
	infoState		bool
	errorState		bool
	warnState		bool
	critState		bool
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

func (logger *Logger) EnableInfo(state bool) (prev bool) {
	prev, logger.infoState = logger.infoState, state
	return
}

func (logger *Logger) EnableError(state bool) (prev bool) {
	prev, logger.errorState = logger.errorState, state
	return
}

func (logger *Logger) EnableWar(state bool) (prev bool) {
	prev, logger.warnState = logger.warnState, state
	return
}

func (logger *Logger) EnableCrit(state bool) (prev bool) {
	prev, logger.critState = logger.critState, state
	return
}

func (logger *Logger) Info(msg string)  {

	if !logger.logConfig.Enabled || !logger.infoState {
		return
	}

	details := logDetails{level:iNFO}
	_, file, line, ok := runtime.Caller(logger.callStackSkip)
	if ok {
		_, filename := filepath.Split(file)
		details.message = fmt.Sprintf(" ▶  %s ( %s:%d )",msg, filename, line)
	}else {
		details.message = msg
	}

	logger.logQueue.Add(details)
}

func (logger *Logger) Warn(msg string)  {

	if !logger.logConfig.Enabled || !logger.warnState {
		return
	}

	details := logDetails{level:wARN}
	_, file, line, ok := runtime.Caller(logger.callStackSkip)
	if ok {
		_, filename := filepath.Split(file)
		details.message = fmt.Sprintf(" ▶  %s ( %s:%d )",msg, filename, line)
	}else {
		details.message = msg
	}

	logger.logQueue.Add(details)
}

func (logger *Logger) Error(msg string)  {

	if !logger.logConfig.Enabled || !logger.errorState {
		return
	}

	details := logDetails{level:eRROR}
	_, file, line, ok := runtime.Caller(logger.callStackSkip)
	if ok {
		_, filename := filepath.Split(file)
		details.message = fmt.Sprintf(" ▶  %s ( %s:%d )",msg, filename, line)
	}else {
		details.message = msg
	}

	logger.logQueue.Add(details)
}

func (logger *Logger) Critical(msg string)  {

	if !logger.logConfig.Enabled || !logger.critState {
		return
	}

	details := logDetails{level:cRITICAL}
	_, file, line, ok := runtime.Caller(logger.callStackSkip)
	if ok {
		_, filename := filepath.Split(file)
		details.message = fmt.Sprintf(" ▶  %s ( %s:%d )",msg, filename, line)
	}else {
		details.message = msg
	}

	logger.logQueue.Add(details)
}

func New(config LoggerConfig) *Logger {
	logger := &Logger{}
	logger.callStackSkip = 1
	logger.logConfig = config
	logger.logConfig.Enabled = logger.logConfig.Enabled && (logger.logConfig.AllowFileLog || logger.logConfig.AllowConsoleLog)
	logger.critState = true
	logger.errorState = true
	logger.infoState = true
	logger.warnState = true

	logger.initilize()

	return logger
}

func (logger *Logger) SetCallStackSkip(skip int)  {
	logger.callStackSkip = skip
}

func (logger *Logger) GetCallStackSkip() int {
	return logger.callStackSkip
}

func (logger *Logger) initilize() {

	if "" == logger.logConfig.Filename {
		b := make([]byte, 5)
		rand.Read(b)
		logger.logConfig.Filename = os.Args[0] + fmt.Sprintf("%x",b)
	}

	var ll *lumberjack.Logger

	if logger.logConfig.Filename == "" {
		ll = &lumberjack.Logger{
			Filename:   "logs/" + os.Args[0] + ".log",
			MaxAge:     30,
			MaxBackups: 20,
			MaxSize:    1,
		}
	} else {
		ll = &lumberjack.Logger{
			Filename:   logger.logConfig.Filename,
			MaxAge:     30,
			MaxBackups: 20,
			MaxSize:    1,
		}
	}

	var backends []logging.Backend

	logger.logQueue = NewLinkedList(true)
	logger.log = logging.MustGetLogger(logger.logConfig.Filename)

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
