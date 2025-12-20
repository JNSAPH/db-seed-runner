package internal

import (
	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	// Setup the logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	})
}
