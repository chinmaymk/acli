package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "acli",
	Short: "Atlassian CLI - manage Jira, Confluence, and Bitbucket from the terminal",
	Long: `acli is a command-line interface for Atlassian Cloud products.

Manage your Jira issues, Confluence pages, and Bitbucket repositories
without leaving the terminal.`,
}

func init() {
	rootCmd.AddCommand(jiraCmd)
	rootCmd.AddCommand(confluenceCmd)
	rootCmd.AddCommand(bitbucketCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringP("profile", "p", "default", "configuration profile to use")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("acli %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

func Execute() error {
	return rootCmd.Execute()
}
