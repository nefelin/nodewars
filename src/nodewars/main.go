package main

import (
	"log"
	"net/http"
	"nwmodel"
)

func main() {

	log.Println("Starting " + nwmodel.VersionTag + " server...")

	// Start Webserver
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", nwmodel.HandleConnections)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
