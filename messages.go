package main

import (
  "github.com/Mischanix/applog"
  "net/http"
  "net/url"
  "time"
)

// messages: all sort received_-1 : /messages?
// by user=?                                  user={user}
// by channel=?                               channel={channel}
// received range:                            start={time}
// end requires start to have effect          end={time}
// then filter regex/i                        match={string}
func handleMessages(w http.ResponseWriter, r *http.Request) {
  applog.Info("/messages: uri=%q", r.RequestURI)

  uri, err := url.ParseRequestURI(r.RequestURI)
  if err != nil {
    applog.Info("/messages: ParseRequestURI failed: %v", err)
    http.Error(w, badRequest, 400)
    return
  }

  query := uri.Query()
  findQuery := dbM{}
  // whether matching is allowed -- i.e. is the result set sufficiently filtered
  // such that this won't murder the server
  canMatch := false
  if user := query.Get("user"); user != "" {
    findQuery["user"] = user
    canMatch = true
  }
  if channel := query.Get("channel"); channel != "" {
    findQuery["channel"] = channel
    canMatch = true
  }
  if timeRange, duration := buildTimeRange(query); timeRange != nil {
    findQuery["received"] = timeRange
    if duration > 12*time.Hour {
      canMatch = true
    }
  }
  if match := query.Get("match"); canMatch && match != "" {
    findQuery["message"] = dbM{
      "$regex":   regexEscape(match),
      "$options": "i",
    }
  }

  applog.Debug("/messages: query built: %v", findQuery)

  var result struct {
    Count    int      `json:"count"`
    Messages []msgDoc `json:"messages"`
  }
  dbQuery := db.msgColl.Find(findQuery).Sort("-received")

  result.Count, err = dbQuery.Count()
  if err != nil {
    applog.Error("/messages: query.Count failed: %v", err)
    http.Error(w, serverError, 500)
    return
  }

  err = querySkipLimit(
    dbQuery, query.Get("offset"), query.Get("limit"),
  ).All(&result.Messages)
  if err != nil {
    applog.Error("/messages: query.All failed: %v", err)
    http.Error(w, serverError, 500)
    return
  }

  writeJson(w, &result)
}
