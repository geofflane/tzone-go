package main

import (
	"fmt"
	"github.com/geofflane/tzone-go/data"
	"github.com/geofflane/tzone-go/util"
	"log"
	"net/http"
	"time"
)

const TimeFormat = "2006-01-02T15:04:05.9999"

var userDb data.UserDb

// Run the program
func main() {
	var err error
	userDb, err = data.NewUserDb()
	if nil != err {
		log.Fatal("Couldn't connect to userDb ", err)
		return
	}
	defer userDb.Close()

	http.Handle("/convertCurrent", WithSecurityCheck{userDb, http.HandlerFunc(currentTime)})
	http.Handle("/convertTime", WithSecurityCheck{userDb, http.HandlerFunc(convertBetween)})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// http.Handler to convert the current time to the time in a different timezone.
// Expects an HTTP query paramter of "to" equal to an Olson DB timezone name
func currentTime(w http.ResponseWriter, r *http.Request) {
	tz, err := tzForParam(r, "to")
	if nil != err {
		http.Error(w, "Must pass 'to' parameter with valid timezone", http.StatusBadRequest)
		return
	}

	writeJsonTime(w, time.Now().In(tz))
}

// http.Handler to convert the a given time from one timezone to the time in a different timezone
// Expects an HTTP query paramter of "to" equal to an Olson DB timezone name
// Expects an HTTP query paramter of "from" equal to an Olson DB timezone name
// Expects an HTTP query paramter of "time" equal to a time in ISO format
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

// Get a timezone from a specified request paramter
func tzForParam(r *http.Request, param string) (*time.Location, error) {
	val := r.FormValue(param)
	return time.LoadLocation(val)
}

// Write a time as Json
func writeJsonTime(w http.ResponseWriter, t time.Time) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, util.JsonResponse{"time": t})
}
