package main

import (
  "github.com/Mischanix/applog"
  "net/http"
  "net/url"
)

// statuses: all sort timestamp_-1 :
// by channel=?                     /channel/:channel?
// then by timestamp range                            start={time}
//  requires start                                    end={time}
func handleStatusesChannel(w http.ResponseWriter, r *http.Request) {
  uri, err := url.ParseRequestURI(r.RequestURI)
  if err != nil {
    applog.Info("/channel/: ParseRequestURI failed: %v", err)
    http.Error(w, badRequest, 400)
    return
  }
  channel := uri.Path[len(channelPath):]

  query := uri.Query()
  findQuery := dbM{"channel": channel}
  if timeRange, _ := buildTimeRange(query); timeRange != nil {
    findQuery["timestamp"] = timeRange
  }

  applog.Debug("/channel/: query built: %v", findQuery)

  writeStatuses(w, findQuery, query)
}

func writeStatuses(w http.ResponseWriter, findQuery dbM, query url.Values) {
  var result struct {
    Count    int         `json:"count"`
    Statuses []statusDoc `json:"statuses"`
  }
  dbQuery := db.statusColl.Find(findQuery).Sort("-timestamp")

  var err error
  result.Count, err = dbQuery.Count()
  if err != nil {
    applog.Error("writeStatuses: query.Count failed: %v", err)
    http.Error(w, serverError, 500)
    return
  }

  err = querySkipLimit(
    dbQuery, query.Get("offset"), query.Get("limit"),
  ).All(&result.Statuses)
  if err != nil {
    applog.Error("writeStatuses: query.All failed: %v", err)
    http.Error(w, serverError, 500)
    return
  }

  writeJson(w, &result)
}
