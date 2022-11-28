package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	port = flag.String("port", "8080", "port to listen on")
)

func main() {
	var defaultPort = "8080"
	flag.Parse()
	if port == nil {
		port = &defaultPort
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s : %s : %s\n", r.RemoteAddr, r.Host, r.UserAgent())
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	http.HandleFunc("/.well-known/mta-sts.txt", mtaSTSHandler)

	fmt.Fprintln(os.Stderr, http.ListenAndServe(":"+*port, nil))
	os.Exit(2)
}

// func that handles the request and response for /.well-known/mta-sts.txt path
func mtaSTSHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s : %s : %s\n", r.RemoteAddr, r.Host, r.UserAgent())
	var mtaSTS = `version: STSv1
mode: testing
mx: mail.example.com
max_age: 86400`

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(mtaSTS)))
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, mtaSTS)
}
