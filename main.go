package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	serverPort = flag.String("port", "8080", "port to listen on")
	stsDomain  = flag.String("domain", "", "domain for sts")
	stsMode    = flag.String("stsMode", "testing", "STS mode to use. Options are 'testing' or 'enforce'")
	stsMX      = flag.String("stsMX", "", "MX to use for STS. Comma separated list")
	stsMaxAge  = flag.String("stsMaxAge", "2419200", "STS max age in seconds")
)

func main() {
	flag.Parse()

	var domain = os.Getenv("DOMAIN")
	var mode = os.Getenv("STS_MODE")
	var mx = os.Getenv("STS_MX")
	var maxAge = os.Getenv("STS_MAX_AGE")

	if domain != "" {
		*stsDomain = domain
	}

	if mode != "" {
		*stsMode = mode
	}

	if mx != "" {
		*stsMX = mx
	}

	if maxAge != "" {
		// convert maxAge to int and then assign it to stsMaxAge
		*stsMaxAge = maxAge
	}

	// if stsMX is not set, the find the mx record for the domain and set it
	if *stsMX == "" {
		mx, err := findMX(*stsDomain)
		if err != nil {
			log.Fatal(err)
		}
		// convert the mx as string without [] and remove the last dot if present
		*stsMX = fmt.Sprintf("%s", mx)
		*stsMX = (*stsMX)[1 : len(*stsMX)-1]
		if (*stsMX)[len(*stsMX)-1:] == "." {
			*stsMX = (*stsMX)[:len(*stsMX)-1]
		}

	}

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/.well-known/mta-sts.txt", mtaSTSHandler)
	// create a new server
	s := &http.Server{
		Addr:    ":" + *serverPort,
		Handler: mux,
	}
	// start the server
	log.Printf("Starting server on port %s", *serverPort)
	log.Fatal(s.ListenAndServe())
	os.Exit(2)
}

// fmt.Fprintln(os.Stderr, http.ListenAndServe(":"+*serverPort, nil))

// func that handles the request and response for /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s : %s : %s\n", r.RemoteAddr, r.Host, r.UserAgent())
	// write the response as server is up and running
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Server is up and running for domain "+*stsDomain)
}

// func that handles the request and response for /.well-known/mta-sts.txt path
func mtaSTSHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s : %s : %s\n", r.RemoteAddr, r.Host, r.UserAgent())
	// check if the domain is set
	if *stsDomain == "" {
		log.Println("Domain not set")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Domain not set")
		os.Exit(2)
	}
	// check if the sts mode is set
	if *stsMode == "" {
		log.Println("STS mode not set")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "STS mode not set")
		os.Exit(2)
	}
	// check if the sts mode is valid
	if *stsMode != "testing" && *stsMode != "enforce" {
		log.Println("STS mode not valid")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "STS mode not valid")
		os.Exit(2)
	}
	// check if the sts max age is set
	if *stsMaxAge == "" {
		log.Println("STS max age not set")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "STS max age not set")
		os.Exit(2)
	}
	// check if the sts mx is set
	if *stsMX == "" {
		log.Println("STS mx not set")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "STS mx not set")
		os.Exit(2)
	}

	// write the response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `version: STSv1
mode: %s
mx: %s
max_age: %s`, *stsMode, *stsMX, *stsMaxAge)

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
