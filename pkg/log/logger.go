package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// DefaultLogger returns the configured default logger.
func DefaultLogger() *logrus.Logger {
	return logger
}

// Fields wraps logrus.Fields (a map[string]interface{}).
type Fields logrus.Fields

// SetLogLevel sets the log level for the default logger.
func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

// SetLogFormatter sets the formatter for the default logger.
func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

// Debug logs a message at level Debug.
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args...)
	}
}

// DebugWithFields logs a message with fields at level Debug.
func DebugWithFields(msg interface{}, f Fields) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Debug(msg)
	}
}

// Println logs a message at level Info (alias for Info).
func Println(args ...interface{}) {
	Info("%v", args...)
}

// Info logs a formatted message at level Info.
func Info(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(fmt.Sprintf(format, args...))
	}
}

// InfoWithFields logs a message with fields at level Info.
func InfoWithFields(msg interface{}, f Fields) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Info(msg)
	}
}

// Warn logs a message at level Warn.
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}

// WarnWithFields logs a message with fields at level Warn.
func WarnWithFields(msg interface{}, f Fields) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Warn(msg)
	}
}

// Error logs a message at level Error.
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

// ErrorWithFields logs a message with fields at level Error.
func ErrorWithFields(msg interface{}, f Fields) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Error(msg)
	}
}

// Fatal logs a message at level Fatal and then exits.
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

// FatalWithFields logs a message with fields at level Fatal and then exits.
func FatalWithFields(msg interface{}, f Fields) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(msg)
	}
}

// Panic logs a message at level Panic and then panics.
func Panic(args ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Panic(args...)
	}
}

// PanicWithFields logs a message with fields at level Panic and then panics.
func PanicWithFields(msg interface{}, f Fields) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Panic(msg)
	}
}

// fileInfo returns the file and line number of the caller.
func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
