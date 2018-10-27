package paragliderdb

import (
	"net/http"
	"strings"

	igc "github.com/marni/goigc"
)

// RootHandler handles the root of the api
// And sends the requests where it should be
func RootHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	// if url /paragliding/"" redirect to api
	if parts[1] == "paragliding" && len(parts) == 3 && parts[2] == "" {
		http.Redirect(w, r, parts[0]+"/"+parts[1]+"/"+"api", 301)
	}
	if parts[1] != "paragliding" {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// from here we know the root has to be /paragliding/ or we would have got an 404 error
	// so here we will send it to the time handler
	if parts[2] == "api" {
		if len(parts) == 3 {
			if r.Method == "GET" {
				handlerTime(w, r)

			} else {
				// error status code wrong method from client side
				http.Error(w, http.StatusText(405), 405)
				return
			}
		}

	} else {
		// error some rubbish url
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// here we know if the parts[2] was api but not len(3) we will
	// check if the lenght was 4 and parts[3] is not rubbish

	if parts[3] == "track" {
		if len(parts) == 4 {
			// will respond with array of all tracks id
			if r.Method == "GET" {
				HandlerTrackArray(w, r)
			}
		}
	}

}

/*   This is the json format whenever u use get on the right path
{
	"H_date": <date from File Header, H-record>,
	"pilot": <pilot>,
	"glider": <glider>,
	"glider_id": <glider_id>,
	"track_length": <calculated total track length>
	"track_src_url": <the original URL used to upload the track, ie. the URL used with POST>
	}
*/

func trackIDToTrackMeta(url igc.Track, igc string) TrackMeta {
	return TrackMeta{
		url.Date,
		url.Pilot,
		url.GliderType,
		url.GliderID,
		TotalDistance(url),
		igc,
	}
}

// TotalDistance is the total distance of the track
// and returns an int representing the length
func TotalDistance(distance igc.Track) float64 {
	totalDistance := 0.0
	for i := 0; i < len(distance.Points)-1; i++ {
		totalDistance += distance.Points[i].Distance(distance.Points[i+1])
	}
	return totalDistance
}
