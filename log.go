package logger

import (
    "github.com/sirupsen/logrus"
    "os"
    "github.com/Shopify/logrus-bugsnag"
    "github.com/bugsnag/bugsnag-go"
)

func init() {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetLevel(extractLogLevel())

    bugsnag.Configure(bugsnag.Configuration{
        APIKey: os.Getenv("BUGSNAG_API_KEY"),
    })
    hook, err := logrus_bugsnag.NewBugsnagHook()
    if err != nil {
        panic(err)
    }

    logrus.AddHook(hook)
}

type Fields map[string]interface{}

type Logger struct {
    fields logrus.Fields
}

// WithFields map Fields into logrus Fields
func (l *Logger) WithFields(f Fields) *Logger {
    fields := make(logrus.Fields, len(f))
    for k, v := range f {
        fields[k] = v
    }
    l.fields = fields
    return l
}

func (l *Logger) Error(message string) {
    l.appendFields()
    logrus.Error(message)
}

func (l *Logger) Info(message string) {
    l.appendFields()
    logrus.Info(message)
}

func (l *Logger) Warning(message string) {
    l.appendFields()
    logrus.Warning(message)
}

func (l *Logger) appendFields() {
    if len(l.fields) > 0 {
        logrus.WithFields(l.fields)
    }
}

func extractLogLevel() logrus.Level {
    var level logrus.Level

    switch os.Getenv("LOG_LEVEL") {
    case "DEBUG":
        level = logrus.DebugLevel
        break
    case "WARNING":
        level = logrus.WarnLevel
        break
    case "ERROR":
        level = logrus.ErrorLevel
        break
    default:
        level = logrus.ErrorLevel
    }

    return level
}

logger.WithFields