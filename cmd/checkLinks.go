/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/fredmansky/go-link-checker/internal"
)

// checkLinksCmd represents the checkLinks command
var checkLinksCmd = &cobra.Command{
	Use:   "check-links",
	Short: "Check links if they are broken",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		fmt.Println("Check links from:", url)
		links, err := internal.FetchLinks(url, true)

		if err != nil {
			fmt.Println("❌ Error fetching links:", err)
			return
		}

		fmt.Printf("✅ Successfully found %d links\n", len(links))

		internal.CheckLinks(links);
	},
}

func init() {
	rootCmd.AddCommand(checkLinksCmd)
}
