package acli

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// openURL opens a URL in the user's default browser.
func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func init() {
	// jira issue open <issue-key>
	jiraIssueOpenCmd := &cobra.Command{
		Use:   "open <issue-key>",
		Short: "Open a Jira issue in the browser",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := getProfile(cmd)
			if err != nil {
				return err
			}
			baseURL := strings.TrimRight(profile.AtlassianURL, "/")
			issueURL := fmt.Sprintf("%s/browse/%s", baseURL, args[0])
			fmt.Fprintf(cmd.OutOrStdout(), "Opening %s\n", issueURL)
			return openURL(issueURL)
		},
	}
	jiraIssueCmd.AddCommand(jiraIssueOpenCmd)

	// bb pr open [workspace] <repo-slug> <pr-id>
	bbPROpenCmd := &cobra.Command{
		Use:   "open [workspace] <repo-slug> <pr-id>",
		Short: "Open a pull request in the browser",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, idStr, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", idStr)
			}
			prURL := fmt.Sprintf("https://bitbucket.org/%s/%s/pull-requests/%d", workspace, repoSlug, prID)
			fmt.Fprintf(cmd.OutOrStdout(), "Opening %s\n", prURL)
			return openURL(prURL)
		},
	}
	bbPRCmd.AddCommand(bbPROpenCmd)

	// bb repo open [workspace] <repo-slug>
	bbRepoOpenCmd := &cobra.Command{
		Use:   "open [workspace] <repo-slug>",
		Short: "Open a repository in the browser",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, err := resolveWorkspaceAndRepo(cmd, args)
			if err != nil {
				return err
			}
			repoURL := fmt.Sprintf("https://bitbucket.org/%s/%s", workspace, repoSlug)
			fmt.Fprintf(cmd.OutOrStdout(), "Opening %s\n", repoURL)
			return openURL(repoURL)
		},
	}
	bbRepoCmd.AddCommand(bbRepoOpenCmd)

	// jira project open <project-key>
	jiraProjectOpenCmd := &cobra.Command{
		Use:   "open <project-key>",
		Short: "Open a Jira project board in the browser",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := getProfile(cmd)
			if err != nil {
				return err
			}
			baseURL := strings.TrimRight(profile.AtlassianURL, "/")
			projectURL := fmt.Sprintf("%s/jira/software/projects/%s/board", baseURL, args[0])
			fmt.Fprintf(cmd.OutOrStdout(), "Opening %s\n", projectURL)
			return openURL(projectURL)
		},
	}
	jiraProjectCmd.AddCommand(jiraProjectOpenCmd)
}
