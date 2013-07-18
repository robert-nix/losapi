package main

import (
  "github.com/Mischanix/applog"
  "github.com/Mischanix/wait"
  "labix.org/v2/mgo"
  "strconv"
  // "time"
)

var db struct {
  ready      *wait.Flag
  session    *mgo.Session
  database   *mgo.Database
  msgColl    *mgo.Collection
  statusColl *mgo.Collection
}

type dbM map[string]interface{}

type msgDoc struct {
  User     string `json:"user"`
  Channel  string `json:"channel"`
  Received jsDate `json:"received"`
  Message  string `json:"message"`
}

type statusDoc struct {
  Channel   string   `json:"channel"`
  Timestamp jsDate   `json:"timestamp"`
  Status    string   `json:"status"`
  Users     []string `json:"users"`
  Viewers   int      `json:"viewers"`
}

func dbClient() {
  var err error
  db.ready = wait.NewFlag(false)
  if db.session, err = mgo.Dial(config.DbUrl); err != nil {
    applog.Error("mgo.Dial failure: %v", err)
  }
  db.database = db.session.DB(config.DbName)
  db.msgColl = db.database.C(config.MsgCollection)
  db.statusColl = db.database.C(config.StatusCollection)

  db.ready.Set(true)
}

func querySkipLimit(dbQuery *mgo.Query, skip string, limit string) *mgo.Query {
  skipCount, _ := strconv.Atoi(skip)
  limitCount, _ := strconv.Atoi(limit)
  if limitCount > 500 {
    limitCount = 500
  }
  if limitCount <= 0 {
    limitCount = 100
  }

  return dbQuery.Skip(skipCount).Limit(limitCount)
}
