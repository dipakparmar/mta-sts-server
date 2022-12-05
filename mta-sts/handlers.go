package mta_sts

// intialize the server

import (
	"fmt"
	"net/http"
)

type Handler struct {
}

var handler Handler

// RootHandler is the handler for the root path
func (*Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	// if the request is not for the root path, then return 404
	if r.RequestURI != "/" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		// write the response as server is up and running
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server is up and running")
	}
	LogRequest(r)
}

// MtaSTSHandler is the handler for the mta-sts.txt record
func (*Handler) MtaSTSHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, MtaSTSRecord(params.Mode, params.MX, params.MaxAge))
	LogRequest(r)
}
