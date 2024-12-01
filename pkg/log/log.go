package log

import (
	"fmt"
	"os"
	"time"

	"k8s.io/klog"
)

type LogLevel int

const (
	Info LogLevel = iota
	Warn
	Error
	File
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("unique_requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		klog.Fatalf("Error opening log file: %v", err)
	}
}

func Print(level LogLevel, format string, args ...interface{}) {
	logMsg := fmt.Sprintf(format, args...)

	switch level {
	case Warn:
		klog.Warning(logMsg)
	case Error:
		klog.Error(logMsg)
	case Info:
		fallthrough
	default:
		klog.Info(logMsg)
		if level == File {
			_, err := logFile.WriteString(fmt.Sprintf("%s %s\n", time.Now().Format("2006-01-02 15:04:05"), logMsg))
			if err != nil {
				klog.Errorf("Error writing to log file: %v", err)
			}
		}
	}
}
