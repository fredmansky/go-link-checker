package internal

import (
	"fmt"
	"runtime"
	"net/http"
	"sync"
	"time"
	"github.com/fredmansky/go-link-checker/pkg"
)

type BrokenLink struct {
	url        string
	StatusCode int
}

func CheckLinks(links []string, rateLimit int) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var brokenLinks []BrokenLink
	totalLinks := len(links)
	checkedLinks := 0
	checkedLastSecond := 0

	maxConcurrentRequests := runtime.NumCPU() * 10
	semaphore := make(chan struct{}, maxConcurrentRequests)

	// Rate limiting
	ticker := time.NewTicker(time.Second / time.Duration(rateLimit))
	defer ticker.Stop()

	stopProgress := make(chan bool)
	go showProgress(&checkedLinks, &checkedLastSecond, totalLinks, stopProgress)

	for _, link := range links {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(link string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// wait till ticker makes new token available
			<-ticker.C

			statusCode := checkLink(link)

			if statusCode >= http.StatusBadRequest {
				mu.Lock()
				brokenLinks = append(brokenLinks, BrokenLink{url: link, StatusCode: statusCode})
				mu.Unlock()
			}

			mu.Lock()
			checkedLinks++
			checkedLastSecond++
			mu.Unlock()
		}(link)
	}

	wg.Wait()

	stopProgress <- true

	if brokenLinksLen := len(brokenLinks); brokenLinksLen > 0 {
		fmt.Printf("\n❌ Not reachable links: %d\n", len(brokenLinks))
		for _, link := range brokenLinks {
			fmt.Printf("[%d] %s\n", link.StatusCode, link.url)
		}
	} else {
		fmt.Println("\n✅ All links passed")
	}
}

func checkLink(url string) int {
	const (
		maxAttempts          = 3
		rateLimitWaitSeconds = 5
		defaultWaitSeconds   = 1
	)

	for i := 0; i < maxAttempts; i++ {
		resp, err := pkg.HttpClient.Head(url)

		if err == nil {
			if resp.StatusCode < http.StatusBadRequest {
				return resp.StatusCode // Success
			}
			return resp.StatusCode // Failure
		}

		time.Sleep(defaultWaitSeconds * time.Second)
	}

	return 0
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
			fmt.Printf("\rProgress: %d%% | LPS %d links/sec  ", progress, lps) // spaces needed to prevent wrong output (secc)
		}
	}
}
