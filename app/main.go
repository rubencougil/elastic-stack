package main

import (
	"github.com/pkg/errors"
	"github.com/rubencougil/elastic/app/logger"
)

func main() {
	log := logger.NewLogger()

	log.Info("Hello World")
	log.WithError(errors.New("error in main")).Error("Error in app")
}

