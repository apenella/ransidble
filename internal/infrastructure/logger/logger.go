package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	INFO  = "info"
	ERROR = "error"
	DEBUG = "debug"
	WARN  = "warn"

	// DefaultLoggerLevel represents the default logger level
	DefaultLoggerLevel = logrus.InfoLevel
)

var LevelMap = map[string]logrus.Level{
	INFO:  logrus.InfoLevel,
	ERROR: logrus.ErrorLevel,
	DEBUG: logrus.DebugLevel,
	WARN:  logrus.WarnLevel,
}

// Logger represents a logger
type Logger struct {
	logger *logrus.Logger
}

// NewLogger creates a new logger
func NewLogger() *Logger {

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) WithLogLevel(level string) *Logger {
	level = strings.ToLower(level)
	_, ok := LevelMap[level]
	if !ok {
		level = INFO
	}

	l.logger.SetLevel(LevelMap[level])

	return l
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logger.WithFields(logrus.Fields{"data": fields}).Info(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.logger.WithFields(logrus.Fields{"data": fields}).Error(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.logger.WithFields(logrus.Fields{"data": fields}).Debug(msg)
}

// Warn logs a warn message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.logger.WithFields(logrus.Fields{"data": fields}).Warn(msg)
}
