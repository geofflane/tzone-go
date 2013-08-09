package main

import (
  "net/http"
  "github.com/geofflane/tzone-go/data"
)

type WithSecurityCheck struct {
  userDb data.UserDb
  f http.Handler
}

func (sc WithSecurityCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  token := r.FormValue("token")
  u, err := sc.userDb.Authenticate(token)
  if nil != err {
    http.Error(w, "Need authentication", http.StatusForbidden)
    return
  }

  go sc.RecordUsage(u)
  sc.f.ServeHTTP(w, r)
}

func (sc WithSecurityCheck) RecordUsage(u data.User) {
  sc.userDb.RecordUsage(u)
}

