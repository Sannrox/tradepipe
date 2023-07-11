package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	SetLogLevel("info")

	if os.Getenv("TRACE") == "1" {
		Enable()
	}
}

// SetLogLevel sets the logrus logging level
func SetLogLevel(logLevel string) {
	if len(logLevel) != 0 {
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse logging level: %s\n", logLevel)
			os.Exit(1)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func SetLogFile(filname string) error {
	file, err := os.OpenFile(filname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	return nil
}

func ErrorWrapper(err error, msg string) error {
	logrus.Errorf(msg+": %s", err.Error())
	return err
}
