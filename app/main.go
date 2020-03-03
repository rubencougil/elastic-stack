package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rubencougil/geekshubs/elastic/app/user"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgin"
)

func main() {

	logger := Logger()
	db := Database(logger)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(apmgin.Middleware(r))

	r.Use(ginlogrus.Logger(logger), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := apm.StartSpan(ctx, "getRoot", "")
		defer span.End()
		c.JSON(200, gin.H{})
	})
	r.POST("/test", func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := apm.StartSpan(ctx, "getExpensiveCalc", "custom")
		defer span.End()

		resp, err := http.Get("http://host.docker.internal:3200/expensive_calc")
		if err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("Error with expensive calc")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		logger.WithFields(logrus.Fields{"body": body}).Info("Expensive calc")

		c.JSON(200, gin.H{"message": "ok"})

	})
	r.GET("/error", func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /foo", "request")
		defer tx.End()
		err := errors.New("emit macho dwarf: elf header corrupted")
		e := apm.DefaultTracer.NewError(err)
		e.SetTransaction(tx)
		e.Send()
		c.JSON(200, gin.H{})
	})

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
	logger.SetLevel(logrus.DebugLevel)
	log.SetOutput(logger.Writer())
	logger.SetOutput(io.MultiWriter(os.Stdout))
	return logger
}
