package paragliderdb

import (
	"encoding/json"
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
	GlobalDB.AddURL(igcTrack)
}
