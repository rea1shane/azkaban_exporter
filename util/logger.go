package util

import (
	formatter "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	logger     *log.Logger
	onceLogger sync.Once
)

func GetLogger() *log.Logger {
	onceLogger.Do(func() {
		logger = log.New()
		logger.SetFormatter(&formatter.Formatter{
			TimestampFormat: "2006-01-02 | 15:04:05",
			FieldsOrder: []string{"name", "duration_seconds"},
			HideKeys: true,
		})
	})
	return logger
}
