package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var jiraCmd = &cobra.Command{
	Use:     "jira",
	Aliases: []string{"j"},
	Short:   "Interact with Jira Cloud",
	Long:    "Manage Jira projects, issues, boards, and sprints.",
	RunE:    helpRunE,
}

var jiraIssueCmd = &cobra.Command{
	Use:     "issue",
	Short:   "Manage issues",
	Aliases: []string{"i"},
	RunE:    helpRunE,
}

var jiraProjectCmd = &cobra.Command{
	Use:     "project",
	Short:   "Manage projects",
	Aliases: []string{"p"},
	RunE:    helpRunE,
}

var jiraBoardCmd = &cobra.Command{
	Use:     "board",
	Short:   "Manage boards",
	Aliases: []string{"b"},
	RunE:    helpRunE,
}

var jiraSprintCmd = &cobra.Command{
	Use:     "sprint",
	Short:   "Manage sprints",
	Aliases: []string{"s"},
	RunE:    helpRunE,
}

func init() {
	jiraCmd.AddCommand(jiraIssueCmd)
	jiraCmd.AddCommand(jiraProjectCmd)
	jiraCmd.AddCommand(jiraBoardCmd)
	jiraCmd.AddCommand(jiraSprintCmd)

	jiraIssueCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "List issues",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: list Jira issues")
		},
	})
	jiraIssueCmd.AddCommand(&cobra.Command{
		Use:   "get [issue-key]",
		Short: "Get issue details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("TODO: get Jira issue %s\n", args[0])
		},
	})
	jiraIssueCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create an issue",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: create Jira issue")
		},
	})
}
