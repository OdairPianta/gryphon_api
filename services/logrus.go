package services

import (
	"os"

	"github.com/sirupsen/logrus"
)

var LOGGER = logrus.New()

func InitLogger() {
	file, err := os.OpenFile("temp/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		LOGGER.Out = file
	} else {
		LOGGER.Info("Failed to log to file, using default stderr")
	}
	LOGGER.SetFormatter(&logrus.TextFormatter{})
}

func WriteLog(name string, fields map[string]interface{}) {
	LOGGER.WithFields(fields).Info(name)
}
