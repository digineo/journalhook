package journalhook

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var uniqueSmallMessageID = fmt.Sprintf("%d", time.Now().UnixNano())

type simpleFormatter struct{}

func (f *simpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

type Test struct {
	Field1 string
	Field2 int
	Field3 map[string]float32
}

func TestCheckSmallMessage(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	log := logrus.New()
	hook, err := NewJournalHook()
	hook.Formatter = &simpleFormatter{}
	require.NoError(err)
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"test_id": uniqueSmallMessageID,
		"struct":  Test{"field1", 2, map[string]float32{"3": 0.4, "5": 6.7}},
		"bytes":   []byte{'\n', 0xde, 0xad, 0xbe, 0xef},
	}).Info("SmallMessage")

	time.Sleep(time.Second)

	out, err := exec.Command(
		"journalctl",
		"-o", "json",
		"TEST_ID="+uniqueSmallMessageID,
	).Output()
	t.Log("journalctl output:", string(out))
	require.NoError(err)

	assert.Contains(string(out), `"MESSAGE":"SmallMessage"`)
	assert.Contains(string(out), `"PRIORITY":"6"`)
	assert.Contains(string(out), `"STRUCT":"{field1 2 map[3:0.4 5:6.7]}"`)
	assert.Contains(string(out), `"BYTES":[10,222,173,190,239]`)
}
