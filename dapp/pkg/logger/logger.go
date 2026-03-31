package logger

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	warnLog  *log.Logger
	debugLog *log.Logger
}

var Log *Logger

func Init() {
	Log = &Logger{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		warnLog:  log.New(os.Stdout, "WARN\t", log.Ldate|log.Ltime|log.Lshortfile),
		debugLog: log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func Info(format string, v ...interface{}) {
	Log.infoLog.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	Log.errorLog.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
	Log.warnLog.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	Log.debugLog.Printf(format, v...)
}

func GetCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return filepath.Base(file) + ":" + string(rune(line))
}
