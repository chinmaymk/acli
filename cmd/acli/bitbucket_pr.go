package acli

import (
	"fmt"
	"io"
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
			author, _ := cmd.Flags().GetString("author")
			pOpts := getBBPaginationOpts(cmd)
			prs, err := client.ListPullRequests(workspace, repoSlug, &bitbucket.ListPRsOptions{
				State:   state,
				Author:  author,
				Page:    pOpts.Page,
				PageLen: pOpts.PageLen,
				All:     pOpts.All,
			})
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(prs)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "ID\tTITLE\tSTATE\tAUTHOR\tSOURCE\tDESTINATION")
			for _, pr := range prs {
				_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
					pr.ID, pr.Title, pr.State, pr.Author.DisplayName,
					pr.Source.Branch.Name, pr.Destination.Branch.Name)
			}
			return w.Flush()
		},
	}
	prListCmd.Flags().String("state", "OPEN", "Filter by state (OPEN, MERGED, DECLINED, SUPERSEDED)")
	prListCmd.Flags().String("author", "", "Filter by author Bitbucket username or UUID (e.g. 'jdoe' or '{d301aafa-d676-4ee0-88be-962be7417567}')")
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
		Args:  cobra.RangeArgs(1, 2),
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

			if title == "" || source == "" {
				return fmt.Errorf("--title and --source are required")
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
	prCreateCmd.Flags().String("title", "", "Pull request title (required)")
	prCreateCmd.Flags().String("source", "", "Source branch name (required)")
	prCreateCmd.Flags().String("destination", "", "Destination branch name (defaults to main branch)")
	prCreateCmd.Flags().String("description", "", "Pull request description")
	prCreateCmd.Flags().Bool("close-source-branch", false, "Close source branch after merge")
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
			if _, err := client.ApprovePullRequest(workspace, repoSlug, prID); err != nil {
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
			if _, err := client.DeclinePullRequest(workspace, repoSlug, prID); err != nil {
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
			if _, err := client.RequestChangesPullRequest(workspace, repoSlug, prID); err != nil {
				return err
			}
			return outputResult(cmd, "changes_requested", fmt.Sprintf("%d", prID), fmt.Sprintf("Requested changes on PR #%d", prID), nil)
		},
	})

	// pr comments
	prCommentsCmd := &cobra.Command{
		Use:   "comments [workspace] <repo-slug> <pr-id>",
		Short: "List comments on a pull request",
		Long: `List comments on a pull request.

Each comment is shown with its ID, author, file/line anchor (for inline
comments), and parent comment ID (for replies). Use 'pr review' to see
comments grouped into threads with their linked tasks.`,
		Args: cobra.RangeArgs(2, 3),
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

			comments, err := client.ListPRComments(workspace, repoSlug, prID, getBBPaginationOpts(cmd))
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(comments)
			}

			for _, c := range comments {
				header := fmt.Sprintf("#%d by %s (%s)", c.ID, c.User.DisplayName, c.CreatedOn)
				if c.Parent != nil {
					header += fmt.Sprintf(" reply-to:#%d", c.Parent.ID)
				}
				if c.Resolution != nil {
					header += " [RESOLVED]"
				}
				fmt.Println(header)
				if c.Inline != nil {
					fmt.Printf("  Location: %s\n", commentLocation(&c))
				}
				fmt.Printf("  %s\n\n", c.Content.Raw)
			}
			return nil
		},
	}
	addBBPaginationFlags(prCommentsCmd)
	bbPRCmd.AddCommand(prCommentsCmd)

	// pr review - joined view of comments + tasks
	prReviewCmd := &cobra.Command{
		Use:   "review [workspace] <repo-slug> <pr-id>",
		Short: "Show all PR feedback (comments + tasks) grouped into threads",
		Long: `Fetches all comments and tasks on a pull request in one call and joins
them into a structured review:

  - File threads: top-level inline comments grouped by file and line, with
    their replies and any tasks linked to the thread.
  - General threads: top-level non-inline comments, with replies and tasks.
  - Standalone tasks: tasks not linked to any comment.

This is the recommended command for coding agents that need the full picture
of reviewer feedback with correct associations between comments, replies, and
tasks. Use --output json for a machine-readable payload.`,
		Args: cobra.RangeArgs(2, 3),
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

			review, err := client.GetPRReview(workspace, repoSlug, prID)
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(review)
			}

			printPRReview(cmd, review)
			return nil
		},
	}
	bbPRCmd.AddCommand(prReviewCmd)

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

	// pr commits
	prCommitsCmd := &cobra.Command{
		Use:     "commits [workspace] <repo-slug> <pr-id>",
		Short:   "List commits on a pull request",
		Aliases: []string{"log"},
		Args:    cobra.RangeArgs(2, 3),
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

			commits, err := client.ListPRCommits(workspace, repoSlug, prID, getBBPaginationOpts(cmd))
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(commits)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "HASH\tMESSAGE\tAUTHOR\tDATE")
			for _, c := range commits {
				hash := c.Hash
				if len(hash) > 12 {
					hash = hash[:12]
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					hash, truncate(firstLine(c.Message), 60), c.Author.Raw, c.Date)
			}
			return w.Flush()
		},
	}
	addBBPaginationFlags(prCommitsCmd)
	bbPRCmd.AddCommand(prCommitsCmd)

	// pr diffstat
	prDiffstatCmd := &cobra.Command{
		Use:   "diffstat [workspace] <repo-slug> <pr-id>",
		Short: "Get file-level diff statistics for a pull request",
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

			stats, err := client.GetPRDiffStat(workspace, repoSlug, prID, getBBPaginationOpts(cmd))
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(stats)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "STATUS\tFILE\tADDED\tREMOVED")
			for _, s := range stats {
				filePath := ""
				if s.Status == "renamed" && s.Old != nil && s.New != nil {
					filePath = s.Old.Path + " → " + s.New.Path
				} else if s.New != nil {
					filePath = s.New.Path
				} else if s.Old != nil {
					filePath = s.Old.Path
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t+%d\t-%d\n",
					s.Status, filePath, s.LinesAdded, s.LinesRemoved)
			}
			return w.Flush()
		},
	}
	addBBPaginationFlags(prDiffstatCmd)
	bbPRCmd.AddCommand(prDiffstatCmd)

	// pr activity
	prActivityCmd := &cobra.Command{
		Use:   "activity [workspace] <repo-slug> <pr-id>",
		Short: "List activity on a pull request (comments, approvals, updates)",
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

			activities, err := client.ListPRActivity(workspace, repoSlug, prID, getBBPaginationOpts(cmd))
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(activities)
			}

			for _, a := range activities {
				switch {
				case a.Approval != nil:
					fmt.Printf("[APPROVAL] %s approved on %s\n", a.Approval.User.DisplayName, a.Approval.Date)
				case a.Update != nil:
					fmt.Printf("[UPDATE] %s changed state to %s on %s\n",
						a.Update.Author.DisplayName, a.Update.State, a.Update.Date)
					if a.Update.Reason != "" {
						fmt.Printf("  Reason: %s\n", a.Update.Reason)
					}
				case a.Comment != nil:
					fmt.Printf("[COMMENT] #%d by %s (%s)\n", a.Comment.ID, a.Comment.User.DisplayName, a.Comment.CreatedOn)
					if a.Comment.Inline != nil {
						fmt.Printf("  File: %s\n", a.Comment.Inline.Path)
					}
					fmt.Printf("  %s\n", truncate(a.Comment.Content.Raw, 120))
				default:
					fmt.Println("[UNKNOWN] unrecognized activity type")
				}
				fmt.Println()
			}
			return nil
		},
	}
	addBBPaginationFlags(prActivityCmd)
	bbPRCmd.AddCommand(prActivityCmd)
}

