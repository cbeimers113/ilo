package log

import (
	"fmt"
	"os"
	"time"
)

// ANSI escape codes for colors
const (
	colReset  = "\033[0m"
	colRed    = "\033[31m"
	colGreen  = "\033[32m"
	colYellow = "\033[33m"
	colBlue   = "\033[34m"
)

type logLevel uint8

const (
	logLevelDebug = iota
	logLevelInfo
	logLevelWarn
	logLevelError
)

func log(level logLevel, msg string) {
	prefix := colReset
	switch level {
	case logLevelDebug:
		prefix = colBlue
	case logLevelInfo:
		prefix = colGreen
	case logLevelWarn:
		prefix = colYellow
	case logLevelError:
		prefix = colRed
	}

	timestamp := time.Now().Format(time.Kitchen)
	fmt.Printf("[%s] %s%s%s\n", timestamp, prefix, msg, colReset)
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
