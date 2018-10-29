package paragliderdb

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// HandlerPostURL handles a post request from api/track
func HandlerPostURL(w http.ResponseWriter, r *http.Request) {

	// url is a temporary struct, just getting the url
	type url struct {
		URL string `json:"url"`
	}
	var i url

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if r.Body == nil {
		http.Error(w, "POST request must have a JSON body", http.StatusBadRequest)
		return
	}

	igcTrack, err := trackIDToTrackMeta(i.URL)
	if err != nil {
		http.Error(w, "Something is wrong with url", http.StatusBadRequest)
	}
	ok := GlobalDB.CheckIfURLIsAlreadyTracked(i.URL)
	log.Println(ok)
	if ok {
		u := NewUniqueParagliderID()
		GlobalDB.AddURL(igcTrack, u)

		// with help of marshal we transform this to a json response
		jsonResponse, err := json.Marshal(u)
		if err != nil {
			http.Error(w, "cant transform the new id to bytes", 500)
		}
		// we write a json response with the id
		w.Write(jsonResponse)

	} else {
		w.WriteHeader(http.StatusAlreadyReported)
		io.WriteString(w, "This Track are already requested and stored in database")
	}
}

// WebhookPost handles a webhook post
func WebhookPost(w http.ResponseWriter, r *http.Request) {

	var webreq Webhookinfo

	err := json.NewDecoder(r.Body).Decode(&webreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		http.Error(w, "POST request must have a JSON body", http.StatusBadRequest)
		return
	}

	err = GlobalDB.AddWebhook(webreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(webreq)
	if err != nil {
		http.Error(w, "cant transform the new id to bytes", 500)
	}
	// we write a json response with the with info
	w.Write(jsonResponse)

}