// commentLocation returns a human-readable location string for an inline
// comment (e.g. "path/to/file.go:42"). Multi-line comments render as
// "path:start-end". Returns an empty string if the comment is not inline.
func commentLocation(c *bitbucket.PRComment) string {
	if c == nil || c.Inline == nil || c.Inline.Path == "" {
		return ""
	}
	line := 0
	if c.Inline.To != nil {
		line = *c.Inline.To
	}
	start := 0
	if c.Inline.StartTo != nil {
		start = *c.Inline.StartTo
	}
	switch {
	case start > 0 && line > 0 && start != line:
		return fmt.Sprintf("%s:%d-%d", c.Inline.Path, start, line)
	case line > 0:
		return fmt.Sprintf("%s:%d", c.Inline.Path, line)
	default:
		return c.Inline.Path
	}
}

// printPRReview prints a human-readable rendering of a PRReview.
func printPRReview(cmd *cobra.Command, review *bitbucket.PRReview) {
	out := cmd.OutOrStdout()

	_, _ = fmt.Fprintf(out, "PR #%d review: %d threads, %d comments, %d tasks (%d unresolved)\n\n",
		review.PullRequestID,
		review.Counts.Threads,
		review.Counts.Comments,
		review.Counts.Tasks,
		review.Counts.UnresolvedTasks,
	)

	if len(review.FileThreads) > 0 {
		_, _ = fmt.Fprintln(out, "== File threads ==")
		for i := range review.FileThreads {
			printReviewThread(out, &review.FileThreads[i])
		}
	}

	if len(review.GeneralThreads) > 0 {
		_, _ = fmt.Fprintln(out, "== General threads ==")
		for i := range review.GeneralThreads {
			printReviewThread(out, &review.GeneralThreads[i])
		}
	}

	if len(review.StandaloneTasks) > 0 {
		_, _ = fmt.Fprintln(out, "== Standalone tasks ==")
		for _, t := range review.StandaloneTasks {
			_, _ = fmt.Fprintf(out, "  [task #%d %s] %s — %s\n",
				t.ID, t.State, t.Creator.DisplayName, t.Content.Raw)
		}
		_, _ = fmt.Fprintln(out)
	}
}

