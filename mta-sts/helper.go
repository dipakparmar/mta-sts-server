package mta_sts

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/mbndr/figlet4go"
	"github.com/spf13/viper"
)

type Config struct {
	Domain  string
	MX      string
	Mode    string
	MaxAge  string
	Port    string
	Verbose bool
}

// func that find the mx record for the domain
func FindMX(domain string) ([]string, error) {
	mx, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}
	var mxList []string
	for _, m := range mx {
		// remove the trailing dot from the mx record
		mxList = append(mxList, m.Host[:len(m.Host)-1])
	}
	return mxList, nil
}

// func that returns the mta-sts.txt record from parameters passed
// mode can be "testing" or "enforce"
// mx is the mx records for the domain - can be a list of mx records
// max_age is the max age of the record in seconds
func MtaSTSRecord(mode string, mx string, max_age string) string {
	var mxList string
	if mx == "" {
		mxList = "none"
	} else {
		mxList = mx
	}
	return fmt.Sprintf("version: STSv1\rmode: %s\rmx: %s\rmax_age: %s", mode, mxList, max_age)
}

// func that creates log for the request
func LogRequest(r *http.Request) {
	log.Printf("%s :: %s :: %s :: %s :: %s :: %s\n", r.Method, r.RequestURI, r.RemoteAddr, r.Referer(), r.Host, r.UserAgent())
}

// func that checks if the domain is valid or not
func isValidDomain(domain string) bool {
	// check if the domain is valid
	if len(domain) < 1 {
		return false
	}
	if domain[0] == '.' {
		return false
	}
	if domain[len(domain)-1] == '.' {
		return false
	}
	if domain == "localhost" {
		return false
	}
	return true
}

// func that checks if the mx record is valid or not
func isValidMX(mx string) bool {
	// check if the mx record is valid
	if len(mx) < 1 {
		return false
	}
	if mx[0] == '.' {
		return false
	}
	if mx[len(mx)-1] == '.' {
		return false
	}
	return true
}

// func that checks if the mode is valid or not
func isValidMode(mode string) bool {
	// check if the mode is valid
	if mode == "testing" || mode == "enforce" || mode == "none" {
		return true
	}
	return false
}

// func that checks if the max age is valid or not
func isValidMaxAge(maxAge string) bool {
	// check if the max age is valid - it should be a number and greater than 0 and less than 31536000 (1 year)
	if _, err := strconv.Atoi(maxAge); err == nil {
		if maxAgeInt, err := strconv.Atoi(maxAge); err == nil {
			if maxAgeInt > 0 && maxAgeInt < 31536000 {
				return true
			}
		}
	}
	return false
}

// func that checks if the config is valid or not and returns error if not valid
func IsValidConfig(config Config) error {

	// gloabla error array variable to store the errors if any
	var err []error

	//  suppress domain name validation if it is not set in the config
	if config.Domain != "" {
		if !isValidDomain(config.Domain) {
			// if verbose is set, then print the input domain
			if config.Verbose {
				log.Printf("invalid provided domain - %s", config.Domain)
			}
			// append the error to the global error variable
			err = append(err, fmt.Errorf("invalid domain - domain can not be localhost or start or end with a dot"))
		}
	} else {

		err = append(err, fmt.Errorf("domain name is empty"))
	}
	if config.MX != "" {
		if !isValidMX(config.MX) {
			// if verbose is set, then print the input mx record
			if config.Verbose {
				log.Printf("invalid provided mx record - %s", config.MX)
			}

			// append the error to the global error variable
			err = append(err, fmt.Errorf("invalid mx record - mx record can not start or end with a dot"))
		}
	} else {

		err = append(err, fmt.Errorf("mx record is empty"))
	}
	if config.Mode != "" {

		if !isValidMode(config.Mode) {
			// if verbose is set, then print the input mode
			if config.Verbose {
				log.Printf("invalid provided mode - %s", config.Mode)
			}
			// append the error to the global error variable
			err = append(err, fmt.Errorf("invalid mode - mode can be testing, enforce or none"))
		}
	} else {

		err = append(err, fmt.Errorf("mode is empty"))
	}
	if config.MaxAge != "" {
		if !isValidMaxAge(config.MaxAge) {
			// if verbose is set, then print the input max age
			if config.Verbose {
				log.Printf("invalid provided max age - %s", config.MaxAge)
			}
			// append the error to the global error variable
			err = append(err, fmt.Errorf("invalid max age - max age should be a number and greater than 0 and less than 31536000 (1 year)"))
		}
	} else {

		err = append(err, fmt.Errorf("max age is empty"))
	}

	// if there are any errors, then return the error
	if len(err) > 0 {
		// use the error array string and create string with all the errors in separate lines and return the error
		var errString string
		for _, e := range err {
			errString = errString + "\n" + e.Error()
		}
		return fmt.Errorf(errString)

	} else {
		return nil
	}
}

// func that generate figlet ascii art for the program name and version - MTA-STS-Server v1.0.0
func PrintFiglet() {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()

	options.FontName = "larry3d"
	options.FontColor = []figlet4go.Color{
		// Colors can be given by default ansi color codes...
		figlet4go.ColorMagenta,
		figlet4go.ColorRed,
		figlet4go.ColorCyan,
		figlet4go.ColorBlue,
	}

	// render the ascii art
	renderStr, _ := ascii.RenderOpts("MTA-STS-Server v0.1.0", options)
	// print the ascii art
	fmt.Println(renderStr)

}

func ReadInConfig() {

	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	defaultCfgFilePath = home + "/.config/mta-sts-server/config.yaml"
	// Read the config file from $HOME/.config/mta-sts-server/config.yaml
	viper.SetConfigFile(defaultCfgFilePath)
	viper.ReadInConfig()

	// store the config in the context for later use
	params.Domain = viper.GetString("domain")
	params.Mode = viper.GetString("mode")
	params.MX = viper.GetString("mx")
	params.MaxAge = viper.GetString("max_age")
	params.Port = viper.GetString("port")
	params.Verbose = viper.GetBool("verbose")
}

func PrintConfig() {
	fmt.Println("Domain: " + params.Domain)
	fmt.Println("Mode: " + params.Mode)
	fmt.Println("MX: " + params.MX)
	fmt.Println("Max Age: " + params.MaxAge)
	fmt.Println("Port: " + params.Port)
	fmt.Println("Verbose: " + strconv.FormatBool(params.Verbose))
}
