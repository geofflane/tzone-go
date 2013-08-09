package main

import (
  "github.com/geofflane/tzone-go/data"
  "net/http"
	"net/http/httptest"
  "testing"
  "errors"
  "github.com/stretchr/testify/assert"
  "time"
)


type TestUserDb struct {
  Users map[string]data.User
  Usage int
  Done chan bool
}
func (db *TestUserDb) Close() error {
  // No op
  return nil
}

func (db *TestUserDb) Authenticate(token string) (data.User, error) {
  u, ok:= db.Users[token]
  if ! ok {
    return u, errors.New("No such user")
  }
  return u, nil
}

func (db *TestUserDb) RecordUsage(u data.User) {
  db.Usage += 1
  db.Done <- true
}



func TestPassingGoodTokenCallsBaseMethod(t *testing.T) {
  called := false
  doneChan := make(chan bool)

  w := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "http://example.com?token=test", nil)
  userDb := &TestUserDb{map[string]data.User{"test": data.User{}}, 0, doneChan}

  secCheck := WithSecurityCheck{userDb, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })}
  secCheck.ServeHTTP(w, req)

  waitForDone(doneChan)

  assert.True(t, called, "handler was not called")
  assert.Equal(t, 1, userDb.Usage)
}

func TestPassingBadTokenDoesNotCallBaseMethod(t *testing.T) {
  called := false

  w := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "http://example.com?token=test", nil)
  userDb := &TestUserDb{map[string]data.User{"different": data.User{}}, 0, nil}

  secCheck := WithSecurityCheck{userDb, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })}
  secCheck.ServeHTTP(w, req)

  assert.False(t, called, "handler was not supposed to be called")
  assert.Equal(t, 0, userDb.Usage)
}

func waitForDone(c chan bool) {
  // record usage is async, so need to wait for it to be set
  // supposedly this would handle a timeout

  select {
  case <-time.After(10 * time.Nanosecond):
    println("Timeout")
  case _ = <-c:
    println("Done")
  }
}
