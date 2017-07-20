package log

import (
	"io"
	"os"
)

// TODO: Output should be like: Short Date and Time [RootPrefix][OtherPrefix]: log info

type Level int

// String will provide a string representation of the log Level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case DisabledLevel:
		return ""
	default:
		return "UNKNOWN"
	}
}

const (
	// DebugLevel messages are used when developing. These are generally not enabled in production. This is the lowest level message.
	DebugLevel Level = iota
	// InfoLevel messages are for verbose output, and should provide general information regarding the program state.
	InfoLevel
	// WarningLevel messages should be used if an abnormal condition arises, but may be handled. These are useful
	// for indicating when unexpected behaviour may be experienced.
	WarningLevel
	// ErrorLevel messages should be be used when an unhandled or unrecoverable condition occurs and may prevent
	// execution from continuing as normal. The program should recover, but expected output may not be available.
	ErrorLevel
	// DisabledLevel prevents any messages from being output.
	DisabledLevel
)

type Logger interface {
	// Debug will attempt to log a DebugLevel message. If log level is set to a higher value, these will be discarded.
	Debug(v ...interface{})
	// Info will attempt to log an InfoLevel message. If log level is set to a higher value, these will be discarded.
	Info(v ...interface{})
	// Warn will attempt to log a WarningLevel message. If log level is set to a higher value, these will be discarded.
	Warn(v ...interface{})
	// Error will attempt to log an ErrorLevel message. This is the highest and most severe logging level. If log level
	// is set to DisabledLevel, then these messages will be discarded.
	Error(v ...interface{})
}

type logger struct {
	parent Logger
	level  Level
	prefix string
	out    io.Writer
}

func (l *logger) Debug(v ...interface{}) {

}

func (l *logger) Info(v ...interface{}) {

}

func (l *logger) Warn(v ...interface{}) {

}

func (l *logger) Error(v ...interface{}) {

}

var (
	currentLevel Level = WarningLevel
	rootLogger   Logger
)

func init() {
	rootLogger = &logger{ nil, WarningLevel, "DSA", os.Stderr}
}

