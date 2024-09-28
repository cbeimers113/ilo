package log

import (
	"fmt"
	"os"
	"time"

	"cbeimers113/ilo/internal/constant"
)

type logLevel uint8

const (
	logLevelDebug = iota
	logLevelInfo
	logLevelWarn
	logLevelError
)

func log(level logLevel, msg string) string {
	prefix := constant.ColReset
	switch level {
	case logLevelDebug:
		prefix = constant.ColBlue
	case logLevelInfo:
		prefix = constant.ColGreen
	case logLevelWarn:
		prefix = constant.ColYellow
	case logLevelError:
		prefix = constant.ColRed
	}

	timestamp := time.Now().Format(time.Kitchen)
	str := fmt.Sprintf("[%s] %s%s%s", timestamp, prefix, msg, constant.ColReset)
	fmt.Println(str)
	return str
}

// Debug prints a message at the debug level
func Debug(msg string) {
	log(logLevelDebug, msg)
}

// Info prints a message at the info level
func Info(msg string) {
	log(logLevelInfo, msg)
}

// Warn prints a message at the warn level
func Warn(msg string) {
	log(logLevelWarn, msg)
}

// Error prints a message at the error level
func Error(msg string) {
	log(logLevelError, msg)
}

// Fatal prints a message at the error level and exits the program
func Fatal(msg string) {
	log(logLevelError, msg)
	os.Exit(1)
}
