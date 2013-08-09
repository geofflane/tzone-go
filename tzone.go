package main

import (
  "net/http"
  "fmt"
  "log"
  "time"
  "github.com/geofflane/tzone-go/util"
)

const TimeFormat = "2006-01-02T15:04:05.9999"

func main() {
  http.Handle("/convertCurrent", WithSecurityCheck{http.HandlerFunc(currentTime)})
  http.Handle("/convertTime", WithSecurityCheck{http.HandlerFunc(convertBetween)})

  defer Cleanup()
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func currentTime(w http.ResponseWriter, r *http.Request) {
  tz, err := tzForParam(r, "to")
  if nil != err {
    http.Error(w, "Must pass 'to' parameter with valid timezone", http.StatusBadRequest)
    return
  }

  writeJsonTime(w, time.Now().In(tz))
}

func convertBetween(w http.ResponseWriter, r *http.Request) {
  toTz, err := tzForParam(r, "to")
  if nil != err {
    http.Error(w, "Must pass 'to' parameter with valid timezone", http.StatusBadRequest)
    return
  }
  fromTz, err := tzForParam(r, "from")
  if nil != err {
    http.Error(w, "Must pass 'from' parameter with valid timezone", http.StatusBadRequest)
    return
  }

  timeParam := r.FormValue("time")
  time, err := time.ParseInLocation(TimeFormat, timeParam, fromTz)
  if nil != err {
    http.Error(w, fmt.Sprintf("Must pass 'time' parameter with format %s", TimeFormat), http.StatusBadRequest)
    return
  }

  writeJsonTime(w, time.In(toTz))
}

func tzForParam(r *http.Request, param string) (*time.Location, error) {
  val := r.FormValue(param)
  return time.LoadLocation(val)
}

func writeJsonTime(w http.ResponseWriter, t time.Time) {
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprint(w, util.JsonResponse{"time": t})
}

