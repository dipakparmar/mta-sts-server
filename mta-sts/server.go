package mta_sts

// intialize the server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s : %s : %s\n", r.RemoteAddr, r.Host, r.UserAgent())
	// write the response as server is up and running
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Server is up and running")
}

func MtaSTSHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s : %s : %s\n", r.RemoteAddr, r.Host, r.UserAgent())

	// write the response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `version: STSv1
mode: %s
mx: %s
max_age: %s`, "enforce", "mail.example.com", "86400")

}

func Start(Port string) {

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("/.well-known/mta-sts.txt", MtaSTSHandler)

	// create a new server
	s := &http.Server{
		Addr:    ":" + Port,
		Handler: mux,
	}

	// start the server
	log.Printf("Starting server on port %s", Port)
	log.Fatal(s.ListenAndServe())
	os.Exit(2)

}

// func that find the mx record for the domain
func findMX(domain string) ([]string, error) {
	mx, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}
	var mxList []string
	for _, m := range mx {
		mxList = append(mxList, m.Host)
	}
	return mxList, nil
}
