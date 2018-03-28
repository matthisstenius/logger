package logger

import (
    "github.com/sirupsen/logrus"
    "os"
    "github.com/Shopify/logrus-bugsnag"
    "github.com/bugsnag/bugsnag-go"
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

type Fields map[string]interface{}

type Logger struct {
    fields logrus.Fields
}

// New init new Logger with base fields
func WithFields(f Fields) *Logger {
    fields := make(logrus.Fields, len(f))
    for k, v := range f {
        fields[k] = v
    }
    return &Logger{fields: fields}
}

func (l *Logger) Error(message string) {
    l.appendFields().Error(message)
}

func (l *Logger) Info(message string) {
    l.appendFields().Info(message)
}

func (l *Logger) Warning(message string) {
    l.appendFields().Warning(message)
}

func (l *Logger) Panic(message string) {
    l.appendFields().Panic(message)
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
