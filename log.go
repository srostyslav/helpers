package helpers

import (
	"log"
	"os"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func InitLogger() error {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	output := &lumberjack.Logger{
		Filename:   "./logs.txt",
		MaxSize:    10, // megabytes after which new file is created
		MaxBackups: 10, // number of backups
		MaxAge:     28, //days
	}
	InfoLogger.SetOutput(output)
	WarningLogger.SetOutput(output)
	ErrorLogger.SetOutput(output)

	return nil
}
