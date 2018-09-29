package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// observeconnection is the connection detector. It tries to do an HTTP GET against
// probeURL and if *anything* comes back we consider ourselves to be online, otherwise
// some network issues prevents us from doing the GET and we are likely offline.
func observeconnection(cremote, clocal string, constat chan string) {
	// the endpoint we're using to check if we're online or offline:
	var probeURL string
	ccurrent := cremote
	for {
		probeURL = getAPIServerURL(ccurrent)
		client := http.Client{Timeout: time.Duration(ProbeTimeoutSeconds * time.Second)}
		resp, err := client.Get(probeURL)
		if err != nil {
			fmt.Printf("\x1b[93mConnection detection [%v], probe resulted in %v\x1b[0m\n", StatusOffline, err)
			ccurrent = clocal
			constat <- StatusOffline
			continue
		}
		fmt.Printf("\x1b[93mConnection detection [%v], probe %v resulted in %v\x1b[0m\n", StatusOnline, probeURL, resp.Status)
		constat <- StatusOnline
		time.Sleep(CheckConnectionDelaySeconds * time.Second)
	}
}

// getAPIServerURL looks up the API Server url of the kubectx provided.
func getAPIServerURL(kubectx string) string {
	clustername := clusterfromcontext(kubectx)
	probeURL, err := kubectl(false, false, "config", "view",
		"--output=jsonpath='{.clusters[?(@.name == \""+clustername+"\")]..server}'")
	if err != nil {
		displayerr("Can't cuddle the cluster", err)
	}
	probeURL = strings.Trim(probeURL, "'")
	return probeURL
}

// clusterfromcontext extracts the cluster name part from
// a context name, asssuming it is in the OpenShift format.
func clusterfromcontext(context string) string {
	// In OpenShift, the context naming format is:
	// $PROJECT/$CLUSTERNAME/$USER for example:
	// mh9sandbox/api-pro-us-east-1-openshift-com:443/mhausenb
	re := regexp.MustCompile("(.*)/(.*)/(.*)")
	return re.FindStringSubmatch(context)[2]
}
