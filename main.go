package main

import (
	"fmt"
	"goprojects/paraglider/paragliderdb"
	"log"
	"math/rand"
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
	// this makes a random seed for our program
	// this will accutally help us to get a random
	// number
	rand.Seed(time.Now().UnixNano())
}

func main() {

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	serverStart()

	// Initialising the global db
	paragliderdb.GlobalDB = &paragliderdb.MongoDB{
		"mongodb://Vegard:Mira1234@ds143893.mlab.com:43893/paragliderdb",
		"paragliderDB",
		"tracks",
	}

	paragliderdb.GlobalDB.Init()

	http.HandleFunc("/paragliding/", paragliderdb.RootHandler)
	http.ListenAndServe(addr, nil)

}
