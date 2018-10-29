package paragliderdb

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
	if len(parts) == 3 {
		if parts[2] == "api" {
			if r.Method == "GET" {
				handlerTime(w, r)

			} else {
				// error status code wrong method from client side
				http.Error(w, http.StatusText(405), 405)
				return
			}
		} else {
			// error some rubbish url
			http.Error(w, http.StatusText(404), 404)
			return
		}

	}

	// here we know if the parts[2] was api but not len(3) we will
	// check if the lenght was 4 and parts[3] is not rubbish

	if len(parts) == 4 {
		if parts[3] == "track" && parts[2] == "api" {
			// will respond with array of all tracks id
			if r.Method == "GET" {
				HandlerTrackArray(w, r)
			} else if r.Method == "POST" {
				HandlerPostURL(w, r)
			} else {
				http.Error(w, http.StatusText(405), 405)
			}
		} else if parts[3] == "ticker" && parts[2] == "api" {
			//handles api/ticker and we know parts[1] are paragliding
			if r.Method == "GET" {
				//		HandlesTicker(w, r)
			} else {
				http.Error(w, http.StatusText(405), 405)
			}
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

	}

	// here we know that when len is 5 the url should be
	// api/track/<id> and this handles that request and
	// checks for valid ids
	if len(parts) == 5 {
		if parts[3] == "track" && parts[2] == "api" {
			if r.Method == "GET" {
				if CheckForValidID(parts[4]) {
					HandleOneTrackMeta(w, r, parts[4])
				} else {
					http.Error(w, "Rubbish id, Please enter a valid one", http.StatusBadRequest)
				}
			} else {
				http.Error(w, http.StatusText(405), 405)
			}
		} else if parts[3] == "ticker" && parts[2] == "api" && parts[4] == "latest" {
			// handles latest timestamp
			HandleLatestTimestamp(w, r)

		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

	}

	// here the length is 6 and the url should be api/track/id/field
	if len(parts) == 6 {
		if parts[3] == "track" && parts[2] == "api" {
			if r.Method == "GET" {
				if CheckForValidID(parts[4]) {
					if CheckForValidField(parts[5]) {
						HandlesField(w, r, parts[5], parts[4])
					} else {
						http.Error(w, "Rubbish field, Please enter a valid one", http.StatusBadRequest)
					}
				} else {
					http.Error(w, "Rubbish id, Please enter a valid one", http.StatusBadRequest)
					return
				}
			} else {
				http.Error(w, http.StatusText(405), 405)
				return
			}
		} else {
			http.Error(w, "rubbish url, Please enter a valid path", http.StatusBadRequest)
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

func trackIDToTrackMeta(igcu string) (TrackMeta, error) {

	url, err := igc.ParseLocation(igcu)

	return TrackMeta{
		url.Date,
		url.Pilot,
		url.GliderType,
		url.GliderID,
		TotalDistance(url),
		igcu,
	}, err
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

// NewUniqueParagliderID checks if its any error comming from
// the random string generator, if not it returns the string
func NewUniqueParagliderID() TrackID {
	id, err := GenerateRandomString(MAXLENGTHID)
	if err != nil {
		fmt.Println("Something wrong happend with the random generator", err)
		return TrackID{http.StatusText(500)}
	}
	return TrackID{id}
}

// GenerateRandomBytes with help of make we will initialse a
// a random byte and we will return it if b has the lenght
// equals to our n integer number
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString generates a random string from the
// GenerateRandomBytes and get returned a byte that gets
// encoded into a string, that we will return as a id to the
// paraglider post request
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	if err != nil {
		fmt.Println("Something happend with the random generator")
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}

// CheckForValidID checks if a id is in the db
func CheckForValidID(id string) bool {

	all := GlobalDB.GetAllID()

	for _, data := range all {
		if data.ID == id {
			return true
		}
	}
	return false
}

// CheckForValidField check if the fields is one of the expected one
// or default will set false
func CheckForValidField(field string) bool {
	switch field {
	case "H_date":
		return true

	case "pilot":
		return true

	case "glider":
		return true

	case "glider_id":
		return true

	case "track_length":
		return true

	case "track_src_url":
		return true

	default:
		return false
	}
}
