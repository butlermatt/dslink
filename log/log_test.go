package log

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestSetLevel(t *testing.T) {
	if currentLevel != WarningLvl {
		t.Errorf("Incorrect default level. expected=%q got=%q", WarningLvl, currentLevel)
	}

	SetLevel(DisabledLvl)

	if currentLevel != DisabledLvl {
		t.Errorf("SetLevel failed. expect=%q got=%q", DisabledLvl, currentLevel)
	}
}

func TestSetOutput(t *testing.T) {
	if rootLogger.out != os.Stderr {
		t.Errorf("Incorrect default output. expect=%#v got=%#v", os.Stderr, rootLogger.out)
	}

	buf := new(bytes.Buffer)
	SetOutput(buf)

	if rootLogger.out != buf {
		t.Errorf("SetOutput failed. expect=%#v got=%#v", buf, rootLogger.out)
	}
}

func TestNew(t *testing.T) {
	buf := new(bytes.Buffer)

	SetOutput(buf)
	SetLevel(DebugLvl)

	log := New("Test", nil)
	log.Debug("This is a test", 15)

	str := buf.String()

	if !strings.Contains(str, "[DSA]") {
		t.Error("Logged line does not contain root prefix")
	}

	if !strings.Contains(str, "[Test]") {
		t.Error("Logged line does not contain prefix.")
	}

	buf.Reset()

	l2 := New("Test2", log)
	l2.Debug("This is a test")

	str = buf.String()

	if !strings.Contains(str, "[DSA]") {
		t.Error("Logged line does not contain root prefix")
	}

	if !strings.Contains(str, "[Test]") {
		t.Error("Logged line does not contain parent prefix.")
	}

	if !strings.Contains(str, "[Test2]") {
		t.Error("Logged line does not contain prefix")
	}
}

func TestLogger_Levels(t *testing.T) {
	l := New("Test", nil)
	testLevels(t, DebugLvl, l.Debug, l.prefix)
	testLevels(t, InfoLvl, l.Info, l.prefix)
	testLevels(t, WarningLvl, l.Warn, l.prefix)
	testLevels(t, ErrorLvl, l.Error, l.prefix)
	testFormattedLevels(t, DebugLvl, l.Debugf, l.prefix)
	testFormattedLevels(t, InfoLvl, l.Infof, l.prefix)
	testFormattedLevels(t, WarningLvl, l.Warnf, l.prefix)
	testFormattedLevels(t, ErrorLvl, l.Errorf, l.prefix)
}

func TestRootLevels(t *testing.T) {
	testLevels(t, DebugLvl, Debug, "DSA")
	testLevels(t, InfoLvl, Info, "DSA")
	testLevels(t, WarningLvl, Warn, "DSA")
	testLevels(t, ErrorLvl, Error, "DSA")
	testFormattedLevels(t, DebugLvl, Debugf, "DSA")
	testFormattedLevels(t, InfoLvl, Infof, "DSA")
	testFormattedLevels(t, WarningLvl, Warnf, "DSA")
	testFormattedLevels(t, ErrorLvl, Errorf, "DSA")
}

func testLevels(t *testing.T, l Level, f func(v ...interface{}), prefix string) {
	levels := []Level{DebugLvl, InfoLvl, WarningLvl, ErrorLvl, DisabledLvl}
	line := "This is a test"
	buf := new(bytes.Buffer)
	SetOutput(buf)

	for _, lvl := range levels {
		buf.Reset()
		SetLevel(lvl)
		f(line)
		s := buf.String()
		if lvl <= l {
			if !strings.Contains(s, fmt.Sprintf("[%s]", l)) {
				t.Errorf("Logged line does not contain level tag: %q, got %s", l, s)
			}
			if !strings.Contains(s, "[DSA]") {
				t.Error("Logged line does not contain root prefix")
			}
			if !strings.Contains(s, fmt.Sprintf("[%s]", prefix)) {
				t.Errorf("Logged line does not contain expected prefix: %q, got=%q", prefix, s)
			}
			if !strings.Contains(s, fmt.Sprintf("%s\n", line)) {
				t.Errorf("Logged line does not contain provided. expected=%q got=%q", line, s)
			}
		} else {
			if len(s) > 0 {
				t.Errorf("Unexpected logged line at %q level: %s", lvl, s)
			}
		}
	}
}

func testFormattedLevels(t *testing.T, l Level, f func(format string, v ...interface{}), prefix string) {
	levels := []Level{DebugLvl, InfoLvl, WarningLvl, ErrorLvl, DisabledLvl}
	line := "This is a %s"
	test := "test"
	buf := new(bytes.Buffer)
	SetOutput(buf)

	for _, lvl := range levels {
		buf.Reset()
		SetLevel(lvl)
		f(line, test)
		s := buf.String()
		if lvl <= l {
			if !strings.Contains(s, fmt.Sprintf("[%s]", l)) {
				t.Errorf("Logged line does not contain level tag: %q, got %s", l, s)
			}
			if !strings.Contains(s, "[DSA]") {
				t.Error("Logged line does not contain root prefix")
			}
			if !strings.Contains(s, fmt.Sprintf("[%s]", prefix)) {
				t.Errorf("Logged line does not contain expected prefix: %q, got=%q", prefix, s)
			}
			if !strings.Contains(s, fmt.Sprintf("%s\n", fmt.Sprintf(line, test))) {
				t.Errorf("Logged line does not contain provided line. expected=%q got=%q", line, s)
			}
		} else {
			if len(s) > 0 {
				t.Errorf("Unexpected logged line at %q level: %s", lvl, s)
			}
		}
	}
}
