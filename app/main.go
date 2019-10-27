package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

func main() {
	logger := Logger()

	for {
		logger.Info("Hello World")
		logger.WithError(errors.New("error in main")).Error("Error in app")
		time.Sleep(10 * time.Second)
	}
}

func Logger() *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}
	log.SetOutput(io.MultiWriter(os.Stdout))
	return log
}

