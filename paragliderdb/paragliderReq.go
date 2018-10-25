package paragliderdb

import (
	"net/http"
	"strings"
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
	if parts[2] == "api" && len(parts) == 3 {
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
