package Laas

import (
	"net/http"
	"log"
	"fmt"
	//"appengine"
)

const (
	LogLevelDebug = 1
	LogLevelInfo = 2
	LogLevelWarning = 3
	LogLevelError = 4
	LogLevelCritical = 5
)

type Logger interface{
	LogMessage(msg string, level int, args ...interface{})
}

type FileLogger struct {
	enabled bool
}

func (logger *FileLogger) GetLogFilePath() string {
	return ""
}

func (logger *FileLogger) WriteLog(format string, args ...interface{}) {
	if logger.enabled {
		log.Printf(format, args...) 
	}
}

type GAELogger struct {
	//Context appengine.Context
}

func (logger *GAELogger) Init(r *http.Request) {
	//logger.Context = appengine.NewContext(r)
}

func (logger *GAELogger) LogMessage(msg string, level int, args ...interface{}) {
	fmt.Printf(msg, args)
	/*switch level {
		case LogLevelDebug:
			logger.Context.Debugf(msg, args)
		case LogLevelInfo:
            logger.Context.Infof(msg, args)
		case LogLevelWarning:
            logger.Context.Warningf(msg, args)
		case LogLevelError:
            logger.Context.Errorf(msg, args)
		case LogLevelCritical:
            logger.Context.Criticalf(msg, args)
		default:
            logger.Context.Infof(msg, args)
	}*/
}
