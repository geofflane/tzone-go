package main

import (
  "github.com/geofflane/tzone-go/data"
  "net/http"
	"net/http/httptest"
  "testing"
  "errors"
)


type TestUserDb struct {
  users map[string]data.User
  usage int
}
func (db TestUserDb) Close() error {
  // No op
  return nil
}

func (db TestUserDb) Authenticate(token string) (data.User, error) {
  u, ok:= db.users[token]
  if ! ok {
    return u, errors.New("No such user")
  }
  return u, nil
}

func (db TestUserDb) RecordUsage(u data.User) {
  db.usage += 1
}




func TestPassingGoodTokenCallsBaseMethod(t *testing.T) {
  called := false

  w := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "http://example.com?token=test", nil)
  userDb := TestUserDb{map[string]data.User{"test": data.User{}}, 0}

  secCheck := WithSecurityCheck{userDb, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })}
  secCheck.ServeHTTP(w, req)

  if ! called {
    t.Error("handler was not called")
  }
}

func TestPassingBadTokenDoesNotCallBaseMethod(t *testing.T) {
  called := false

  w := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "http://example.com?token=test", nil)
  userDb := TestUserDb{map[string]data.User{"different": data.User{}}, 0}

  secCheck := WithSecurityCheck{userDb, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })}
  secCheck.ServeHTTP(w, req)

  if called {
    t.Error("handler was not supposed to be called")
  }
}

