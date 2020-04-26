package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Level string
}

func NewLogger(c Config) *logrus.Logger {
	logger := logrus.Logger{}

	switch c.Level {
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		fallthrough
	default:
		logger.SetLevel(logrus.DebugLevel)
	}

	formatter := logrus.JSONFormatter{}
	logger.SetFormatter(&formatter)
	logger.SetOutput(os.Stdout)

	return &logger
}
