package main

import (
  "fmt"
  "github.com/Mischanix/applog"
  "net/http"
  "time"
)

var badRequest = `{"status":"nok", "reason":"Bad request"}`
var serverError = `{"status":"nok", "reason":"Internal server error"}`

var channelPath = "/channel/"
var messagesPath = "/messages"

func log(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    applog.Info("http: %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)
    fn(w, r)
    applog.Info("http: %s %s took %v to process",
      r.Method, r.RequestURI,
      time.Now().Sub(start),
    )
  }
}

func httpServer() {
  http.HandleFunc(messagesPath, log(handleMessages))
  http.HandleFunc(channelPath, log(handleStatusesChannel))

  dialString := fmt.Sprintf(":%d", config.HttpPort)
  if err := http.ListenAndServe(dialString, nil); err != nil {
    applog.Error("httpServer: ListenAndServe error: %v", err)
  }
}
