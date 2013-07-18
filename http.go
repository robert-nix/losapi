package main

import (
  "fmt"
  "github.com/Mischanix/applog"
  "net/http"
)

var badRequest = `{"status":"nok", "reason":"Bad request"}`
var serverError = `{"status":"nok", "reason":"Internal server error"}`

func httpServer() {
  http.HandleFunc("/messages", handleMessages)
  dialString := fmt.Sprintf(":%d", config.HttpPort)
  if err := http.ListenAndServe(dialString, nil); err != nil {
    applog.Error("httpServer: ListenAndServe error: %v", err)
  }
}
