package main

import (
	"github.com/geofflane/tzone-go/data"
	"net/http"
)

// Implements http.Handler interface
// Wrapper to add a Security check to another http.Handler
type WithSecurityCheck struct {
	userDb data.UserDb
	f      http.Handler
}

func (sc WithSecurityCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	u, err := sc.userDb.Authenticate(token)
	if nil != err {
		http.Error(w, "Need authentication", http.StatusForbidden)
		return
	}

	go sc.userDb.RecordUsage(u)
	sc.f.ServeHTTP(w, r)
}
