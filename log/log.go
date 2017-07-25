package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Level represents one of the levels of which logs may be stored at. If a log is called on a level
// below the current level then it will be discarded.
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

var (
	currentLevel Level = WarningLevel
	rootLogger   *Logger
)

// Logger represents an active logging object that generates lines of output to an io.Writer.
// Each logging operation makes a single call to the Writer's Write method. A Logger can be used
// simultaneously from multiple goroutines. It guarantees to serialize access to the Writer.
type Logger struct {
	mu     sync.Mutex
	parent *Logger
	prefix string
	buf    []byte
	out    io.Writer
}

// Debug will attempt to log a DebugLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Debug(v ...interface{}) {
	if currentLevel > DebugLevel {
		return
	}

	s := fmt.Sprintf(" %s", fmt.Sprint(v...))
	l.writeBuf([]byte(s), DebugLevel)
}

// Info will attempt to log an InfoLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Info(v ...interface{}) {
	if currentLevel > InfoLevel {
		return
	}

	s := fmt.Sprintf(" %s", fmt.Sprint(v...))
	l.writeBuf([]byte(s), InfoLevel)
}

// Warn will attempt to log a WarningLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Warn(v ...interface{}) {
	if currentLevel > WarningLevel {
		return
	}

	s := fmt.Sprintf(" %s", fmt.Sprint(v...))
	l.writeBuf([]byte(s), WarningLevel)
}

// Error will attempt to log an ErrorLevel message. This is the highest and most severe logging level. If log level
// is set to DisabledLevel, then these messages will be discarded.
func (l *Logger) Error(v ...interface{}) {
	if currentLevel > ErrorLevel {
		return
	}

	s := fmt.Sprintf(" %s", fmt.Sprint(v...))
	l.writeBuf([]byte(s), ErrorLevel)
}

func (l *Logger) write(lvl Level) {
	if l.out == nil {
		panic("logger has no output destination")
	}

	now := time.Now()
	var b []byte
	b = append(b, fmt.Sprintf("[%s] ", lvl)...)
	b = append(b, now.Format(time.StampMilli)...)
	b = append(b, ' ')
	b = append(b, l.buf...)
	if b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	l.out.Write(b)
}

func (l *Logger) writeBuf(b []byte, lvl Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.addPrefix()
	l.buf = append(l.buf, b...)

	if l.parent != nil {
		l.parent.writeBuf(l.buf, lvl)
	} else {
		l.write(lvl)
	}

	l.buf = l.buf[:0]
}

func (l *Logger) addPrefix() {
	if len(l.prefix) >= 1 {
		l.buf = append(l.buf, '[')
		l.buf = append(l.buf, l.prefix...)
		l.buf = append(l.buf, ']')
	}
}

func init() {
	rootLogger = &Logger{
		parent: nil,
		prefix: "DSA",
		out:    os.Stderr,
	}
}

// SetLevel will specify the global logging level. Any logging calls below the level specified will be ignored.
// The default Level is WarningLevel
func SetLevel(level Level) {
	currentLevel = level
}

// SetOutput will change the output Writer of all loggers to be the one specified. The logger will _not_
// close the Writer and must be closed elsewhere. The default Output destination is os.Stderr
func SetOutput(w io.Writer) {
	rootLogger.out = w
}

// New will create a new logger, with the specified prefix. The prefix will be added the start of every message at
// ever log level. You can also specify a parent logger, which will affix the prefix of the parent before the prefix
// of the returned logger. If parent is `nil` then the parent will be the root DSA logger, and will prefix the message
// messages at each log level with `[LogLevel] Date/Time Stamp [DSA]`
func New(prefix string, parent *Logger) *Logger {
	l := &Logger{prefix: prefix}
	if parent == nil {
		l.parent = rootLogger
	} else {
		l.parent = parent
	}

	return l
}

// Debug will attempt to log a DebugLevel message on the root logger. If log level is set to a higher value, these will be discarded.
func Debug(v ...interface{}) {
	rootLogger.Debug(v...)
}

// Info will attempt to log an InfoLevel message on the root logger. If log level is set to a higher value, these will be discarded.
func Info(v ...interface{}) {
	rootLogger.Info(v...)
}

// Warn will attempt to log a WarningLevel message on the root logger. If log level is set to a higher value, these will be discarded.
func Warn(v ...interface{}) {
	rootLogger.Warn(v...)
}

// Error will attempt to log an ErrorLevel message on the root logger. This is the highest and most severe logging level. If log level
// is set to DisabledLevel, then these messages will be discarded.
func Error(v ...interface{}) {
	rootLogger.Error(v...)
}
