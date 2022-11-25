package logger

import (
	"fmt"
	"os"
)

var l *logger

// Logger ...
type Logger interface {
	Info(m ...interface{})
	Debug(m ...interface{})
	Warn(m ...interface{})
	Error(m ...interface{})
	Fatal(m ...interface{})

	Debugf(format string, m ...interface{})
	Infof(format string, m ...interface{})
	Warnf(format string, m ...interface{})
	Errorf(format string, m ...interface{})
	Fatalf(format string, m ...interface{})

	GetLogLevel() string
}

// GetLogger ...
func GetLogger(defaultLogLevel LogLevel) Logger {
	if l != nil {
		return l
	}
	l = &logger{
		level: defaultLogLevel,
	}
	l.Debugf("log level set to '%s'", l.GetLogLevel())
	return l
}

// LogLevel ...
type LogLevel int

// LogLevels ...
const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type logger struct {
	level LogLevel
}

func (l *logger) Debug(m ...interface{}) {
	if l.level > LogLevelDebug {
		return
	}
	fmt.Fprintln(os.Stdout, append([]interface{}{"[DEBUG]"}, m...)...)
}

func (l *logger) Debugf(format string, m ...interface{}) {
	if l.level > LogLevelDebug {
		return
	}
	fs := fmt.Sprintf(format, m...)
	fmt.Fprintln(os.Stdout, "[DEBUG]", fs)
}

func (l *logger) Info(m ...interface{}) {
	if l.level > LogLevelInfo {
		return
	}
	fmt.Fprintln(os.Stdout, append([]interface{}{"[INFO ]"}, m...)...)
}

func (l *logger) Infof(format string, m ...interface{}) {
	if l.level > LogLevelInfo {
		return
	}
	fs := fmt.Sprintf(format, m...)
	fmt.Fprintln(os.Stdout, "[INFO ]", fs)
}

func (l *logger) Warn(m ...interface{}) {
	if l.level > LogLevelWarn {
		return
	}
	fmt.Fprintln(os.Stdout, append([]interface{}{"[WARN ]"}, m...)...)
}

func (l *logger) Warnf(format string, m ...interface{}) {
	if l.level > LogLevelWarn {
		return
	}
	fs := fmt.Sprintf(format, m...)
	fmt.Fprintln(os.Stdout, "[WARN ]", fs)
}

func (l *logger) Error(m ...interface{}) {
	if l.level > LogLevelError {
		return
	}
	fmt.Fprintln(os.Stdout, append([]interface{}{"[ERROR]"}, m...)...)
}

func (l *logger) Errorf(format string, m ...interface{}) {
	if l.level > LogLevelError {
		return
	}
	fs := fmt.Sprintf(format, m...)
	fmt.Fprintln(os.Stdout, "[ERROR]", fs)
}

func (l *logger) Fatal(m ...interface{}) {
	fmt.Fprintln(os.Stderr, append([]interface{}{"[FATAL]"}, m...)...)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, m ...interface{}) {
	fs := fmt.Sprintf(format, m...)
	fmt.Fprintln(os.Stdout, "[FATAL]", fs)
}

func (l *logger) GetLogLevel() string {
	switch l.level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNDEFINED"
	}
}
