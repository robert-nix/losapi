package main

import (
  "github.com/Mischanix/applog"
  "net/http"
  "net/url"
  "strings"
)

// messages: all sort received_-1 : /messages?
// by user=?                                  user={user}
// by channel=?                               channel={channel}
// by command=?                               command={command}
// filter out/only commands                   is_command=[true|false]
// received range:                            start={time}
// end requires start to have effect          end={time}
// then filter regex/i                        match={string}
func handleMessages(w http.ResponseWriter, r *http.Request) {
  uri, err := url.ParseRequestURI(r.RequestURI)
  if err != nil {
    applog.Info("/messages: ParseRequestURI failed: %v", err)
    errJson(w, badRequest, 400)
    return
  }

  query := uri.Query()
  findQuery := dbM{}

  commandQueried := false
  canQuery := false
  if user := query.Get("user"); user != "" {
    canQuery = true
    findQuery["user"] = strings.ToLower(user)

    // user implies command can exist
    if command := query.Get("command"); command != "" {
      findQuery["command"] = strings.ToUpper(command)
      commandQueried = true
    } else if isCommand := query.Get("is_command"); isCommand != "" {
      commandQueried = true
      if isCommand == "true" {
        findQuery["command"] = dbM{"$ne": nil}
      } else {
        findQuery["command"] = nil
      }
    }
  }

  // channel implies command cannot exist
  if channel := query.Get("channel"); !commandQueried && channel != "" {
    canQuery = true
    findQuery["channel"] = strings.ToLower(channel)
  }
  if timeRange, _ := buildTimeRange(query); timeRange != nil {
    findQuery["received"] = timeRange
  }
  if match := query.Get("match"); match != "" {
    findQuery["message"] = dbM{
      "$regex":   regexEscape(match),
      "$options": "i",
    }
  }

  if !canQuery {
    applog.Info("/messages: not enough parameters")
    errJson(w, "/messages requires user or channel", 403)
    return
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
    errJson(w, serverError, 500)
    return
  }

  err = querySkipLimit(
    dbQuery, query.Get("offset"), query.Get("limit"),
  ).All(&result.Messages)
  if err != nil {
    applog.Error("/messages: query.All failed: %v", err)
    errJson(w, serverError, 500)
    return
  }

  writeJson(w, &result)
}
