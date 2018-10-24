package paragliderdb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandlerRoot_AllURLRespondWith404(t *testing.T) {

	testvalue := []string{"/paraglider", "/hhdhfhshasaj", "/rubbish", "/paraglider", "/PAraglider"}

	ts := httptest.NewServer(http.HandlerFunc(RootHandler))
	defer ts.Close()

	for i := 0; i < len(testvalue); i++ {
		fmt.Println(ts.URL)
		resp, err := http.Get(ts.URL + testvalue[i])

		if err != nil {
			t.Error("Something is wrong with the request")
			return
		}
		if resp.StatusCode != 404 {
			t.Errorf("For get %s, expected StatusCode %d, received %d", testvalue[i],
				404, resp.StatusCode)
			return
		}

	}

}

func Test_HandlerTime_api(t *testing.T) {
	testvalue := []string{"/api", "/hhdhfhshasaj", "/rubbish", "/paraglider", "/PAraglider"}

	ts := httptest.NewServer(http.HandlerFunc(RootHandler))
	defer ts.Close()

	for i := 0; i < len(testvalue); i++ {
		fmt.Println(ts.URL + "/paraglider" + testvalue[i])
		resp, err := http.Get(ts.URL + testvalue[i])

		if err != nil {
			t.Error("Something is wrong with the request")
			return
		}
		if testvalue[i] != "api" {
			if resp.StatusCode != 404 {
				t.Errorf("For get %s, expected StatusCode %d, received %d", testvalue[i],
					404, resp.StatusCode)
				return
			}
		} else {
			if resp.StatusCode != 200 {
				t.Errorf("For get /paraglider%s, expected StatusCode %d, received %d", testvalue[i],
					200, resp.StatusCode)
			}
		}

	}

}
