package journalhook

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
	"github.com/ssgreg/journald"
)

// JournalHook is the systemd-journald hook for logrus.
type JournalHook struct {
}

// Fire writes a message to the systemd journal.
func (h *JournalHook) Fire(entry *logrus.Entry) error {
	return journald.Send(
		entry.Message,
		mapLevelToPriority(entry.Level),
		// Journal wants uppercase strings.
		stringifyEntries(entry.Data))
}

// Levels returns a slice of Levels the hook is fired for.
func (h *JournalHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Enable adds the Journal hook if journal is enabled.
// Sets log output to ioutil.Discard so stdout isn't captured.
func Enable() {
	if journald.IsNotExists() {
		logrus.Warning("Journal not available but user requests we log to it. Ignoring")
	} else {
		logrus.AddHook(&JournalHook{})
		logrus.SetOutput(ioutil.Discard)
	}
}

func mapLevelToPriority(l logrus.Level) journald.Priority {
	switch l {
	case logrus.DebugLevel:
		return journald.PriorityDebug
	case logrus.InfoLevel:
		return journald.PriorityInfo
	case logrus.WarnLevel:
		return journald.PriorityWarning
	case logrus.ErrorLevel:
		return journald.PriorityErr
	case logrus.FatalLevel:
		return journald.PriorityCrit
	case logrus.PanicLevel:
		return journald.PriorityEmerg
	}
	return journald.PriorityNotice
}

func stringifyOp(r rune) rune {
	switch {
	case r >= 'A' && r <= 'Z':
		return r
	case unicode.IsDigit(r):
		return r
	case r >= 'a' && r <= 'z':
		return unicode.ToUpper(r)
	default:
		return rune('_')
	}
}

func stringifyKey(key string) string {
	key = strings.Map(stringifyOp, key)
	if strings.HasPrefix(key, "_") {
		key = strings.TrimPrefix(key, "_")
	}
	return key
}

// Journal wants strings but logrus takes anything.
func stringifyEntries(data map[string]interface{}) map[string]string {
	entries := make(map[string]string)
	for k, v := range data {
		key := stringifyKey(k)
		entries[key] = fmt.Sprint(v)
	}
	return entries
}
