package main

import (
	"errors"
	"github.com/geofflane/tzone-go/data"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestUserDb struct {
	Users map[string]data.User
	Usage int
	Done  chan bool
}

func (db *TestUserDb) Close() error {
	// No op
	return nil
}

func (db *TestUserDb) Authenticate(token string) (data.User, error) {
	u, ok := db.Users[token]
	if !ok {
		return u, errors.New("No such user")
	}
	return u, nil
}

func (db *TestUserDb) RecordUsage(u data.User) {
	db.Usage += 1
	db.Done <- true
}

func buildHttpParams() (w http.ResponseWriter, r *http.Request) {
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "http://example.com?token=test", nil)
	return
}

func buildHandlerFunc(called *bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
	})
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

func TestPassingGoodTokenCallsBaseMethod(t *testing.T) {
	var called bool
	doneChan := make(chan bool)
	defer close(doneChan)

	w, req := buildHttpParams()
	userDb := &TestUserDb{map[string]data.User{"test": data.User{}}, 0, doneChan}

	secCheck := WithSecurityCheck{userDb, buildHandlerFunc(&called)}
	secCheck.ServeHTTP(w, req)

	waitForDone(doneChan)

	assert.True(t, called, "handler was not called")
	assert.Equal(t, 1, userDb.Usage)
}

func TestPassingBadTokenDoesNotCallBaseMethod(t *testing.T) {
	var called bool

	w, req := buildHttpParams()
	userDb := &TestUserDb{map[string]data.User{"different": data.User{}}, 0, nil}

	secCheck := WithSecurityCheck{userDb, buildHandlerFunc(&called)}
	secCheck.ServeHTTP(w, req)

	assert.False(t, called, "handler was not supposed to be called")
	assert.Equal(t, 0, userDb.Usage)
}
