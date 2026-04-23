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

// printIssueRow prints a single issue row to a tabwriter (shared by board/sprint).
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

// issueFieldHeaders maps Jira field names (as passed to --fields) to display column headers.
var issueFieldHeaders = map[string]string{
	"summary":     "SUMMARY",
	"status":      "STATUS",
	"priority":    "PRIORITY",
	"issuetype":   "TYPE",
	"assignee":    "ASSIGNEE",
	"reporter":    "REPORTER",
	"creator":     "CREATOR",
	"created":     "CREATED",
	"updated":     "UPDATED",
	"duedate":     "DUE DATE",
	"labels":      "LABELS",
	"resolution":  "RESOLUTION",
	"components":  "COMPONENTS",
	"fixversions": "FIX VERSIONS",
	"versions":    "VERSIONS",
	"project":     "PROJECT",
	"parent":      "PARENT",
}

// issueFieldHeader returns the column header for a Jira field name.
func issueFieldHeader(name string) string {
	if h, ok := issueFieldHeaders[strings.ToLower(name)]; ok {
		return h
	}
	return strings.ToUpper(name)
}

// renderIssueField returns the display value for a given Jira field on an issue.
func renderIssueField(name string, issue jira.IssueDetailed) string {
	f := issue.Fields
	switch strings.ToLower(name) {
	case "summary":
		return f.Summary
	case "status":
		if f.Status != nil {
			return f.Status.Name
		}
	case "priority":
		if f.Priority != nil {
			return f.Priority.Name
		}
	case "issuetype":
		if f.IssueType != nil {
			return f.IssueType.Name
		}
	case "assignee":
		if f.Assignee != nil {
			return f.Assignee.DisplayName
		}
	case "reporter":
		if f.Reporter != nil {
			return f.Reporter.DisplayName
		}
	case "creator":
		if f.Creator != nil {
			return f.Creator.DisplayName
		}
	case "created":
		return f.Created
	case "updated":
		return f.Updated
	case "duedate":
		return f.DueDate
	case "labels":
		return strings.Join(f.Labels, ",")
	case "resolution":
		if f.Resolution != nil {
			return f.Resolution.Name
		}
	case "components":
		names := make([]string, 0, len(f.Components))
		for _, c := range f.Components {
			names = append(names, c.Name)
		}
		return strings.Join(names, ",")
	case "fixversions":
		names := make([]string, 0, len(f.FixVersions))
		for _, v := range f.FixVersions {
			names = append(names, v.Name)
		}
		return strings.Join(names, ",")
	case "versions":
		names := make([]string, 0, len(f.Versions))
		for _, v := range f.Versions {
			names = append(names, v.Name)
		}
		return strings.Join(names, ",")
	case "project":
		if f.Project != nil {
			return f.Project.Key
		}
	case "parent":
		if f.Parent != nil {
			return f.Parent.Key
		}
	}
	return ""
}
