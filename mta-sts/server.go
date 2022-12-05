package mta_sts

// intialize the server

import (
	"log"
	"net/http"
	"os"
)

type MtaSTSRecordParams struct {
	Domain string
	Mode   string
	MX     []string
	MaxAge string
}

var params MtaSTSRecordParams

func Start(Port string, Domain string, Mode string, MX []string, MaxAge string) {

	// set the parameters to the context
	params.Domain = Domain
	params.Mode = Mode
	params.MX = MX
	params.MaxAge = MaxAge

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.RootHandler)
	mux.HandleFunc("/.well-known/mta-sts.txt", handler.MtaSTSHandler)

	// create a new server
	server := &http.Server{
		Addr:    ":" + Port,
		Handler: mux,
	}

	PrintFiglet()
	// start the server
	log.Printf("Starting server on port %s", Port)
	log.Fatal(server.ListenAndServe())
	os.Exit(2)

}
