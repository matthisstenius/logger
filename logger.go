package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	debugLevel   = "DEBUG"
	infoLevel    = "INFO"
	warningLevel = "WARNING"
	errorLevel   = "ERROR"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(extractLogLevel())
}

// Fields custom data to be logged
type Fields map[string]interface{}

// Logger
type Logger struct {
	fields logrus.Fields
}

// New makes new Logger object
func New() *Logger {
	fields := logrus.Fields{}
	fptr, file, row, ok := runtime.Caller(1)
	if ok {
		fields["function"] = runtime.FuncForPC(fptr).Name()
		fields["file"] = file
		fields["row"] = row
	}
	return &Logger{fields: fields}
}

// WithFields makes new Logger with custom fields
func WithFields(f Fields) *Logger {
	logger := New()
	for k, v := range f {
		logger.fields[k] = v
	}
	return logger
}

// Error log on error level
func (l *Logger) Error(message string) {
	l.fields["severity"] = errorLevel
	l.appendFields().Error(fmt.Sprintf("[ERROR] %s", message))
}

// Info log in info level
func (l *Logger) Info(message string) {
	l.fields["severity"] = infoLevel
	l.appendFields().Info(fmt.Sprintf("[INFO] %s", message))
}

// Warning log on warning level
func (l *Logger) Warning(message string) {
	l.fields["severity"] = warningLevel
	l.appendFields().Warning(fmt.Sprintf("[WARNING] %s", message))
}

// Debug log on debug level
func (l *Logger) Debug(message string) {
	l.fields["severity"] = debugLevel
	l.appendFields().Warning(fmt.Sprintf("[DEBUG] %s", message))
}

// Panic log on panic level
func (l *Logger) Panic(message string) {
	l.appendFields().Panic(fmt.Sprintf("[PANIC] %s", message))
}

func (l *Logger) appendFields() *logrus.Entry {
	return logrus.WithFields(l.fields)
}

func extractLogLevel() logrus.Level {
	switch os.Getenv("LOG_LEVEL") {
	case debugLevel:
		return logrus.DebugLevel
	case infoLevel:
		return logrus.InfoLevel
	case warningLevel:
		return logrus.WarnLevel
	case errorLevel:
		return logrus.ErrorLevel
	default:
		return logrus.ErrorLevel
	}
}
