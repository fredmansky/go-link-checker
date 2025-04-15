package internal

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
	"github.com/fredmansky/go-link-checker/pkg"
)

type BrokenLink struct {
	url          string
	StatusCode   int
	ResponseTime time.Duration
}

func CheckLinks(links []string, rateLimit int) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var brokenLinks []BrokenLink
	totalLinks := len(links)
	checkedLinks := 0
	checkedLastSecond := 0
	var totalResponseTime time.Duration
	successfulRequests := 0 // Nur erfolgreiche Requests zählen

	maxConcurrentRequests := runtime.NumCPU() * 10
	semaphore := make(chan struct{}, maxConcurrentRequests)

	// Rate limiting
	ticker := time.NewTicker(time.Second / time.Duration(rateLimit))
	defer ticker.Stop()

	stopProgress := make(chan bool)
	go showProgress(&checkedLinks, &checkedLastSecond, totalLinks, stopProgress)

	startTime := time.Now() // Startzeitpunkt für den gesamten Prozess

	for _, link := range links {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(link string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// Warte auf nächstes Rate-Limit-Intervall
			<-ticker.C

			statusCode, responseTime := checkLink(link)

			mu.Lock()
			if statusCode >= http.StatusBadRequest {
				brokenLinks = append(brokenLinks, BrokenLink{url: link, StatusCode: statusCode, ResponseTime: responseTime})
			} else {
				totalResponseTime += responseTime
				successfulRequests++
			}
			checkedLinks++
			checkedLastSecond++
			mu.Unlock()
		}(link)
	}

	wg.Wait()
	stopProgress <- true

	// Gesamtzeit berechnen
	elapsed := time.Since(startTime).Seconds()
	requestsPerSecond := float64(checkedLinks) / elapsed

	// **Korrekte Berechnung der durchschnittlichen Antwortzeit**
	var avgResponseTime float64
	if successfulRequests > 0 {
		avgResponseTime = totalResponseTime.Seconds() / float64(successfulRequests)
	}

	fmt.Printf("\n📊 Durchschnittliche Anfragen pro Sekunde: %.2f\n", requestsPerSecond)
	fmt.Printf("⏳ Durchschnittliche Antwortzeit (nur erfolgreiche Anfragen): %.2f Sekunden\n", avgResponseTime)

	if len(brokenLinks) > 0 {
		fmt.Printf("\n❌ Nicht erreichbare Links: %d\n", len(brokenLinks))
		for _, link := range brokenLinks {
			fmt.Printf("[%d] %s (⏱️ %.2fs)\n", link.StatusCode, link.url, link.ResponseTime.Seconds())
		}
	} else {
		fmt.Println("\n✅ Alle Links sind erreichbar")
	}
}

func checkLink(url string) (int, time.Duration) {
	const (
		maxAttempts        = 3
		defaultWaitSeconds = 1
	)

	for i := 0; i < maxAttempts; i++ {
		start := time.Now()
		resp, err := pkg.HttpClient.Head(url)
		responseTime := time.Since(start)

		if err == nil {
			if resp.StatusCode < http.StatusBadRequest {
				return resp.StatusCode, responseTime // Erfolg
			}
			return resp.StatusCode, responseTime // Fehlerhafte Antwort
		}

		time.Sleep(defaultWaitSeconds * time.Second)
	}

	return 0, 0
}

func showProgress(checked *int, checkedLastSecond *int, total int, stop chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			progress := (*checked * 100) / total
			lps := *checkedLastSecond
			*checkedLastSecond = 0
			fmt.Printf("\rProgress: %d%% | LPS %d links/sec  ", progress, lps)
		}
	}
}
