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

// NewJournalHook creates a hook to be added to an instance of logger.
func NewJournalHook() (*JournalHook, error) {
	if journald.IsNotExist() {
		return nil, errors.New("systemd journal is not exist")
	}
	return &JournalHook{journald.Journal{
		NormalizeFieldNameFn: strings.ToUpper,
	}}, nil
}

// JournalHook is the systemd-journald hook for logrus.
type JournalHook struct {
	Journal journald.Journal
}

// Fire writes a message to the systemd journal.
func (h *JournalHook) Fire(entry *logrus.Entry) error {
	return h.Journal.Send(entry.Message, levelToPriority(entry.Level), entry.Data)
}

// Levels returns a slice of Levels the hook is fired for.
func (h *JournalHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func levelToPriority(l logrus.Level) journald.Priority {
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
