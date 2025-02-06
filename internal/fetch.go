package internal

import (
	"encoding/xml"
	"fmt"
	"io"
	"github.com/fredmansky/go-link-checker/pkg"
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

func FetchLinks(url string, recursive bool) ([]string, error) {
	resp, err := pkg.HttpClient.Get(url)
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
		fmt.Println("Successfully recogniced a sitemap index...")

		if recursive {
			var allLinks []string
			for _, sitemapEntry := range index.Sitemaps {
				subLinks, err := FetchLinks(sitemapEntry.Loc, recursive)
				if err != nil {
					fmt.Printf("⚠️ Warning: Failed to fetch %s: %v\n", sitemapEntry.Loc, err)
					continue
				}
				allLinks = append(allLinks, subLinks...)
			}
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
