package logger

// fix this
import (
	"log"
	"os"
	"path/filepath"

	oz "github.com/go-ozzo/ozzo-log"
)

type Logger struct {
	ozzo       *oz.Logger
	logDirPath string
}

func NewLogger() *Logger {
	var logger Logger
	logger.ozzo = oz.NewLogger()
	logger.logDirPath = "../logs"

	err := os.Mkdir("../logs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("cannot create logger path: %s", err)
	}

	logFile, err := os.OpenFile(filepath.Join(logger.logDirPath, "server.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("cannot open logger file: %s", err)
	}

	targetConsole := oz.NewConsoleTarget()
	targetFile := oz.NewFileTarget()
	targetFile.FileName = logFile.Name()

	logger.ozzo.Targets = append(logger.ozzo.Targets, targetFile, targetConsole)

	logger.ozzo.Open()

	log.Print("Logger is created")

	return &logger
}

func (l *Logger) Error(errorMsg string) {
	l.ozzo.Error(errorMsg)
}

func (l *Logger) Debug(debugMsg string) {
	l.ozzo.Debug(debugMsg)
}
