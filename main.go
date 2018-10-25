package main

import (
	"fmt"
	"goprojects/paraglider/paragliderdb"
	"log"
	"net/http"
	"os"
	"time"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func serverStart() {
	paragliderdb.StartTime = time.Now()
}

func main() {

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	serverStart()
	http.HandleFunc("/paragliding/", paragliderdb.RootHandler)
	http.ListenAndServe(addr, nil)

}
