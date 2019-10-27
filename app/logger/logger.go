package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	path := "./var/logs/"
	log.Formatter = &logrus.JSONFormatter{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	filePath, _ := filepath.Abs(path + "app.log")
	logFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.WithError(err).Fatal("Cannot open log file")
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	return log
}
