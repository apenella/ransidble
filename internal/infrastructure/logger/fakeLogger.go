package logger

// FakeLogger represents a fake logger
type FakeLogger struct{}

// NewFakeLogger creates a new fake logger
func NewFakeLogger() *FakeLogger {
	return &FakeLogger{}
}

// Info logs an info message
func (l *FakeLogger) Info(msg string, fields ...interface{}) {}

// Error logs an error message
func (l *FakeLogger) Error(msg string, fields ...interface{}) {}

// Debug logs a debug message
func (l *FakeLogger) Debug(msg string, fields ...interface{}) {}

// Warn logs a warn message
func (l *FakeLogger) Warn(msg string, fields ...interface{}) {}
