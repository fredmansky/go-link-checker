package internal

import (
	"encoding/xml"
	"fmt"
	"github.com/fredmansky/go-link-checker/pkg"
	"io"
	"sync"
)

type Sitemap struct {
	URLs []SitemapURL `xml:"url"`
}
type SitemapURL struct {
	Loc string `xml:"loc"`
}

type SitemapIndex struct {
	Sitemaps []SitemapEntry `xml:"sitemap"`
}
type SitemapEntry struct {
	Loc string `xml:"loc"`
}

func FetchLinks(url string, recursive bool, maxRequests int) ([]string, error) {
	resp, err := pkg.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error accessing sitemap: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading sitemap: %v", err)
	}

	// Check if it is a normal sitemap
	var sitemap Sitemap
	if err := xml.Unmarshal(data, &sitemap); err == nil && len(sitemap.URLs) > 0 {
		// Normal sitemap -> return links
		links := make([]string, len(sitemap.URLs))
		for i, url := range sitemap.URLs {
			links[i] = url.Loc
		}
		return links, nil
	}

	// Check if it is a sitemap index
	var index SitemapIndex
	if err := xml.Unmarshal(data, &index); err == nil && len(index.Sitemaps) > 0 {
		if recursive {
			fmt.Printf("✅ %v Sitemaps found\n", len(index.Sitemaps))

			var (
				allLinks []string
				mu       sync.Mutex
				wg       sync.WaitGroup
				sem      = make(chan struct{}, maxRequests)
			)
			for _, entry := range index.Sitemaps {
				wg.Add(1)
				// Wait till place is free in waitgroup
				sem <- struct{}{}

				go func(loc string) {
					defer wg.Done()
					defer func() { <-sem }() // Free up space in waitgroup

					subLinks, subErr := FetchLinks(loc, true, maxRequests)
					if subErr != nil {
						fmt.Printf("❌ Failed to fetch %s: %v\n", loc, subErr)
						return
					}
					fmt.Printf("Fetching links from %s\n", loc)

					mu.Lock()
					allLinks = append(allLinks, subLinks...)
					mu.Unlock()
				}(entry.Loc)
			}

			wg.Wait()

			return allLinks, nil
		}

		// Return sitemap urls if it is recursive flag is false
		sitemapLinks := make([]string, len(index.Sitemaps))
		for i, entry := range index.Sitemaps {
			sitemapLinks[i] = entry.Loc
		}
		return sitemapLinks, nil
	}

	return nil, fmt.Errorf("invalid sitemap format: %s", url)
}
