package main

import (
	"github.com/pkg/errors"
	"github.com/rubencougil/elastic/app/logger"
	"time"
)

func main() {
	log := logger.NewLogger()

	for {
		log.Info("Hello World")
		log.WithError(errors.New("error in main")).Error("Error in app")
		time.Sleep(10 * time.Second)
	}
}

