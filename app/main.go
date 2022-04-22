package main

import (
	"bytes"
	"fmt"
	"go.elastic.co/ecslogrus"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	elastic_logrus "github.com/interactive-solutions/go-logrus-elasticsearch"
	"github.com/jmoiron/sqlx"
	"github.com/olivere/elastic/v7"
	"github.com/rubencougil/geekshubs/elastic/app/user"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"go.elastic.co/apm/module/apmgin"
)

func main() {

	logger := Logger()

	db := Database(logger)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(ginlogrus.Logger(logger), gin.Recovery())
	r.Use(apmgin.Middleware(r))

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
	logger.SetFormatter(&ecslogrus.Formatter{
		DataKey: "labels",
	})
	logger.ReportCaller = true

	switch os.Getenv("LOG_TO") {
	case "elasticsearch":
		logger.SetOutput(&bytes.Buffer{})
		logger.Info("Logging to Elasticsearch")
		client, err := elastic.NewClient(
			elastic.SetURL("http://elasticsearch:9200"),
			elastic.SetBasicAuth("elastic", "elastic"),
			elastic.SetSniff(false))
		if err != nil {
			logger.WithError(err).Fatal("Failed to construct elasticsearch client")
		}
		hook, err := elastic_logrus.NewElasticHook(client, "some-host", logrus.InfoLevel, func() string {
			return fmt.Sprintf("%s-%s", "some-index", time.Now().Format("2006-01-02"))
		}, time.Second*5)

		if err != nil {
			logger.WithError(err).Fatal("Failed to create elasticsearch hook for logger")
		}

		logger.Hooks.Add(hook)
	case "stdout":
		logger.Info("Logging to Stdout")
		logger.SetLevel(logrus.DebugLevel)
		logger.SetOutput(io.MultiWriter(os.Stdout))
	}
	return logger
}
