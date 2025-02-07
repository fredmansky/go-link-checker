package pkg

import (
	"net"
	"net/http"
	"time"
)

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