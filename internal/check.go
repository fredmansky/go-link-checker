package internal

import (
	"fmt"
	"github.com/fredmansky/go-link-checker/pkg"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func CheckLinks(links []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var brokenLinks []string
	totalLinks := len(links)
	checkedLinks := 0
	checkedLastSecond := 0

	maxConcurrentRequests := runtime.NumCPU() * 10
	semaphore := make(chan struct{}, maxConcurrentRequests)

	stopProgress := make(chan bool)
	go showProgress(&checkedLinks, &checkedLastSecond, totalLinks, stopProgress)

	for _, link := range links {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(link string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if !checkLink(link) {
				mu.Lock()
				brokenLinks = append(brokenLinks, link)
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
			fmt.Println(link)
		}
	} else {
		fmt.Println("\n✅ All links passed")
	}
}

func checkLink(url string) bool {
	const (
		maxAttempts          = 3
		rateLimitWaitSeconds = 5
		defaultWaitSeconds   = 1
	)

	for i := 0; i < maxAttempts; i++ {
		resp, err := pkg.HttpClient.Head(url)

		// Return true if request was successful and no error occurred
		if err == nil && resp.StatusCode < http.StatusBadRequest {
			return true
		}

		if err != nil {
			time.Sleep(defaultWaitSeconds * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("Server rate-limited us! Waiting %ds...\n", rateLimitWaitSeconds)
			time.Sleep(rateLimitWaitSeconds * time.Second)
		} else {
			time.Sleep(defaultWaitSeconds * time.Second)
		}
	}

	return false
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
