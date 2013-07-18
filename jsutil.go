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

func writeJson(w http.ResponseWriter, result interface{}) {
  if err := json.NewEncoder(w).Encode(result); err != nil {
    applog.Error("writeJson: json Encode failed: %v", err)
    http.Error(w, serverError, 500)
  }
}
