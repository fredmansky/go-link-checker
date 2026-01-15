package pkg

import (
	"encoding/base64"
	"net"
	"net/http"
	"time"
)

// Pre-computed Basic Auth header (empty if no auth)
var basicAuthHeader string

// SetBasicAuth sets the Basic Auth credentials for all HTTP requests.
// The header is pre-computed once for performance.
func SetBasicAuth(username, password string) {
	credentials := username + ":" + password
	basicAuthHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}

var HttpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns: 1000,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout: 60 * time.Second,
		TLSHandshakeTimeout: 15 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
		ForceAttemptHTTP2: true,
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
			KeepAlive: 180 * time.Second,
			DualStack: true,
		}).DialContext,
		DisableKeepAlives: false,
		DisableCompression: false,
	},
	Timeout: 40 * time.Second,
}

// Get performs a GET request with optional Basic Auth.
func Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if basicAuthHeader != "" {
		req.Header.Set("Authorization", basicAuthHeader)
	}
	return HttpClient.Do(req)
}

// Head performs a HEAD request with optional Basic Auth.
func Head(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	if basicAuthHeader != "" {
		req.Header.Set("Authorization", basicAuthHeader)
	}
	return HttpClient.Do(req)
}