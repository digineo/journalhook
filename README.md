# `Systemd's` Journal hook for `logrus` <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

[![Build Status](https://travis-ci.org/ssgreg/journalhook.svg?branch=master)](https://travis-ci.org/ssgreg/journalhook)
[![Go Report Status](https://goreportcard.com/badge/github.com/ssgreg/journalhook)](https://goreportcard.com/report/github.com/ssgreg/journalhook)

## Installation

Install the package with:

```shell
go get -u github.com/ssgreg/journalhook
```

## Usage

```go
package main

import (
    "os"

    "github.com/sirupsen/logrus"
    "github.com/ssgreg/journalhook"

fun  main() {
    log := logrus.New()
    hook, err := journalhook.NewJournalHook()
    if err == nil {
        log.Hooks.Add(hook)
    }

    logEntry := log.WithFields(logrus.Fields{
        "n_goroutine": runtime.NumGoroutine(),
        "executable":  os.Args[0],
        "trace":       runtime.ReadTrace(),
    })

    logEntry.Info("Hello World!")

    // Make sure that journal connection will be successfully closed
    // and no message will be lost.
    logrus.Exit(0)
}
```

This is how it will look like:

```json
{
  "__CURSOR": "s=f81e8eb7fd0941b089528d889c929c1f;i=f1;b=40582011948a4b1998bf5ca928517a0f;m=2345f2fc69;t=557cb96d25735;x=e26ade807f4bab79",
  "__REALTIME_TIMESTAMP": "1503906803898165",
  "__MONOTONIC_TIMESTAMP": "151497407593",
  "_BOOT_ID": "40582011948a4b1998bf5ca928517a0f",
  "PRIORITY": "6",
  "_UID": "0",
  "_GID": "0",
  "_CAP_EFFECTIVE": "a80425fb",
  "_MACHINE_ID": "78b67a34f030403aa4dc97056d5efced",
  "_HOSTNAME": "64620c2b0d13",
  "_TRANSPORT": "journal",
  "MESSAGE": "Hello World!",
  "N_GOROUTINE": "2",
  "TRACE": [ 103, 111, 32, 49, 46, 56, 32, 116, 114, 97, 99, 101, 0, 0, 0, 0 ],
  "_COMM": "usage",
  "EXECUTABLE": "./usage",
  "_PID": "7533",
  "_SOURCE_REALTIME_TIMESTAMP": "1503906803896890"
}
```