/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/fredmansky/go-link-checker/pkg"
	"github.com/spf13/cobra"
)

var (
	username string
	password string
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-link-checker",
	Short: "A link checker working with seomatic sitemap",
	Long: "A link checker that fetches urls based on your seomatic sitemap and checks all urls.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Basic Auth flags (persistent = available for all subcommands)
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username for Basic Auth")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for Basic Auth")

	// Set credentials before any command runs
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if username != "" && password != "" {
			pkg.SetBasicAuth(username, password)
		}
	}
}


