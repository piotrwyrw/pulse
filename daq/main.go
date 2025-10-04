package main

import (
	"daq/internal/bootstrap"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.DebugLevel)

	bootstrap.Bootstrap()

	logrus.Info("Goodbye!")
}
