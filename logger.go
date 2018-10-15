package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Shopify/logrus-bugsnag"
	"github.com/bugsnag/bugsnag-go"
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

	if os.Getenv("BUGSNAG_API_KEY") != "" {
		addBugsnagHook()
	}
}

// Fields custom data to be logged
type Fields map[string]interface{}

// Logger
type Logger struct {
	fields logrus.Fields
}

// New init new Logger with base fields
func WithFields(f Fields) *Logger {
	fptr, file, row, ok := runtime.Caller(1)
	if ok {
		f["function"] = runtime.FuncForPC(fptr).Name()
		f["file"] = file
		f["row"] = row
	}

	fields := make(logrus.Fields, len(f))
	for k, v := range f {
		fields[k] = v
	}
	return &Logger{fields: fields}
}

// Error log on error level
func (l *Logger) Error(message string) {
	l.appendFields().Error(fmt.Sprintf("[ERROR] %s", message))
}

// Info log in info level
func (l *Logger) Info(message string) {
	l.appendFields().Info(fmt.Sprintf("[INFO] %s", message))
}

// Warning log on warning level
func (l *Logger) Warning(message string) {
	l.appendFields().Warning(fmt.Sprintf("[WARNING] %s", message))
}

// Panic log on panic level
func (l *Logger) Panic(message string) {
	l.appendFields().Panic(fmt.Sprintf("[PANIC] %s", message))
}

func (l *Logger) appendFields() *logrus.Entry {
	return logrus.WithFields(l.fields)
}

func extractLogLevel() logrus.Level {
	var level logrus.Level

	switch os.Getenv("LOG_LEVEL") {
	case debugLevel:
		level = logrus.DebugLevel
		break
	case infoLevel:
		level = logrus.InfoLevel
		break
	case warningLevel:
		level = logrus.WarnLevel
		break
	case errorLevel:
		level = logrus.ErrorLevel
		break
	default:
		level = logrus.ErrorLevel
	}
	return level
}

func addBugsnagHook() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: os.Getenv("BUGSNAG_API_KEY"),
	})

	hook, err := logrus_bugsnag.NewBugsnagHook()
	if err != nil {
		panic(err)
	}
	logrus.AddHook(hook)
}
