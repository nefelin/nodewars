package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"protocol"
)

const localhost = "http://localhost:8080/"
const localport = ":8080"

var host string
var port string

func main() {

	certfile := os.Getenv("CERTFILE")
	keyfile := os.Getenv("KEYFILE")
	prod := os.Getenv("PROD")
	if prod == "" {
		host = localhost
		port = localport
	} else {
		host = prod
		port = ":443"
	}

	log.Println("Starting " + protocol.VersionTag + " server...")

	d := protocol.NewDispatcher()

	//// Start Webserver WO CORS

	mux := http.NewServeMux()

	pagefs := http.FileServer(http.Dir("public"))
	// gamefs := http.FileServer(http.Dir("fe_overhaul/build"))
	// mux.HandleFunc("/", play)
	// mux.Handle("/static/", gamefs)
	mux.Handle("/", pagefs)
	// mux.Handle("/play/", http.StripPrefix("/play/", gamefs))
	mux.HandleFunc("/ws", // wrap the func to pass the dispatcher
		func(w http.ResponseWriter, req *http.Request) {
			protocol.HandleConnections(w, req, d)
		})

	if host == localhost { // aka env var not set
		log.Fatal(
			http.ListenAndServe(port, mux))
	}

	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	log.Fatal(
		http.ListenAndServeTLS(port, certfile, keyfile, mux))

	// Start Webserver WITH CORS
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", index)
	// mux.HandleFunc("/ws", // wrap the func to pass the dispatcher
	// 	func(w http.ResponseWriter, req *http.Request) {
	// 		protocol.HandleConnections(w, req, d)
	// 	})

	// // c := cors.New(cors.Options{
	// // 	AllowedOrigins: []string{"http://localhost:9009"},
	// // 	// AllowOriginFunc: func(origin string) bool { return true },
	// // })
	// // handler := c.Handler(mux)

	// handler := cors.AllowAll().Handler(mux)

	// if host == localhost { // aka env var not set
	// 	fmt.Println("checkcheck")
	// 	log.Fatal(
	// 		http.ListenAndServe(port, handler))
	// }

	// go http.ListenAndServe(":80", handler)
	// log.Fatal(
	// 	http.ListenAndServeTLS(port, certfile, keyfile, mux))

}

// Redirect all HTTP requests to HTTPS
func redirect(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Redirecting connection to https")
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	log.Printf("redirect to %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// Reject requests that don't have the correct referer header unless they are for the root.
func index(w http.ResponseWriter, req *http.Request) {
	log.Println("Received index request")
	//		if req.Header.Get("Referer") != host &&
	//		req.URL.Path != "/" {
	//		log.Printf("404: %s", req.URL.String())
	//		http.NotFound(w, req)
	//		return
	//	}

	http.FileServer(http.Dir("fe_overhaul/build")).ServeHTTP(w, req)
}
func play(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Received play request")
	//		if req.Header.Get("Referer") != host &&
	//		req.URL.Path != "/" {
	//		log.Printf("404: %s", req.URL.String())
	//		http.NotFound(w, req)
	//		return
	//	}

	// http.FileServer(http.Dir("fe_overhaul/build")).ServeHTTP(w, req)
}
