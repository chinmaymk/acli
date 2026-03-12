package acli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbPRCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
	RunE:  helpRunE,
}

func init() {
	// pr list
	prListCmd := &cobra.Command{
		Use:     "list [workspace] <repo-slug>",
		Short:   "List pull requests",
		Aliases: []string{"ls"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, err := resolveWorkspaceAndRepo(cmd, args)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			state, _ := cmd.Flags().GetString("state")
			page, _ := cmd.Flags().GetInt("page")
			pagelen, _ := cmd.Flags().GetInt("pagelen")
			all, _ := cmd.Flags().GetBool("all")
			prs, err := client.ListPullRequests(workspace, repoSlug, &bitbucket.ListPRsOptions{
				State:   state,
				Page:    page,
				PageLen: pagelen,
				All:     all,
			})
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(prs)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tSTATE\tAUTHOR\tSOURCE\tDESTINATION")
			for _, pr := range prs {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
					pr.ID, pr.Title, pr.State, pr.Author.DisplayName,
					pr.Source.Branch.Name, pr.Destination.Branch.Name)
			}
			return w.Flush()
		},
	}
	prListCmd.Flags().String("state", "", "Filter by state (OPEN, MERGED, DECLINED, SUPERSEDED)")
	addBBPaginationFlags(prListCmd)
	bbPRCmd.AddCommand(prListCmd)

	// pr get
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "get [workspace] <repo-slug> <pr-id>",
		Short: "Get pull request details",
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

			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			pr, err := client.GetPullRequest(workspace, repoSlug, prID)
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(pr)
			}

			fmt.Printf("ID:           %d\n", pr.ID)
			fmt.Printf("Title:        %s\n", pr.Title)
			fmt.Printf("State:        %s\n", pr.State)
			fmt.Printf("Author:       %s\n", pr.Author.DisplayName)
			fmt.Printf("Source:       %s\n", pr.Source.Branch.Name)
			fmt.Printf("Destination:  %s\n", pr.Destination.Branch.Name)
			fmt.Printf("Comments:     %d\n", pr.CommentCount)
			fmt.Printf("Tasks:        %d\n", pr.TaskCount)
			fmt.Printf("Created:      %s\n", pr.CreatedOn)
			fmt.Printf("Updated:      %s\n", pr.UpdatedOn)
			fmt.Printf("URL:          %s\n", pr.Links.HTML.Href)
			if pr.Description != "" {
				fmt.Printf("\nDescription:\n%s\n", pr.Description)
			}
			return nil
		},
	})

	// pr create
	prCreateCmd := &cobra.Command{
		Use:   "create [workspace] <repo-slug>",
		Short: "Create a pull request",
		Long: `Create a pull request. Use --from-issue to auto-populate the title and
description from a Jira issue (e.g. --from-issue PROJ-123). The PR description
will include a link back to the Jira ticket.`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, err := resolveWorkspaceAndRepo(cmd, args)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			title, _ := cmd.Flags().GetString("title")
			source, _ := cmd.Flags().GetString("source")
			dest, _ := cmd.Flags().GetString("destination")
			desc, _ := cmd.Flags().GetString("description")
			closeBranch, _ := cmd.Flags().GetBool("close-source-branch")
			fromIssue, _ := cmd.Flags().GetString("from-issue")

			// If --from-issue is specified, fetch issue details and populate title/description
			if fromIssue != "" {
				jiraClient, err := getJiraClient(cmd)
				if err != nil {
					return fmt.Errorf("--from-issue requires a configured Jira profile: %w", err)
				}
				issue, err := jiraClient.GetIssue(fromIssue, []string{"summary", "description"}, nil)
				if err != nil {
					return fmt.Errorf("fetching Jira issue %s: %w", fromIssue, err)
				}

				if title == "" {
					title = fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary)
				}
				if desc == "" {
					profile, _ := getProfile(cmd)
					issueURL := fmt.Sprintf("%s/browse/%s",
						strings.TrimRight(profile.AtlassianURL, "/"), issue.Key)
					desc = fmt.Sprintf("Resolves [%s](%s)\n\n%s",
						issue.Key, issueURL, issue.Fields.Summary)
				}
				if source == "" {
					// Generate a branch name from the issue key
					source = strings.ToLower(fmt.Sprintf("feature/%s", issue.Key))
				}
			}

			if title == "" || source == "" {
				return fmt.Errorf("--title and --source are required (or use --from-issue to auto-populate)")
			}

			pr, err := client.CreatePullRequest(workspace, repoSlug, &bitbucket.CreatePRRequest{
				Title:             title,
				Description:       desc,
				SourceBranch:      source,
				DestinationBranch: dest,
				CloseSourceBranch: closeBranch,
			})
			if err != nil {
				return err
			}

			return outputResult(cmd, "created", fmt.Sprintf("%d", pr.ID), fmt.Sprintf("Created PR #%d: %s", pr.ID, pr.Title), pr)
		},
	}
	prCreateCmd.Flags().String("title", "", "Pull request title (required, or use --from-issue)")
	prCreateCmd.Flags().String("source", "", "Source branch name (required, or use --from-issue)")
	prCreateCmd.Flags().String("destination", "", "Destination branch name (defaults to main branch)")
	prCreateCmd.Flags().String("description", "", "Pull request description")
	prCreateCmd.Flags().Bool("close-source-branch", false, "Close source branch after merge")
	prCreateCmd.Flags().String("from-issue", "", "Jira issue key to populate title/description from (e.g. PROJ-123)")
	bbPRCmd.AddCommand(prCreateCmd)

	// pr approve
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "approve [workspace] <repo-slug> <pr-id>",
		Short: "Approve a pull request",
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
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.ApprovePullRequest(workspace, repoSlug, prID); err != nil {
				return err
			}
			return outputResult(cmd, "approved", fmt.Sprintf("%d", prID), fmt.Sprintf("Approved PR #%d", prID), nil)
		},
	})

	// pr unapprove
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "unapprove [workspace] <repo-slug> <pr-id>",
		Short: "Remove approval from a pull request",
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
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.UnapprovePullRequest(workspace, repoSlug, prID); err != nil {
				return err
			}
			return outputResult(cmd, "unapproved", fmt.Sprintf("%d", prID), fmt.Sprintf("Removed approval from PR #%d", prID), nil)
		},
	})

	// pr decline
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "decline [workspace] <repo-slug> <pr-id>",
		Short: "Decline a pull request",
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
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeclinePullRequest(workspace, repoSlug, prID); err != nil {
				return err
			}
			return outputResult(cmd, "declined", fmt.Sprintf("%d", prID), fmt.Sprintf("Declined PR #%d", prID), nil)
		},
	})

	// pr merge
	prMergeCmd := &cobra.Command{
		Use:   "merge [workspace] <repo-slug> <pr-id>",
		Short: "Merge a pull request",
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
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			strategy, _ := cmd.Flags().GetString("strategy")
			message, _ := cmd.Flags().GetString("message")
			closeBranch, _ := cmd.Flags().GetBool("close-source-branch")

			req := &bitbucket.MergePRRequest{
				MergeStrategy: strategy,
				Message:       message,
			}
			if cmd.Flags().Changed("close-source-branch") {
				req.CloseSourceBranch = &closeBranch
			}

			pr, err := client.MergePullRequest(workspace, repoSlug, prID, req)
			if err != nil {
				return err
			}
			return outputResult(cmd, "merged", fmt.Sprintf("%d", pr.ID), fmt.Sprintf("Merged PR #%d: %s", pr.ID, pr.Title), pr)
		},
	}
	prMergeCmd.Flags().String("strategy", "", "Merge strategy (merge_commit, squash, fast_forward)")
	prMergeCmd.Flags().String("message", "", "Merge commit message")
	prMergeCmd.Flags().Bool("close-source-branch", false, "Close source branch after merge")
	bbPRCmd.AddCommand(prMergeCmd)

	// pr request-changes
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "request-changes [workspace] <repo-slug> <pr-id>",
		Short: "Request changes on a pull request",
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
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.RequestChangesPullRequest(workspace, repoSlug, prID); err != nil {
				return err
			}
			return outputResult(cmd, "changes_requested", fmt.Sprintf("%d", prID), fmt.Sprintf("Requested changes on PR #%d", prID), nil)
		},
	})

	// pr comments
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "comments [workspace] <repo-slug> <pr-id>",
		Short: "List comments on a pull request",
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
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			comments, err := client.ListPRComments(workspace, repoSlug, prID)
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(comments)
			}

			for _, c := range comments {
				fmt.Printf("#%d by %s (%s)\n", c.ID, c.User.DisplayName, c.CreatedOn)
				if c.Inline != nil {
					fmt.Printf("  File: %s\n", c.Inline.Path)
				}
				fmt.Printf("  %s\n\n", c.Content.Raw)
			}
			return nil
		},
	})

	// pr comment (add a comment)
	prCommentCmd := &cobra.Command{
		Use:   "comment [workspace] <repo-slug> <pr-id>",
		Short: "Add a comment to a pull request",
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
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			body, _ := cmd.Flags().GetString("body")
			if body == "" {
				return fmt.Errorf("--body is required")
			}

			filePath, _ := cmd.Flags().GetString("file")
			line, _ := cmd.Flags().GetInt("line")

			if filePath != "" && line == 0 {
				return fmt.Errorf("--line is required when --file is specified")
			}
			if line != 0 && filePath == "" {
				return fmt.Errorf("--file is required when --line is specified")
			}

			var inline *bitbucket.InlineCommentParams
			if filePath != "" {
				inline = &bitbucket.InlineCommentParams{
					Path: filePath,
					To:   line,
				}
			}

			comment, err := client.CreatePRCommentInline(workspace, repoSlug, prID, body, inline)
			if err != nil {
				return err
			}
			return outputResult(cmd, "created", fmt.Sprintf("%d", comment.ID), fmt.Sprintf("Added comment #%d to PR #%d", comment.ID, prID), comment)
		},
	}
	prCommentCmd.Flags().String("body", "", "Comment body (required)")
	prCommentCmd.Flags().String("file", "", "File path for an inline comment")
	prCommentCmd.Flags().Int("line", 0, "Line number in the new version of the file (requires --file)")
	bbPRCmd.AddCommand(prCommentCmd)

	// pr diff
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "diff [workspace] <repo-slug> <pr-id>",
		Short: "Get the diff of a pull request",
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
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			diff, err := client.GetPRDiff(workspace, repoSlug, prID)
			if err != nil {
				return err
			}
			fmt.Print(diff)
			return nil
		},
	})
}
