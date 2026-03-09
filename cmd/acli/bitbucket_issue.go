package acli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbIssueCmd = &cobra.Command{
	Use:     "issue",
	Short:   "Manage repository issues (Bitbucket issue tracker)",
	Aliases: []string{"i"},
	RunE:    helpRunE,
}

func init() {
	// issue list
	issueListCmd := &cobra.Command{
		Use:     "list <workspace> <repo-slug>",
		Short:   "List issues",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			q, _ := cmd.Flags().GetString("query")
			sort, _ := cmd.Flags().GetString("sort")

			issues, err := client.ListIssues(args[0], args[1], &bitbucket.ListIssuesOptions{
				Q:    q,
				Sort: sort,
			})
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tSTATE\tPRIORITY\tKIND\tASSIGNEE")
			for _, issue := range issues {
				assignee := ""
				if issue.Assignee != nil {
					assignee = issue.Assignee.DisplayName
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
					issue.ID, issue.Title, issue.State, issue.Priority, issue.Kind, assignee)
			}
			return w.Flush()
		},
	}
	issueListCmd.Flags().String("query", "", "Filter issues (Bitbucket query syntax)")
	issueListCmd.Flags().String("sort", "", "Sort field (e.g. -priority)")
	bbIssueCmd.AddCommand(issueListCmd)

	// issue get
	bbIssueCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug> <issue-id>",
		Short: "Get issue details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			issueID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[2])
			}

			issue, err := client.GetIssue(args[0], args[1], issueID)
			if err != nil {
				return err
			}

			fmt.Printf("ID:        %d\n", issue.ID)
			fmt.Printf("Title:     %s\n", issue.Title)
			fmt.Printf("State:     %s\n", issue.State)
			fmt.Printf("Priority:  %s\n", issue.Priority)
			fmt.Printf("Kind:      %s\n", issue.Kind)
			fmt.Printf("Reporter:  %s\n", issue.Reporter.DisplayName)
			if issue.Assignee != nil {
				fmt.Printf("Assignee:  %s\n", issue.Assignee.DisplayName)
			}
			if issue.Component != nil {
				fmt.Printf("Component: %s\n", issue.Component.Name)
			}
			if issue.Milestone != nil {
				fmt.Printf("Milestone: %s\n", issue.Milestone.Name)
			}
			fmt.Printf("Votes:     %d\n", issue.Votes)
			fmt.Printf("Created:   %s\n", issue.CreatedOn)
			fmt.Printf("Updated:   %s\n", issue.UpdatedOn)
			fmt.Printf("URL:       %s\n", issue.Links.HTML.Href)
			if issue.Content.Raw != "" {
				fmt.Printf("\nContent:\n%s\n", issue.Content.Raw)
			}
			return nil
		},
	})

	// issue create
	issueCreateCmd := &cobra.Command{
		Use:   "create <workspace> <repo-slug>",
		Short: "Create an issue",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			title, _ := cmd.Flags().GetString("title")
			content, _ := cmd.Flags().GetString("content")
			kind, _ := cmd.Flags().GetString("kind")
			priority, _ := cmd.Flags().GetString("priority")

			if title == "" {
				return fmt.Errorf("--title is required")
			}

			req := &bitbucket.CreateIssueRequest{
				Title:    title,
				Kind:     kind,
				Priority: priority,
			}
			if content != "" {
				req.Content = &struct {
					Raw string `json:"raw"`
				}{Raw: content}
			}

			issue, err := client.CreateIssue(args[0], args[1], req)
			if err != nil {
				return err
			}

			fmt.Printf("Created issue #%d: %s\n", issue.ID, issue.Title)
			fmt.Printf("URL: %s\n", issue.Links.HTML.Href)
			return nil
		},
	}
	issueCreateCmd.Flags().String("title", "", "Issue title (required)")
	issueCreateCmd.Flags().String("content", "", "Issue content/description")
	issueCreateCmd.Flags().String("kind", "bug", "Issue kind (bug, enhancement, proposal, task)")
	issueCreateCmd.Flags().String("priority", "major", "Issue priority (trivial, minor, major, critical, blocker)")
	bbIssueCmd.AddCommand(issueCreateCmd)

	// issue update
	issueUpdateCmd := &cobra.Command{
		Use:   "update <workspace> <repo-slug> <issue-id>",
		Short: "Update an issue",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			issueID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[2])
			}

			title, _ := cmd.Flags().GetString("title")
			state, _ := cmd.Flags().GetString("state")
			kind, _ := cmd.Flags().GetString("kind")
			priority, _ := cmd.Flags().GetString("priority")

			req := &bitbucket.UpdateIssueRequest{
				Title:    title,
				State:    state,
				Kind:     kind,
				Priority: priority,
			}

			issue, err := client.UpdateIssue(args[0], args[1], issueID, req)
			if err != nil {
				return err
			}

			fmt.Printf("Updated issue #%d: %s\n", issue.ID, issue.Title)
			return nil
		},
	}
	issueUpdateCmd.Flags().String("title", "", "New title")
	issueUpdateCmd.Flags().String("state", "", "New state (new, open, resolved, on hold, invalid, duplicate, wontfix, closed)")
	issueUpdateCmd.Flags().String("kind", "", "New kind")
	issueUpdateCmd.Flags().String("priority", "", "New priority")
	bbIssueCmd.AddCommand(issueUpdateCmd)

	// issue delete
	bbIssueCmd.AddCommand(&cobra.Command{
		Use:   "delete <workspace> <repo-slug> <issue-id>",
		Short: "Delete an issue",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			issueID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[2])
			}
			if err := client.DeleteIssue(args[0], args[1], issueID); err != nil {
				return err
			}
			fmt.Printf("Deleted issue #%d\n", issueID)
			return nil
		},
	})

	// issue comments
	bbIssueCmd.AddCommand(&cobra.Command{
		Use:   "comments <workspace> <repo-slug> <issue-id>",
		Short: "List comments on an issue",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			issueID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[2])
			}

			comments, err := client.ListIssueComments(args[0], args[1], issueID)
			if err != nil {
				return err
			}

			for _, c := range comments {
				fmt.Printf("#%d by %s (%s)\n", c.ID, c.User.DisplayName, c.CreatedOn)
				fmt.Printf("  %s\n\n", c.Content.Raw)
			}
			return nil
		},
	})

	// issue comment (add)
	issueCommentCmd := &cobra.Command{
		Use:   "comment <workspace> <repo-slug> <issue-id>",
		Short: "Add a comment to an issue",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			issueID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[2])
			}

			body, _ := cmd.Flags().GetString("body")
			if body == "" {
				return fmt.Errorf("--body is required")
			}

			comment, err := client.CreateIssueComment(args[0], args[1], issueID, body)
			if err != nil {
				return err
			}
			fmt.Printf("Added comment #%d to issue #%d\n", comment.ID, issueID)
			return nil
		},
	}
	issueCommentCmd.Flags().String("body", "", "Comment body (required)")
	bbIssueCmd.AddCommand(issueCommentCmd)
}
