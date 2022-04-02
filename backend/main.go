package main

import (
	"CoWordle/backend/session"
	"flag"
	"log"
	"net/http"
	"regexp"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var controller *session.Controller = session.NewController()

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got one connection%+v", *r)
		session.ServeUserWebsocket(controller, w, r)
	})
	fileServer := http.FileServer(http.Dir("../frontend/build"))
	fileMatcher := regexp.MustCompile(`\.[a-zA-Z]*$`)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !fileMatcher.MatchString(r.URL.Path) {
			http.ServeFile(w, r, "../frontend/build/index.html")
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})
	// http.HandleFunc("/", home)
	log.Printf("Listening on http://%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
