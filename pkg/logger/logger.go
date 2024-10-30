package logger

import (
	"fmt"
	"log"
	"os"
)

type CustomLogger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
	Fatal       *log.Logger
}

func NewLogger() *CustomLogger {
	return &CustomLogger{}
}

func (LoggerObject *CustomLogger) GetLoggerObject(infoFilePath, errorFilePath, debugFilePath, part string) *CustomLogger {

	var option int
	if part == "HTTP" {
		option = os.O_CREATE | os.O_RDWR | os.O_APPEND | os.O_TRUNC
	} else {
		option = os.O_CREATE | os.O_RDWR | os.O_APPEND
	}
	file, err := os.OpenFile(infoFilePath, option, 0666)
	if err != nil {
		log.Fatalln("Error opening info log file: ", err)
	}
	LoggerObject.InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)

	file, err = os.OpenFile(errorFilePath, option, 0666)
	if err != nil {
		log.Fatalln("Error opening error log file: ", err)
	}
	LoggerObject.ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)

	file, err = os.OpenFile(debugFilePath, option, 0666)
	if err != nil {
		log.Fatalln("Error opening error log file: ", err)
	}
	LoggerObject.DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile)

	return LoggerObject

}

func ErrorWrapper(layer, functionName, context string, err error) error {
	return fmt.Errorf("%s %w\n", fmt.Sprintf("[Layer:%s,Function: %s,Context: %s]--->", layer, functionName, context), err)
}
