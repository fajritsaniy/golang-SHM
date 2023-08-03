package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/config"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/sirupsen/logrus"
)

func LogRequestMiddleware(log *logrus.Logger) gin.HandlerFunc {
	cfg, err := config.NewConfig()
	file, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Set up logrus
	log.SetOutput(file)

	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Since(startTime)
		entryLog := model.RequestLog{
			EndTime:      endTime,
			StatusCode:   c.Writer.Status(),
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			RelativePath: c.Request.URL.Path,
			UserAgent:    c.Request.UserAgent(),
		}

		switch {
		case c.Writer.Status() >= 500:
			log.Error(entryLog)
		case c.Writer.Status() >= 400:
			log.Warn(entryLog)
		default:
			log.Info(entryLog)
		}
	}
}
