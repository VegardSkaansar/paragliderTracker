package main

import (
	"fmt"
	"goprojects/paraglider/paragliderdb"
	"log"
	"net/http"
	"os"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/paraglider", paragliderdb.RootHandler)
	http.ListenAndServe(addr, nil)
}
