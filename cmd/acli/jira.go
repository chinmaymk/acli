package acli

import (
	"github.com/spf13/cobra"
)

var jiraCmd = &cobra.Command{
	Use:     "jira",
	Aliases: []string{"j"},
	Short:   "Interact with Jira Cloud",
	Long:    "Manage Jira projects, issues, boards, sprints, and more.",
	RunE:    helpRunE,
}

func init() {
	// Core resources
	jiraCmd.AddCommand(jiraIssueCmd)
	jiraCmd.AddCommand(jiraProjectCmd)
	jiraCmd.AddCommand(jiraBoardCmd)
	jiraCmd.AddCommand(jiraSprintCmd)
	jiraCmd.AddCommand(jiraEpicCmd)
	jiraCmd.AddCommand(jiraBacklogMoveCmd)

	// Search & Filters
	jiraCmd.AddCommand(jiraSearchCmd)
	jiraCmd.AddCommand(jiraFilterCmd)

	// Users & Groups
	jiraCmd.AddCommand(jiraUserCmd)
	jiraCmd.AddCommand(jiraGroupCmd)

	// Dashboards
	jiraCmd.AddCommand(jiraDashboardCmd)
}
