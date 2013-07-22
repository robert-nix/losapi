package main

import (
  "encoding/json"
  "github.com/Mischanix/applog"
  "labix.org/v2/mgo/bson"
  "net/http"
  "strconv"
  "time"
)

type jsDate struct {
  t time.Time
}

func (d *jsDate) SetBSON(raw bson.Raw) error {
  return raw.Unmarshal(&d.t)
}

func (d *jsDate) MarshalJSON() (result []byte, err error) {
  n := d.t.UnixNano() / 1e6
  return strconv.AppendInt(nil, n, 10), nil
}

func jsDateToTime(date string) time.Time {
  ms, _ := strconv.ParseInt(date, 10, 64)
  return time.Unix(0, ms*1e6)
}

// regexEscape escapes a string's regex metachars. Deliberately copied from
// mono's System.Text.RegularExpressions implementation
func regexEscape(regex string) (result string) {
  for _, c := range regex {
    switch c {
    case '\\', '*', '+', '?', '|', '{', '[',
      '(', ')', '^', '$', '.', '#', ' ':
      result += string('\\' + c)

    case '\t':
      result += "\\t"
    case '\n':
      result += "\\n"
    case '\r':
      result += "\\r"
    case '\f':
      result += "\\f"

    default:
      result += string(c)
    }
  }
  return result
}

func responseHeaders(w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
}

func writeJson(w http.ResponseWriter, result interface{}) {
  responseHeaders(w)
  w.WriteHeader(200)
  if err := json.NewEncoder(w).Encode(result); err != nil {
    applog.Error("writeJson: json Encode failed: %v", err)
  }
}

func errJson(w http.ResponseWriter, error string, code int) {
  responseHeaders(w)
  w.WriteHeader(code)
  if err := json.NewEncoder(w).Encode(struct {
    Status string `json:"status"`
    Reason string `json:"reason"`
  }{"nok", error}); err != nil {
    applog.Error("writeJson: json Encode failed: %v", err)
  }
}
