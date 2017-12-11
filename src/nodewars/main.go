package main

import (
	"log"
	"net/http"
	"nwmodel"
	"os"
)

func main() {
	certfile := os.Getenv("CERTFILE")
	keyfile := os.Getenv("KEYFILE")
	prod := os.Getenv("PROD")

	log.Println("Starting " + nwmodel.VersionTag + " server...")

	// Start Webserver
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.FileServer(http.Dir("public")).ServeHTTP)
	mux.HandleFunc("/ws", nwmodel.HandleConnections)

	if prod == "" { // aka env var not set
		log.Fatal(
			http.ListenAndServe(":8080", mux))
	}

	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	log.Fatal(
		http.ListenAndServeTLS(":443", certfile, keyfile, mux))
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	log.Printf("redirect to %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

/*
func index(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		log.Printf("404: %s", req.URL.String())
		http.NotFound(w, req)
		return
	}
}
*/
