package paragliderdb

import (
	"net/http"
	"strings"
)

// RootHandler handles the root of the api
// And sends the requests where it should be
func RootHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if parts[1] != "paraglider" || len(parts) <= 2 {
		http.Error(w, http.StatusText(404), 404)
		return
	} else {
		if parts[2] == "api" && len(parts) == 3 {
			handlerTime(w, r)
		}
	}
}
