package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type LogLevel string

const (
	DEBUG LogLevel = "debug"
	INFO  LogLevel = "info"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
	FATAL LogLevel = "fatal"
)

func NewLogger() *Logger {
	logger := logrus.New()

	// Set formatter
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	logger.SetOutput(os.Stdout)

	// Set level
	logger.SetLevel(logrus.InfoLevel)

	// Report caller
	logger.SetReportCaller(true)

	return &Logger{logger}
}

// NewFileLogger creates logger with file output
func NewFileLogger(filename string) *Logger {
	logger := NewLogger()

	// Create logs directory if not exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Open log file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger.SetOutput(file)
	return logger
}

func (l *Logger) SetLevel(level LogLevel) {
	switch level {
	case DEBUG:
		l.Logger.SetLevel(logrus.DebugLevel)
	case INFO:
		l.Logger.SetLevel(logrus.InfoLevel)
	case WARN:
		l.Logger.SetLevel(logrus.WarnLevel)
	case ERROR:
		l.Logger.SetLevel(logrus.ErrorLevel)
	case FATAL:
		l.Logger.SetLevel(logrus.FatalLevel)
	default:
		l.Logger.SetLevel(logrus.InfoLevel)
	}
}

// WithFields adds fields to log entry
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithError adds error to log entry
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithRequest logs HTTP request details
func (l *Logger) WithRequest(method, path, userAgent, ip string) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields{
		"method":     method,
		"path":       path,
		"user_agent": userAgent,
		"ip":         ip,
	})
}

// LogError logs error with context
func (l *Logger) LogError(err error, context string, fields ...map[string]interface{}) {
	entry := l.Logger.WithError(err).WithField("context", context)

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Error("Error occurred")
}

// LogTransaction logs transaction details
func (l *Logger) LogTransaction(transactionID, userID, action string, amount float64, status string) {
	l.Logger.WithFields(logrus.Fields{
		"transaction_id": transactionID,
		"user_id":        userID,
		"action":         action,
		"amount":         amount,
		"status":         status,
	}).Info("Transaction log")
}

// LogPayment logs payment details
func (l *Logger) LogPayment(paymentID, method, provider, status string, amount float64) {
	l.Logger.WithFields(logrus.Fields{
		"payment_id": paymentID,
		"method":     method,
		"provider":   provider,
		"status":     status,
		"amount":     amount,
	}).Info("Payment log")
}
