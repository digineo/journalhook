package journalhook

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var uniqueSmallMessageID = "1"

func TestMain(m *testing.M) {
	uniqueSmallMessageID = time.Now().String()

	os.Exit(m.Run())
}

type Test struct {
	Field1 string
	Field2 int
	Field3 map[string]float32
}

func TestCheckSmallMessage(t *testing.T) {
	log := logrus.New()
	hook, err := NewJournalHook()
	require.NoError(t, err)
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"test_id": uniqueSmallMessageID,
		"struct":  Test{"field1", 2, map[string]float32{"3": 0.4, "5": 6.7}},
		"bytes":   []byte{'\n', 0xde, 0xad, 0xbe, 0xef},
	}).Info("SmallMessage")

	time.Sleep(time.Second * 5)

	out, err := exec.Command("sh", "-c", fmt.Sprintf("journalctl 'TEST_ID=%s' -o json", uniqueSmallMessageID)).Output()
	fmt.Println(string(out))

	require.NoError(t, err)
	require.True(t, strings.Contains(string(out), "MESSAGE\" : \"SmallMessage"))
	require.True(t, strings.Contains(string(out), "PRIORITY\" : \"6"))
	require.True(t, strings.Contains(string(out), "STRUCT\" : \"{field1 2 map[3:0.4 5:6.7]}"))
	require.True(t, strings.Contains(string(out), "BYTES\" : [ 10, 222, 173, 190, 239 ]"))
}

func TestCheckErrMessage(t *testing.T) {
	log := logrus.New()
	hook, err := NewJournalHookWithErrToMsg()
	require.NoError(t, err)
	log.Hooks.Add(hook)

	err = errors.New("something something dark side")

	log.WithField("test_id", uniqueSmallMessageID).WithError(err).Error()

	time.Sleep(time.Second * 5)

	out, err := exec.Command("sh", "-c", fmt.Sprintf("journalctl 'TEST_ID=%s' -o json", uniqueSmallMessageID)).Output()
	fmt.Println(string(out))

	require.NoError(t, err)
	require.True(t, strings.Contains(string(out), "MESSAGE\" : \"something something dark side"))
	require.True(t, strings.Contains(string(out), "PRIORITY\" : \"3"))
}

func TestCheckErrMessageWithMessage(t *testing.T) {
	log := logrus.New()
	hook, err := NewJournalHookWithErrToMsg()
	require.NoError(t, err)
	log.Hooks.Add(hook)

	err = errors.New("something something dark side")

	log.WithField("test_id", uniqueSmallMessageID).WithError(err).Error("are we there yet")

	time.Sleep(time.Second * 5)

	out, err := exec.Command("sh", "-c", fmt.Sprintf("journalctl 'TEST_ID=%s' -o json", uniqueSmallMessageID)).Output()
	fmt.Println(string(out))

	require.NoError(t, err)
	require.True(t, strings.Contains(string(out), "ERROR\" : \"something something dark side"))
	require.True(t, strings.Contains(string(out), "MESSAGE\" : \"are we there yet"))
	require.True(t, strings.Contains(string(out), "PRIORITY\" : \"3"))
}
