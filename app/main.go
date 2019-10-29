package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"io"
	"os"
)

func main() {

	logger := Logger()
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.Use(ginlogrus.Logger(logger), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	_ = r.Run(":80")
}

func Logger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(io.MultiWriter(os.Stdout))
	return logger
}

