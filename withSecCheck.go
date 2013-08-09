package main

import (
  "net/http"
  "log"
  "github.com/geofflane/tzone-go/data"
)


var userDb data.UserDb
var usageChan chan data.User

func init() {
  usageChan = make(chan data.User)

  var err error
  userDb, err = data.NewUserDb()
  if nil != err {
    log.Fatal("Couldn't connect to userDb: %s", err)
    return
  }
  go LogUsage(usageChan)
}

type WithSecurityCheck struct {
  f http.Handler
}

func (sc WithSecurityCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  token := r.FormValue("token")
  u, err := userDb.Authenticate(token)
  if nil != err {
    http.Error(w, "Need authentication", http.StatusForbidden)
    return
  }

  usageChan <- u
  sc.f.ServeHTTP(w, r)
}

func Cleanup() {
  userDb.Close()
  close(usageChan)
}

func LogUsage(c chan data.User) {
  for {
    u := <- c
    userDb.RecordUsage(&u)
  }
}


