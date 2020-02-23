package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rubencougil/geekshubs/elastic/app/user"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"io"
	"log"
	"os"
)

func main() {

	logger := Logger()
	db := Database(logger)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(ginlogrus.Logger(logger), gin.Recovery())

	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	r.POST("/create", user.CreateUserHandler(logger, user.NewUserStore(db, logger)))

	_ = r.Run(":80")
}

func Database(logger *logrus.Logger) *sqlx.DB {
	db, err := sqlx.Connect("mysql", "user:password@(db:3306)/db")
	if err != nil {
		logger.Fatalf("Cannot connect to the database: %v", err)
	}
	return db
}

func Logger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)
	log.SetOutput(logger.Writer())
	logger.SetOutput(io.MultiWriter(os.Stdout))
	return logger
}
