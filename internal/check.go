package internal

import (
	"fmt"
	"github.com/fredmansky/go-link-checker/pkg"
	"sync"
	"time"
)

func CheckLinks(links []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var brokenLinks []string

	maxConcurrentRequests := 50
	semaphore := make(chan struct{}, maxConcurrentRequests)

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
		}(link)
	}

	wg.Wait()

	if brokenLinksLen := len(brokenLinks); brokenLinksLen > 0 {
		fmt.Printf("🚨 Not reachable links: %d\n", len(brokenLinks))
		for _, link := range brokenLinks {
			fmt.Println(link)
		}
	} else {
		fmt.Printf("✅ All links passed")
	}
}

func checkLink(url string) bool {
	for i := 0; i < 3; i++ {
		resp, err := pkg.HttpClient.Head(url)
		if err == nil && resp.StatusCode < 400 {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}
