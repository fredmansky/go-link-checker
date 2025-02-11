package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/fredmansky/go-link-checker/internal"
)

var fetchLinksCmd = &cobra.Command{
	Use:   "fetch-links [URL]",
	Short: "Returns all URLs based on the provided website",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		fmt.Println("Fetching links from:", url)
		links, err := internal.FetchLinks(url, true, rateLimit)

		if err != nil {
			fmt.Println("Error fetching links:", err)
			return
		}

		fmt.Printf("\nAmount of found links: %d\n\n", len(links))
		// print out links
		for _, link := range links {
			fmt.Println(link)
		}
	},
}

func init() {
	fetchLinksCmd.Flags().IntVarP(&rateLimit, "rate-limit", "l", 200, "Limit link checks per second")
	rootCmd.AddCommand(fetchLinksCmd)
}
