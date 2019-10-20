package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}
	filePath, _ := filepath.Abs("../var/log/app.log")
	logFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.WithError(err).Fatal("Cannot open log file")
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	return log
}
