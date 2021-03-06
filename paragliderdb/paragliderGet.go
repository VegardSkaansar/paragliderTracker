package paragliderdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ------------------------------------------------------------------------
// -                 -- Handles api --                                    -
//-------------------------------------------------------------------------

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

// HandleOneTrackMeta handles all the api/track/<id> requests
func HandleOneTrackMeta(w http.ResponseWriter, r *http.Request, id string) {
	single := GlobalDB.GetTrackMeta(id)
	json.NewEncoder(w).Encode(&single)
}

// HandlesField handles all the api/track/<id>/<field> requests
func HandlesField(w http.ResponseWriter, r *http.Request, field string, id string) {
	single := GlobalDB.GetTrackMeta(id)

	// maybe a better way to solve this, couldnt find one since i kept
	// getting plain/text what every i tryed, in this switch we make
	// a temporary struct for the field we are going give back as json
	switch field {
	case "H_date":
		json.NewEncoder(w).Encode(single.Date)
	case "pilot":
		json.NewEncoder(w).Encode(single.Pilot)

	case "glider":
		json.NewEncoder(w).Encode(single.Glider)

	case "glider_id":
		json.NewEncoder(w).Encode(single.GliderID)

	case "track_length":
		json.NewEncoder(w).Encode(single.TrackLength)

	case "track_src_url":
		json.NewEncoder(w).Encode(single.TrackSrcURL)

	default:
		// if this happends something is wrong
		panic("shouldnt happend with the fields above")
	}

}

// ------------------------------------------------------------------------
// -                 -- Handles ticker --                                 -
//-------------------------------------------------------------------------

// HandleLatestTimestamp check the db for the latest added track
func HandleLatestTimestamp(w http.ResponseWriter, r *http.Request) {
	one := GlobalDB.GetLatestObjectID()
	log.Println(one)
	json.NewEncoder(w).Encode(one.Time().Unix())
}

//HandlesTicker handles a ticker struct and format it to json for a get request
func HandlesTicker(w http.ResponseWriter, r *http.Request, startprocess int64) {
	process := (time.Now().UnixNano() / int64(time.Millisecond)) - startprocess
	start, stop := GlobalDB.GetTicker()
	var tra []string

	if len(GlobalDB.GetAllID()) <= TICKERIDLENGTH {
		for i := 0; i < len(GlobalDB.GetAllID()); i++ {
			tra = append(tra, GlobalDB.GetAllID()[i].ID)
		}
	} else {
		for i := 0; i <= TICKERIDLENGTH; i++ {
			tra = append(tra, GlobalDB.GetAllID()[i].ID)
		}
	}

	newTicker := Ticker{start.Time().Unix(), stop.Time().Unix(), GlobalDB.GetLatestObjectID().Time().Unix(), tra, process}

	json.NewEncoder(w).Encode(&newTicker)
}

// HandleTimestamp handles timestamp from client
func HandleTimestamp(w http.ResponseWriter, r *http.Request, i int64, startprocess int64) {
	process := (time.Now().UnixNano() / int64(time.Millisecond)) - startprocess
	start, stop := GlobalDB.GetTicker()
	strID := GlobalDB.RequestTimestamp(i)
	var tra []string

	if len(strID) <= TICKERIDLENGTH {
		for i := 0; i < len(strID); i++ {
			tra = append(tra, strID[i])
		}
	} else {
		for i := 0; i <= TICKERIDLENGTH; i++ {
			tra = append(tra, strID[i])
		}
	}

	newTicker := Ticker{start.Time().Unix(), stop.Time().Unix(), GlobalDB.GetLatestObjectID().Time().Unix(), tra, process}

	json.NewEncoder(w).Encode(&newTicker)

}

// WebhookGetHandler handles a webhook get id
func WebhookGetHandler(w http.ResponseWriter, r *http.Request, id string) {
	str, ok := GlobalDB.GetWebhook(id)

	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	json.NewEncoder(w).Encode(str)
}

// WebhookDeleteHandler handles deleting a webhook
func WebhookDeleteHandler(w http.ResponseWriter, r *http.Request, id string) {

	str, ok := GlobalDB.GetWebhook(id)

	if ok == false {
		http.Error(w, "couldnt find that id", 500)
		return
	}

	okk := GlobalDB.DeleteWebhook(id)

	if !okk {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	jsonResponse, err := json.Marshal(str)
	if err != nil {
		http.Error(w, "cant transform the new id to bytes", 500)
	}
	// we write a json response with the with info
	w.Write(jsonResponse)

}

// HandlerAdmin delete all records
func HandlerAdmin(w http.ResponseWriter, r *http.Request) {
	GlobalDB.AdmimDeleteAll()
}
