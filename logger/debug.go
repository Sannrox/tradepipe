package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Enable sets the DEBUG env var to true
// and makes the logger to log at debug level.
func Enable() {
	logrus.Info("Enable debug")
	os.Setenv("DEBUG", "1")
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

// Disable sets the DEBUG env var to false
// and makes the logger to log at info level.
func Disable() {
	os.Setenv("DEBUG", "")
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(false)
}

// IsEnabled checks whether the debug flag is set or not.
func IsEnabled() bool {
	return len(os.Getenv("DEBUG")) != 0
}
