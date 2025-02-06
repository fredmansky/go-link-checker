package pkg

import (
	"net"
	"net/http"
	"time"
)

// Globaler HTTP-Client für alle Anfragen
var HttpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        500,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     300 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
	},
	Timeout: 10 * time.Second,
}
