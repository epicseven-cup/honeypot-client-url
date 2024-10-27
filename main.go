package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Only assuming desktop for now

// Mapping of operating systemm and the template that should be
var osMap = map[string]string{
	"window": "(Windows NT x.y; rv:10.0)",
	"mac":    "(Macintosh; Intel Mac OS X x.y; rv:10.0)",
	"linux":  "(X11; Linux x86_64; rv:10.0)",
}

var broswerMap = map[string]string{
	"firefox": "Mozilla/5.0 <operating-system> Gecko/20100101 Firefox/10.0",
}

const (
	LOGGER_FILE_NAME             = "honeypot-client.log"
	OPERATING_SYSTEM_PLACEHOLDER = "<operating-system>"
)

// mimic broswer
// only supports firefox for now
// honeypot-client operating-system browser-type url
func main() {

	if len(os.Args) != 4 {
		fmt.Println("Expected arguments of length 4")
	}

	// url that will be used to access as the honey pot client
	osType, broswer, url := os.Args[1], os.Args[2], os.Args[3]
	logger := log.Default()
	file, err := os.OpenFile(LOGGER_FILE_NAME, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	logger.SetOutput(file)

	visted := map[string]bool{}
	// http client that will be act like a honeypot
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			_, ok := visted[req.URL.RawPath]
			if ok {
				return errors.New("Duplicate redirect urls")
			}
			logger.Println("Status: ", req.Response.Status)
			logger.Println("HEADER: ", req.Header)
			logger.Println("PATH: ", req.URL)
			logger.Println("BODY: ", req.Body)
			return nil
		},
	}

	// Creating the request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	// Create User Agent
	ua, ok := broswerMap[broswer]

	if !ok {
		fmt.Println("Operating System is not supported")
	}

	// opearting system tag
	osType, ok = osMap[osType]

	if !ok {
		fmt.Println("Browser type is not supported")
	}

	// Replace the operating system in template
	strings.Replace(ua, OPERATING_SYSTEM_PLACEHOLDER, osType, 1)
	req.Header.Set("User-Agent", ua)

	// Sending the rqeuest and logs everything

	logger.Println("Sending request now")
	// Checking the first respond
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	logger.Println("Final STATUS: ", res.Status)
	logger.Println("Final HEADER: ", res.Header)
	logger.Println("Final BODY: ", res.Body)
}
