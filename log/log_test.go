package log

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

type syncWriter struct {
	buf *bytes.Buffer
	ch  chan bool
}

func newWriter() *syncWriter {
	return &syncWriter{buf: new(bytes.Buffer), ch: make(chan bool)}
}

func (sw *syncWriter) Write(p []byte) (n int, err error) {
	n, err = sw.buf.Write(p)
	sw.ch <- true
	return n, err
}

func (sw *syncWriter) String() string {
	return sw.buf.String()
}

func TestSetLevel(t *testing.T) {
	if rootLogger.level != WarningLvl {
		t.Errorf("Incorrect default level. expected=%q got=%q", WarningLvl, rootLogger.level)
	}

	SetLevel(DisabledLvl)

	if rootLogger.level != DisabledLvl {
		t.Errorf("SetLevel failed. expect=%q got=%q", DisabledLvl, rootLogger.level)
	}

	SetLevel(WarningLvl) // Reset logging
}

func TestSetOutput(t *testing.T) {
	if out != os.Stdout {
		t.Errorf("Incorrect default output. expect=%#v got=%#v", os.Stdout, out)
	}

	buf := newWriter()
	SetOutput(buf)
	Warn("ignore")
	<-buf.ch

	if out != buf {
		t.Errorf("SetOutput failed. expect=%#v got=%#v", buf, out)
	}
}

func TestNewChild(t *testing.T) {
	buf := newWriter()

	log := New("Test")
	SetOutput(buf)
	log.SetLevel(DebugLvl)

	log.Debug("This is a test")

	<-buf.ch
	str := buf.String()

	if !strings.Contains(str, "DSA") {
		t.Errorf("Logged line does not contain root prefix: %q", str)
	}

	if !strings.Contains(str, "Test") {
		t.Error("Logged line does not contain prefix.")
	}

	buf.buf.Reset()

	l2 := log.Child("Test2")
	l2.Debug("This is a test")

	<-buf.ch
	str = buf.String()

	if !strings.Contains(str, "[DSA.Test.Test2]") {
		t.Errorf("unexpected logger name. expected=%q, got=%q", "DSA.Test.Test2", l2.name)
	}
}

type logLevelTest struct {
	level  Level
	action func(string)
}

type formattedLogTest struct {
	level  Level
	action func(format string, args ...interface{})
}

func TestLogger_Levels(t *testing.T) {
	l := New("Test")
	tests := []logLevelTest{
		{TraceLvl, l.Trace},
		{DebugLvl, l.Debug},
		{FineLvl, l.Fine},
		{WarningLvl, l.Warn},
		{InfoLvl, l.Info},
		{ErrorLvl, l.Error},
		{AdminLvl, l.Admin},
		{FatalLvl, l.Fatal},
	}

	formatTests := []formattedLogTest{
		{TraceLvl, l.Tracef},
		{DebugLvl, l.Debugf},
		{FineLvl, l.Finef},
		{WarningLvl, l.Warnf},
		{InfoLvl, l.Infof},
		{ErrorLvl, l.Errorf},
		{AdminLvl, l.Adminf},
		{FatalLvl, l.Fatalf},
	}

	testLevels(t, l, tests, "DSA.Test")
	testFormattedLevels(t, l, formatTests, "DSA.Test")
}

func TestRootLevels(t *testing.T) {
	tests := []logLevelTest{
		{TraceLvl, Trace},
		{DebugLvl, Debug},
		{FineLvl, Fine},
		{WarningLvl, Warn},
		{InfoLvl, Info},
		{ErrorLvl, Error},
		{AdminLvl, Admin},
		{FatalLvl, Fatal},
	}

	formatTests := []formattedLogTest{
		{TraceLvl, Tracef},
		{DebugLvl, Debugf},
		{FineLvl, Finef},
		{WarningLvl, Warnf},
		{InfoLvl, Infof},
		{ErrorLvl, Errorf},
		{AdminLvl, Adminf},
		{FatalLvl, Fatalf},
	}

	testLevels(t, rootLogger, tests, "DSA")
	testFormattedLevels(t, rootLogger, formatTests, "DSA")
}

func testLevels(t *testing.T, logger *Logger, tests []logLevelTest, name string) {
	t.Helper()
	levels := []Level{TraceLvl, DebugLvl, FineLvl, WarningLvl, InfoLvl, ErrorLvl, AdminLvl, FatalLvl, DisabledLvl}
	buf := newWriter()
	SetOutput(buf)

	for _, lvl := range levels {
		for _, tt := range tests {
			buf.buf.Reset()
			logger.SetLevel(lvl)
			tt.action("this is a test")
			if tt.level >= lvl {
				<-buf.ch
				s := buf.String()
				if !strings.Contains(s, fmt.Sprintf("%s", tt.level)) {
					t.Errorf("Logged line does not contain level tag: %q, got %s", tt.level, s)
				}
				if !strings.Contains(s, "DSA") {
					t.Error("Logged line does not contain root prefix")
				}
				if !strings.Contains(s, name) {
					t.Errorf("Logged line does not contain expected prefix: %q, got=%q", name, s)
				}
				if !strings.Contains(s, "this is a test") {
					t.Errorf("Logged line does not contain provided. expected=%q got=%q", "this is a test", s)
				}
			} else {
				go func() { buf.ch <- true }()
				<-buf.ch
				s := buf.String()
				if len(s) != 0 {
					t.Errorf("logged line contains unexpected content: %q", s)
				}
			}
		}
	}
}

func testFormattedLevels(t *testing.T, logger *Logger, tests []formattedLogTest, name string) {
	t.Helper()
	levels := []Level{TraceLvl, DebugLvl, FineLvl, WarningLvl, InfoLvl, ErrorLvl, AdminLvl, FatalLvl, DisabledLvl}
	buf := newWriter()
	SetOutput(buf)

	for _, lvl := range levels {
		for _, tt := range tests {
			buf.buf.Reset()
			logger.SetLevel(lvl)
			tt.action("this is a test")
			if tt.level >= lvl {
				<-buf.ch
				s := buf.String()
				if !strings.Contains(s, fmt.Sprintf("%s", tt.level)) {
					t.Errorf("Logged line does not contain level tag: %q, got %s", tt.level, s)
				}
				if !strings.Contains(s, "DSA") {
					t.Error("Logged line does not contain root prefix")
				}
				if !strings.Contains(s, name) {
					t.Errorf("Logged line does not contain expected prefix: %q, got=%q", name, s)
				}
				if !strings.Contains(s, "this is a test") {
					t.Errorf("Logged line does not contain provided. expected=%q got=%q", "this is a test", s)
				}
			} else {
				go func() { buf.ch <- true }()
				<-buf.ch
				s := buf.String()
				if len(s) != 0 {
					t.Errorf("logged line contains unexpected content: %q", s)
				}
			}
		}
	}
}
