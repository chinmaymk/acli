package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bitbucketCmd = &cobra.Command{
	Use:     "bitbucket",
	Aliases: []string{"bb"},
	Short:   "Interact with Bitbucket Cloud",
	Long:    "Manage Bitbucket repositories, pull requests, and pipelines.",
	RunE:    helpRunE,
}

var bbRepoCmd = &cobra.Command{
	Use:     "repo",
	Short:   "Manage repositories",
	Aliases: []string{"r"},
	RunE:    helpRunE,
}

var bbPRCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
	RunE:  helpRunE,
}

var bbPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Short:   "Manage pipelines",
	Aliases: []string{"pipe"},
	RunE:    helpRunE,
}

func init() {
	bitbucketCmd.AddCommand(bbRepoCmd)
	bitbucketCmd.AddCommand(bbPRCmd)
	bitbucketCmd.AddCommand(bbPipelineCmd)

	bbRepoCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "List repositories",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: list Bitbucket repos")
		},
	})

	bbPRCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "List pull requests",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: list Bitbucket PRs")
		},
	})
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "get [pr-id]",
		Short: "Get pull request details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("TODO: get Bitbucket PR %s\n", args[0])
		},
	})
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a pull request",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: create Bitbucket PR")
		},
	})
}