func printReviewThread(out io.Writer, t *bitbucket.PRReviewThread) {
	loc := commentLocation(&t.Comment)
	if loc == "" {
		loc = "(general)"
	}
	resolution := ""
	if t.Comment.Resolution != nil {
		resolution = " [RESOLVED]"
	}
	_, _ = fmt.Fprintf(out, "\n-- %s --%s\n", loc, resolution)
	_, _ = fmt.Fprintf(out, "  #%d by %s (%s)\n", t.Comment.ID, t.Comment.User.DisplayName, t.Comment.CreatedOn)
	_, _ = fmt.Fprintf(out, "    %s\n", indent(t.Comment.Content.Raw, "    "))

	for _, r := range t.Replies {
		_, _ = fmt.Fprintf(out, "  ↳ #%d by %s (%s)\n", r.ID, r.User.DisplayName, r.CreatedOn)
		_, _ = fmt.Fprintf(out, "    %s\n", indent(r.Content.Raw, "    "))
	}

	for _, task := range t.Tasks {
		_, _ = fmt.Fprintf(out, "  ☑ task #%d [%s] by %s: %s\n",
			task.ID, task.State, task.Creator.DisplayName, task.Content.Raw)
	}
}

// indent prepends `prefix` to each subsequent line of s.
func indent(s, prefix string) string {
	if s == "" {
		return s
	}
	// Replace each newline with newline+prefix so every continuation line is aligned.
	out := ""
	first := true
	for _, line := range strings.Split(s, "\n") {
		if first {
			out = line
			first = false
			continue
		}
		out += "\n" + prefix + line
	}
	return out
}
