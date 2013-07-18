package main

import (
  "github.com/Mischanix/applog"
  "github.com/Mischanix/evconf"
  "github.com/Mischanix/wait"
  "os"
)

var ready = wait.NewFlag(false)

var logStdout = false

var config struct {
  DbUrl            string `json:"db_url"`
  DbName           string `json:"db_name"`
  MsgCollection    string `json:"msg_collection"`
  StatusCollection string `json:"status_collection"`
  HttpPort         int    `json:"http_port"`
}

func defaultConfig() {
  config.DbUrl = "localhost"
  config.DbName = "tl-dev"
  config.MsgCollection = "messages"
  config.StatusCollection = "statuses"
  config.HttpPort = 9002
}

func main() {
  // Log setup
  applog.Level = applog.DebugLevel
  if logStdout {
    applog.SetOutput(os.Stdout)
  } else {
    if logFile, err := os.OpenFile(
      "losapi.log",
      os.O_WRONLY|os.O_CREATE|os.O_APPEND,
      os.ModeAppend|0666,
    ); err != nil {
      applog.SetOutput(os.Stdout)
      applog.Error("Unable to open log file: %v", err)
    } else {
      applog.SetOutput(logFile)
    }
  }
  applog.Info("starting...")

  // Config setup
  conf := evconf.New("losapi.json", &config)
  conf.OnLoad(func() {
    ready.Set(true)
  })
  conf.StopWatching()
  defaultConfig()
  go func() {
    conf.Ready()
  }()

  ready.WaitFor(true)

  // Application
  go dbClient()
  httpServer()
}
