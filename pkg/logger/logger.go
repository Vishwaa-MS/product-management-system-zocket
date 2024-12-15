package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	
	// Set log format
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	
	// Set log level
	Log.SetLevel(logrus.InfoLevel)
	
	// Set output to stdout
	Log.SetOutput(os.Stdout)
}

func LogAPIRequest(method string, path string, status int, responseTime time.Duration) {
	Log.WithFields(logrus.Fields{
		"method":        method,
		"path":          path,
		"status_code":   status,
		"response_time": responseTime,
	}).Info("API Request")
}

func LogImageProcessingEvent(imageURL string, success bool) {
	Log.WithFields(logrus.Fields{
		"image_url": imageURL,
		"success":   success,
	}).Info("Image Processing Event")
}