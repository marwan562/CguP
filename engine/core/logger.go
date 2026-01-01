package core

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogger() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(format string, v ...interface{}) {
	if InfoLogger == nil {
		InitLogger()
	}
	InfoLogger.Printf(format, v...)
}

func LogError(format string, v ...interface{}) {
	if ErrorLogger == nil {
		InitLogger()
	}
	ErrorLogger.Printf(format, v...)
}

func LogFatal(format string, v ...interface{}) {
	if ErrorLogger == nil {
		InitLogger()
	}
	ErrorLogger.Fatalf(format, v...)
}
