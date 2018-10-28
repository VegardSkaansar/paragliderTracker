package paragliderdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func handlerTime(w http.ResponseWriter, r *http.Request) {
	Timediff := getServerTime()

	service := ServiceTime{
		Timediff,
		"Service for Paragliding tracks.",
		"v1",
	}
	json.NewEncoder(w).Encode(&service)
}

func getServerTime() string {
	// Time conversions
	// 60 Seconds (minute)
	// 3600 seconds (hour)
	// 86400 seconds (day)
	// Aproximate 2592000 seconds (months)
	// 31536000 seconds (year)
	var years int
	var months int
	var days int
	var hours int
	var minute int
	var seconds int
	var rest int

	elapsed := time.Now().Sub(StartTime)
	sec := elapsed.Seconds()

	years = int(sec) / 31536000
	rest = (int(sec) - (years * 31536000))
	months = rest / 2592000
	days = (rest - (months * 2592000)) / 86400
	rest -= (months * 2592000)
	hours = (rest - (days * 86400)) / 3600
	rest -= (days * 86400)
	minute = (rest - (hours * 3600)) / 60
	rest -= (hours * 3600)
	seconds = (rest - (minute * 60))

	return fmt.Sprintf("P%dY%dM%dD%dH%dM%dS", years, months, days, hours, minute, seconds)
}

//HandlerTrackArray handles the database request and gives us
// all the ids of the tracks stored in a array
func HandlerTrackArray(w http.ResponseWriter, r *http.Request) {
	all := GlobalDB.GetAllID()
	json.NewEncoder(w).Encode(&all)

}
