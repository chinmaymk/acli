package acli

import (
	"fmt"
	"os"
	"strings"
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

// defaultBBProject returns the --project flag value if set, otherwise falls back to the profile default BB project.
func defaultBBProject(cmd *cobra.Command) (string, error) {
	project, _ := cmd.Flags().GetString("project")
	if project != "" {
		return project, nil
	}
	profile, err := getProfile(cmd)
	if err != nil {
		return "", err
	}
	return profile.Defaults.BBProject, nil
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

// firstLine returns the first line of a string.
func firstLine(s string) string {
	for i, c := range s {
		if c == '\n' || c == '\r' {
			return s[:i]
		}
	}
	return s
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
	_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
		issue.Key, issueType, status, priority, assignee, issue.Fields.Summary)
}

// issueColumn describes a table column for a given Jira field.
type issueColumn struct {
	header  string
	extract func(jira.IssueDetailed) string
}

// issueColumnFor returns the column definition for a Jira field name.
// Unknown fields fall back to an uppercased header with an empty value
// (the custom field data isn't parsed into IssueFields).
func issueColumnFor(field string) issueColumn {
	switch strings.ToLower(field) {
	case "summary":
		return issueColumn{"SUMMARY", func(i jira.IssueDetailed) string { return i.Fields.Summary }}
	case "issuetype", "type":
		return issueColumn{"TYPE", func(i jira.IssueDetailed) string {
			if i.Fields.IssueType == nil {
				return ""
			}
			return i.Fields.IssueType.Name
		}}
	case "status":
		return issueColumn{"STATUS", func(i jira.IssueDetailed) string {
			if i.Fields.Status == nil {
				return ""
			}
			return i.Fields.Status.Name
		}}
	case "priority":
		return issueColumn{"PRIORITY", func(i jira.IssueDetailed) string {
			if i.Fields.Priority == nil {
				return ""
			}
			return i.Fields.Priority.Name
		}}
	case "assignee":
		return issueColumn{"ASSIGNEE", func(i jira.IssueDetailed) string {
			if i.Fields.Assignee == nil {
				return ""
			}
			return i.Fields.Assignee.DisplayName
		}}
	case "reporter":
		return issueColumn{"REPORTER", func(i jira.IssueDetailed) string {
			if i.Fields.Reporter == nil {
				return ""
			}
			return i.Fields.Reporter.DisplayName
		}}
	case "creator":
		return issueColumn{"CREATOR", func(i jira.IssueDetailed) string {
			if i.Fields.Creator == nil {
				return ""
			}
			return i.Fields.Creator.DisplayName
		}}
	case "resolution":
		return issueColumn{"RESOLUTION", func(i jira.IssueDetailed) string {
			if i.Fields.Resolution == nil {
				return ""
			}
			return i.Fields.Resolution.Name
		}}
	case "project":
		return issueColumn{"PROJECT", func(i jira.IssueDetailed) string {
			if i.Fields.Project == nil {
				return ""
			}
			return i.Fields.Project.Key
		}}
	case "labels":
		return issueColumn{"LABELS", func(i jira.IssueDetailed) string {
			return strings.Join(i.Fields.Labels, ",")
		}}
	case "created":
		return issueColumn{"CREATED", func(i jira.IssueDetailed) string { return i.Fields.Created }}
	case "updated":
		return issueColumn{"UPDATED", func(i jira.IssueDetailed) string { return i.Fields.Updated }}
	case "duedate":
		return issueColumn{"DUE DATE", func(i jira.IssueDetailed) string { return i.Fields.DueDate }}
	}
	return issueColumn{strings.ToUpper(field), func(jira.IssueDetailed) string { return "" }}
}

// printIssueTable renders a tabular view of search results using only the
// requested fields as columns. KEY is always the leading column.
func printIssueTable(issues []jira.IssueDetailed, fields []string) {
	cols := []issueColumn{{"KEY", func(i jira.IssueDetailed) string { return i.Key }}}
	for _, f := range fields {
		if strings.EqualFold(f, "key") {
			continue
		}
		cols = append(cols, issueColumnFor(f))
	}

	w := newTabWriter()
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.header
	}
	_, _ = fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, issue := range issues {
		row := make([]string, len(cols))
		for i, c := range cols {
			row[i] = c.extract(issue)
		}
		_, _ = fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	_ = w.Flush()
}
