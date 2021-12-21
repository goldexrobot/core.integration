package helpers

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type Consoler interface {
	Errorf(f string, args ...interface{})
	Warningf(f string, args ...interface{})
	Debugf(f string, args ...interface{})
	Printf(f string, args ...interface{})
}

type LogrusConsoleHook struct {
	Consoler Consoler
	NoData   bool
}

func (h LogrusConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h LogrusConsoleHook) Fire(e *logrus.Entry) error {
	msg := fmt.Sprintf("%-45s%s", e.Message, h.data(e))
	switch e.Level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		h.Consoler.Errorf(msg)
	case logrus.WarnLevel:
		h.Consoler.Warningf(msg)
	case logrus.InfoLevel:
		h.Consoler.Printf(msg)
	case logrus.DebugLevel, logrus.TraceLevel:
		h.Consoler.Debugf(msg)
	default:
		h.Consoler.Printf(msg)
	}
	return nil
}

func (h LogrusConsoleHook) data(e *logrus.Entry) string {
	if h.NoData || len(e.Data) == 0 {
		return ""
	}

	var keys []string
	for k := range e.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sw := bytes.NewBuffer(nil)
	for _, k := range keys {
		if k == "" {
			continue
		}
		sw.WriteString("[")
		sw.WriteString(k)
		val := fmt.Sprint(e.Data[k])
		if val != "" {
			sw.WriteString("=")
			sw.WriteString(val)
		}
		sw.WriteString("] ")
	}

	return "   " + color.HiBlackString(strings.TrimSpace(sw.String()))
}
