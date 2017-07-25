package log

import (
	"bytes"
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
	buf    *bytes.Buffer
	out    io.Writer
}

// Debug will attempt to log a DebugLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Debug(v ...interface{}) {
	if currentLevel > DebugLevel {
		return
	}

	buf := bytes.NewBufferString(" ")
	fmt.Fprint(buf, v...)
	l.writeBuf(buf.Bytes(), DebugLevel)
}

// Debugf will attempt to log a formatted DebugLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if currentLevel > DebugLevel {
		return
	}

	buf := bytes.NewBufferString(" ")
	fmt.Fprintf(buf, format, v...)
	l.writeBuf(buf.Bytes(), DebugLevel)
}

// Info will attempt to log an InfoLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Info(v ...interface{}) {
	if currentLevel > InfoLevel {
		return
	}

	buf := bytes.NewBufferString(" ")
	fmt.Fprint(buf, v...)
	l.writeBuf(buf.Bytes(), InfoLevel)
}

// Infof will attempt to log a formatted InfoLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Infof(format string, v ...interface{}) {
	if currentLevel > InfoLevel {
		return
	}

	buf := bytes.NewBufferString(" ")
	fmt.Fprintf(buf, format, v...)
	l.writeBuf(buf.Bytes(), InfoLevel)
}

// Warn will attempt to log a WarningLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Warn(v ...interface{}) {
	if currentLevel > WarningLevel {
		return
	}

	buf := bytes.NewBufferString(" ")
	fmt.Fprint(buf, v...)
	l.writeBuf(buf.Bytes(), WarningLevel)
}

// Warnf will attempt to log a formatted WarningLevel message. If log level is set to a higher value, these will be discarded.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if currentLevel > WarningLevel {
		return
	}

	buf := bytes.NewBufferString(" ")
	fmt.Fprintf(buf, format, v...)
	l.writeBuf(buf.Bytes(), WarningLevel)
}

// Error will attempt to log an ErrorLevel message. This is the highest and most severe logging level. If log level
// is set to DisabledLevel, then these messages will be discarded.
func (l *Logger) Error(v ...interface{}) {
	if currentLevel > ErrorLevel {
		return
	}

	buf := bytes.NewBufferString(" ")
	fmt.Fprint(buf, v...)
	l.writeBuf(buf.Bytes(), ErrorLevel)
}

// Errorf will attempt to log a formatted ErrorLevel message. This is the highest and most severe logging level. If log level
// is set to DisabledLevel, then these messages will be discarded.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if currentLevel > ErrorLevel {
		return
	}

	buf := bytes.NewBufferString(" ")
	fmt.Fprintf(buf, format, v...)
	l.writeBuf(buf.Bytes(), ErrorLevel)
}

func (l *Logger) write(lvl Level) {
	if l.out == nil {
		panic("logger has no output destination")
	}

	now := time.Now()
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "[%s] ", lvl)
	fmt.Fprint(b, now.Format(time.StampMilli))
	b.WriteByte(' ')
	l.buf.WriteTo(b)
	if b.Bytes()[len(b.Bytes())-1] != '\n' {
		b.WriteByte('\n')
	}

	b.WriteTo(l.out)
}

func (l *Logger) writeBuf(b []byte, lvl Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.addPrefix()
	l.buf.Write(b)

	if l.parent != nil {
		l.parent.writeBuf(l.buf.Bytes(), lvl)
	} else {
		l.write(lvl)
	}

	l.buf.Reset()
}

func (l *Logger) addPrefix() {
	if l.buf == nil {
		l.buf = new(bytes.Buffer)
	}
	if len(l.prefix) >= 1 {
		l.buf.WriteByte('[')
		l.buf.WriteString(l.prefix)
		l.buf.WriteByte(']')
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

// Debugf will attempt to log a formatted DebugLevel message on the root logger. If log level is set to a higher value, these will be discarded.
func Debugf(format string, v ...interface{}) {
	rootLogger.Debugf(format, v...)
}

// Infof will attempt to log a formatted InfoLevel message on the root logger. If log level is set to a higher value, these will be discarded.
func Infof(format string, v ...interface{}) {
	rootLogger.Infof(format, v...)
}

// Warnf will attempt to log a formatted WarningLevel message on the root logger. If log level is set to a higher value, these will be discarded.
func Warnf(format string, v ...interface{}) {
	rootLogger.Warnf(format, v...)
}

// Errorf will attempt to log a formatted ErrorLevel message on the root logger. This is the highest and most severe logging level. If log level
// is set to DisabledLevel, then these messages will be discarded.
func Errorf(format string, v ...interface{}) {
	rootLogger.Errorf(format, v...)
}
