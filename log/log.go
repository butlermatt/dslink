package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

// Time will be formatted as YYYY-MM-DD HH:MM:SS:ssss
const formatStr = "2006-01-02 15:04:05.000"

// Level represents one of the levels of which logs may be stored at. If a log is called on a level
// below the current level then it will be discarded.
type Level int

// String will provide a string representation of the log Level
func (l Level) String() string {
	switch l {
	case TraceLvl:
		return "TRACE"
	case DebugLvl:
		return "DEBUG"
	case FineLvl:
		return "FINE "
	case WarningLvl:
		return "WARN "
	case InfoLvl:
		return "INFO "
	case ErrorLvl:
		return "ERROR"
	case AdminLvl:
		return "ADMIN"
	case FatalLvl:
		return "FATAL"
	case DisabledLvl:
		return ""
	default:
		return "UNKNOWN"
	}
}

const (
	// TraceLvl messages are used when debugging or developing. They are generally not enabled in production. This is the lowest level message.
	TraceLvl Level = iota
	// DebugLvl messages are used when developing. These are generally not enabled in production.
	DebugLvl
	// FineLvl messages are operating detail messages. Not generally exposed but useful for high level debugging.
	FineLvl
	// WarningLvl messages should be used if an abnormal condition arises, but may be handled. These are useful
	// for indicating when unexpected behaviour may be experienced.
	WarningLvl
	// InfoLvl messages are for verbose output, and should provide general information regarding the program state.
	InfoLvl
	// ErrorLvl messages should be be used when an unhandled or unrecoverable condition occurs and may prevent
	// execution from continuing as normal. The program should recover, but expected output may not be available.
	ErrorLvl
	// AdminLvl messages are high level warning or error messages that many need an administer to resolve.
	AdminLvl
	// FatalLvl messages are worst case messages which should be sent when something fails completely.
	FatalLvl
	// DisabledLvl prevents any messages from being output.
	DisabledLvl
)

var (
	currentLevel Level = WarningLvl
	rootLogger   *Logger
	ch           chan *LogRecord
	out          io.Writer
)

type LogRecord struct {
	// Time will be formatted as YYYY-MM-DD HH:MM:SS:ssss
	Time       time.Time
	Level      Level
	LoggerName string
	Format     string
	Args       []interface{}
}

func newRecord(lvl Level, logger, format string, args ...interface{}) *LogRecord {
	return &LogRecord{
		Time:       time.Now(),
		Level:      lvl,
		LoggerName: logger,
		Format:     format,
		Args:       args,
	}
}

func (lr *LogRecord) String() string {
	var buf bytes.Buffer

	buf.WriteByte('[')
	buf.WriteString(lr.Time.Format(formatStr))
	buf.WriteByte(']')
	buf.WriteString(" " + lr.Level.String() + " ")
	buf.WriteString("[" + lr.LoggerName + "] ")
	buf.WriteString(fmt.Sprintf(lr.Format, lr.Args...))

	return buf.String()
}

// Logger represents an active logging object that generates lines of output to an io.Writer.
// Each logging operation makes a single call to the Writer's Write method. A Logger can be used
// simultaneously from multiple goroutines. It guarantees to serialize access to the Writer.
type Logger struct {
	Name string
}

// New creates a new Logger with the specified name.
func New(name string) *Logger {
	return &Logger{Name: name}
}

// Logf creates a new log entry at the specified level with the format string specified for this logger.
func (l *Logger) Logf(lvl Level, format string, args ...interface{}) {
	r := newRecord(lvl, l.Name, format, args...)
	ch <- r
}

// Trace creates a Trace Level log entry with the specified string.
func (l *Logger) Trace(message string) {
	l.Logf(TraceLvl, message)
}

// Tracef creates a Trace Level log entry with the format string specified for this logger.
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Logf(TraceLvl, format, args...)
}

// Debug creates a Debug Level log entry with the specified string.
func (l *Logger) Debug(message string) {
	l.Logf(DebugLvl, message)
}

// Debugf creates a Debug Level log entry with the format string specified for this logger.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logf(DebugLvl, format, args...)
}

// Fine creates a Fine Level log entry with the specified string.
func (l *Logger) Fine(message string) {
	l.Logf(FineLvl, message)
}

// Finef creates a Fine Level log entry with the format string specified for this logger.
func (l *Logger) Finef(format string, args ...interface{}) {
	l.Logf(FineLvl, format, args...)
}

// Warn creates a Warn Level log entry with the specified string.
func (l *Logger) Warn(message string) {
	l.Logf(WarningLvl, message)
}

// Warnf creates a Warning Level log entry with the format string specified for this logger.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logf(WarningLvl, format, args...)
}

// Info creates a Info Level log entry with the specified string.
func (l *Logger) Info(message string) {
	l.Logf(InfoLvl, message)
}

// Infof creates a Info Level log entry with the format string specified for this logger.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logf(InfoLvl, format, args...)
}

// Error creates a Error Level log entry with the specified string.
func (l *Logger) Error(message string) {
	l.Logf(ErrorLvl, message)
}

// Errorf creates a Error Level log entry with the format string specified for this logger.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(ErrorLvl, format, args...)
}

// Admin creates a Admin Level log entry with the specified string.
func (l *Logger) Admin(message string) {
	l.Logf(AdminLvl, message)
}

// Adminf creates a Admin Level log entry with the format string specified for this logger.
func (l *Logger) Adminf(format string, args ...interface{}) {
	l.Logf(AdminLvl, format, args...)
}

// Fatal creates a Fatal Level log entry with the specified string.
func (l *Logger) Fatal(message string) {
	l.Logf(FatalLvl, message)
}

// Fatalf creates a Fatal Level log entry with the format string specified for this logger.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logf(FatalLvl, format, args...)
}

type LogHandler func(record *LogRecord)

func init() {
	ch = make(chan *LogRecord, 10)
	out = os.Stdout

	go printLog()
}

func printLog() {
	for r := range ch {
		_, _ = fmt.Fprint(out, r.String())
	}
}
