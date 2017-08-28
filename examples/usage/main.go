package main

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/ssgreg/journalhook"
)

func main() {
	log := logrus.New()
	hook, err := journalhook.NewJournalHook()
	if err == nil {
		log.Hooks.Add(hook)
	}

	log.WithFields(logrus.Fields{
		"n_goroutine": runtime.NumGoroutine(),
		"executable":  os.Args[0],
		"trace":       runtime.ReadTrace(),
	}).Info("Hello World!")

	// Make sure that journal connection will be successfully closed
	// and no log entries will be lost.
	logrus.Exit(0)
}
