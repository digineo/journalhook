// Copyright 2017 Grigory Zubankov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.
//

package journalhook

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ssgreg/journald"
)

var errJournalMissing = errors.New("systemd journal is not exist")

// NormalizeFieldName returns field name acceptable by journald.
func NormalizeFieldName(s string) string {
	return strings.ToUpper(strings.Replace(s, "-", "_", -1))
}

// NewJournalHookWithLevels creates a hook to be added to an instance of logger.
// It's also allowed to specify logrus levels to fire events for.
func NewJournalHookWithLevels(levels []logrus.Level) (*JournalHook, error) {
	if journald.IsNotExist() {
		return nil, errJournalMissing
	}
	j := &journald.Journal{
		NormalizeFieldNameFn: NormalizeFieldName,
	}
	// Register an exit handler to be sure that journal connection will
	// be successfully closed and no log entries will be lost.
	// Client should call logrus.Exit() at exit.
	logrus.RegisterExitHandler(func() {
		j.Close()
	})

	return &JournalHook{
		Journal:      j,
		LogrusLevels: levels,
		Formatter:    &SyslogFormatter{},
	}, nil
}

// NewJournalHook creates a hook to be added to an instance of logger.
func NewJournalHook() (*JournalHook, error) {
	return NewJournalHookWithLevels(logrus.AllLevels)
}

// JournalHook is the systemd-journald hook for logrus.
type JournalHook struct {
	Journal      *journald.Journal
	LogrusLevels []logrus.Level
	Formatter    logrus.Formatter
}

// Fire writes a log entry to the systemd journal.
func (h *JournalHook) Fire(entry *logrus.Entry) error {
	line, _ := h.Formatter.Format(entry)

	return h.Journal.Send(string(line), levelToPriority(entry.Level), entry.Data)
}

// Levels returns a slice of Levels the hook is fired for.
func (h *JournalHook) Levels() []logrus.Level {
	return h.LogrusLevels
}

func levelToPriority(l logrus.Level) journald.Priority {
	switch l {
	case logrus.DebugLevel, logrus.TraceLevel:
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
