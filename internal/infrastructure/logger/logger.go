package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	// INFO represents the info log level
	INFO = "info"
	// ERROR represents the error log level
	ERROR = "error"
	// DEBUG represents the debug log level
	DEBUG = "debug"
	// WARN represents the warn log level
	WARN = "warn"

	// DefaultLoggerLevel represents the default logger level
	DefaultLoggerLevel = logrus.InfoLevel
)

// LevelMap maps the log level to logrus level
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

// WithLogLevel sets the log level
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
