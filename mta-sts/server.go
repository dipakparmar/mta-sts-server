package mta_sts

import (
	"log"
	"net/http"
	"os"
)

var params Config = Config{
	Port:    "8080",
	Verbose: false,
}

var defaultCfgFilePath string

// Start starts the server and listens on the port specified in the config file or the default port 8080
func Start() {

	// read the config file
	ReadInConfig()

	// print the config to the console if verbose is set to true in the config file
	if params.Verbose {
		PrintConfig()
	}

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.RootHandler)                          // root handler
	mux.HandleFunc("/.well-known/mta-sts.txt", handler.MtaSTSHandler) // mta-sts handler for the mta-sts.txt file

	// create a new server
	server := &http.Server{
		Addr:    ":" + params.Port,
		Handler: mux,
	}

	// print the figlet logo to the console
	PrintFiglet()

	// start the server and listen on the port specified in the config file or the default port 8080
	log.Printf("Starting server on port %s", params.Port)
	log.Fatal(server.ListenAndServe())
	os.Exit(2)
}
