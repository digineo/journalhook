package journalhook

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"sort"

	"github.com/sirupsen/logrus"
)

const FieldKeyPrefix = "prefix"

// SyslogFormatter formats logs into text
type SyslogFormatter struct {
	// Force quoting of all values
	ForceQuote bool

	// DisableQuote disables quoting for all values.
	// DisableQuote will have a lower priority than ForceQuote.
	// If both of them are set to true, quote will be forced on all values.
	DisableQuote bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)
}

// Format renders a single log entry
func (f *SyslogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var prefix string

	data := make(logrus.Fields)
	keys := make([]string, 0, len(data))
	for k, v := range entry.Data {
		if k == logrus.FieldKeyTime || k == logrus.FieldKeyLevel {
			continue // skip this field
		}
		if k == FieldKeyPrefix {
			prefix = v.(string)
			continue
		}

		data[k] = v
		keys = append(keys, k)
	}

	var funcVal, fileVal string

	fixedKeys := make([]string, 0, 4+len(data))

	if entry.Message != "" {
		fixedKeys = append(fixedKeys, logrus.FieldKeyMsg)
	}

	if entry.HasCaller() {
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		} else {
			funcVal = entry.Caller.Function
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}

		if funcVal != "" {
			fixedKeys = append(fixedKeys, logrus.FieldKeyFunc)
		}
		if fileVal != "" {
			fixedKeys = append(fixedKeys, logrus.FieldKeyFile)
		}
	}

	// Sort
	if f.SortingFunc == nil {
		sort.Strings(keys)
		fixedKeys = append(fixedKeys, keys...)
	} else {
		fixedKeys = append(fixedKeys, keys...)
		f.SortingFunc(fixedKeys)
	}

	b := &bytes.Buffer{}
	if prefix != "" {
		fmt.Fprintf(b, "[%s]", prefix)
	}

	for _, key := range fixedKeys {
		var value interface{}
		switch {
		case logrus.FieldKeyFunc != "" && entry.HasCaller():
			value = funcVal
		case logrus.FieldKeyFile != "" && entry.HasCaller():
			value = fileVal
		default:
			value = data[key]
		}

		if value != nil {
			f.appendKeyValue(b, key, value)
		}
	}

	if entry.Message != "" {
		f.appendKeyValue(b, "msg", entry.Message)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *SyslogFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *SyslogFormatter) appendValue(b io.StringWriter, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func (f *SyslogFormatter) needsQuoting(text string) bool {
	if f.ForceQuote {
		return true
	}
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.DisableQuote {
		return false
	}

	return hasSpecialCharacter(text)
}

func hasSpecialCharacter(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}
