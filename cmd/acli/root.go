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
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(commandsCmd)
	rootCmd.AddCommand(myCmd)

	rootCmd.PersistentFlags().StringP("profile", "p", "", "configuration profile to use (defaults to the default profile)")
	rootCmd.PersistentFlags().StringP("output", "o", "text", "output format: text or json (json is recommended for programmatic/agent use)")
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
