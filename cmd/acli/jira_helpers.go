package acli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/chinmaymk/acli/internal/config"
	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

func getProfile(cmd *cobra.Command) (config.Profile, error) {
	profileName, _ := cmd.Flags().GetString("profile")
	cfg, err := config.Load()
	if err != nil {
		return config.Profile{}, fmt.Errorf("loading config: %w", err)
	}
	return cfg.GetProfile(profileName)
}

func getJiraClient(cmd *cobra.Command) (*jira.Client, error) {
	profile, err := getProfile(cmd)
	if err != nil {
		return nil, err
	}
	return jira.NewClient(profile)
}

func getBitbucketClient(cmd *cobra.Command) (*bitbucket.Client, error) {
	profile, err := getProfile(cmd)
	if err != nil {
		return nil, err
	}
	return bitbucket.NewClient(profile)
}

// defaultProject returns the flag value if set, otherwise falls back to the profile default.
func defaultProject(cmd *cobra.Command) (string, error) {
	project, _ := cmd.Flags().GetString("project")
	if project != "" {
		return project, nil
	}
	profile, err := getProfile(cmd)
	if err != nil {
		return "", err
	}
	return profile.Defaults.Project, nil
}

// defaultWorkspace returns the arg if provided, otherwise falls back to the profile default.
// Returns the workspace and an error if no workspace could be resolved.
func defaultWorkspace(cmd *cobra.Command, args []string, argIndex int) (string, error) {
	if argIndex < len(args) {
		return args[argIndex], nil
	}
	profile, err := getProfile(cmd)
	if err != nil {
		return "", err
	}
	if profile.Defaults.Workspace != "" {
		return profile.Defaults.Workspace, nil
	}
	return "", fmt.Errorf("workspace is required: provide it as an argument or set a default with 'acli config set-defaults'")
}

// resolveWorkspaceAndRepo handles the common pattern of [workspace] <repo> args.
// With 2 args: workspace=args[0], repo=args[1].
// With 1 arg: workspace from profile default, repo=args[0].
func resolveWorkspaceAndRepo(cmd *cobra.Command, args []string) (string, string, error) {
	if len(args) >= 2 {
		return args[0], args[1], nil
	}
	workspace, err := defaultWorkspace(cmd, nil, 0)
	if err != nil {
		return "", "", err
	}
	return workspace, args[0], nil
}

// resolveWorkspaceRepoAndID handles the pattern of [workspace] <repo> <id> args.
// With 3 args: workspace=args[0], repo=args[1], id=args[2].
// With 2 args: workspace from profile default, repo=args[0], id=args[1].
func resolveWorkspaceRepoAndID(cmd *cobra.Command, args []string) (string, string, string, error) {
	if len(args) >= 3 {
		return args[0], args[1], args[2], nil
	}
	workspace, err := defaultWorkspace(cmd, nil, 0)
	if err != nil {
		return "", "", "", err
	}
	return workspace, args[0], args[1], nil
}

func newTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// truncate truncates a string to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// printIssueRow prints a single issue row to a tabwriter (shared by board/sprint/search).
func printIssueRow(w *tabwriter.Writer, issue jira.IssueDetailed) {
	issueType := ""
	if issue.Fields.IssueType != nil {
		issueType = issue.Fields.IssueType.Name
	}
	status := ""
	if issue.Fields.Status != nil {
		status = issue.Fields.Status.Name
	}
	priority := ""
	if issue.Fields.Priority != nil {
		priority = issue.Fields.Priority.Name
	}
	assignee := ""
	if issue.Fields.Assignee != nil {
		assignee = issue.Fields.Assignee.DisplayName
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
		issue.Key, issueType, status, priority, assignee, issue.Fields.Summary)
}
