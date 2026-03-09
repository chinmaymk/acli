package acli

import (
	"fmt"
	"os"
	"strconv"
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
		Use:     "list <workspace> <repo-slug>",
		Short:   "List pull requests",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			state, _ := cmd.Flags().GetString("state")
			prs, err := client.ListPullRequests(args[0], args[1], &bitbucket.ListPRsOptions{
				State: state,
			})
			if err != nil {
				return err
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
	bbPRCmd.AddCommand(prListCmd)

	// pr get
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug> <pr-id>",
		Short: "Get pull request details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
			}

			pr, err := client.GetPullRequest(args[0], args[1], prID)
			if err != nil {
				return err
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
		Use:   "create <workspace> <repo-slug>",
		Short: "Create a pull request",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
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

			pr, err := client.CreatePullRequest(args[0], args[1], &bitbucket.CreatePRRequest{
				Title:             title,
				Description:       desc,
				SourceBranch:      source,
				DestinationBranch: dest,
				CloseSourceBranch: closeBranch,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Created PR #%d: %s\n", pr.ID, pr.Title)
			fmt.Printf("URL: %s\n", pr.Links.HTML.Href)
			return nil
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
		Use:   "approve <workspace> <repo-slug> <pr-id>",
		Short: "Approve a pull request",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
			}
			if err := client.ApprovePullRequest(args[0], args[1], prID); err != nil {
				return err
			}
			fmt.Printf("Approved PR #%d\n", prID)
			return nil
		},
	})

	// pr unapprove
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "unapprove <workspace> <repo-slug> <pr-id>",
		Short: "Remove approval from a pull request",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
			}
			if err := client.UnapprovePullRequest(args[0], args[1], prID); err != nil {
				return err
			}
			fmt.Printf("Removed approval from PR #%d\n", prID)
			return nil
		},
	})

	// pr decline
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "decline <workspace> <repo-slug> <pr-id>",
		Short: "Decline a pull request",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
			}
			if err := client.DeclinePullRequest(args[0], args[1], prID); err != nil {
				return err
			}
			fmt.Printf("Declined PR #%d\n", prID)
			return nil
		},
	})

	// pr merge
	prMergeCmd := &cobra.Command{
		Use:   "merge <workspace> <repo-slug> <pr-id>",
		Short: "Merge a pull request",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
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

			pr, err := client.MergePullRequest(args[0], args[1], prID, req)
			if err != nil {
				return err
			}
			fmt.Printf("Merged PR #%d: %s\n", pr.ID, pr.Title)
			return nil
		},
	}
	prMergeCmd.Flags().String("strategy", "", "Merge strategy (merge_commit, squash, fast_forward)")
	prMergeCmd.Flags().String("message", "", "Merge commit message")
	prMergeCmd.Flags().Bool("close-source-branch", false, "Close source branch after merge")
	bbPRCmd.AddCommand(prMergeCmd)

	// pr request-changes
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "request-changes <workspace> <repo-slug> <pr-id>",
		Short: "Request changes on a pull request",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
			}
			if err := client.RequestChangesPullRequest(args[0], args[1], prID); err != nil {
				return err
			}
			fmt.Printf("Requested changes on PR #%d\n", prID)
			return nil
		},
	})

	// pr comments
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "comments <workspace> <repo-slug> <pr-id>",
		Short: "List comments on a pull request",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
			}

			comments, err := client.ListPRComments(args[0], args[1], prID)
			if err != nil {
				return err
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
		Use:   "comment <workspace> <repo-slug> <pr-id>",
		Short: "Add a comment to a pull request",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
			}

			body, _ := cmd.Flags().GetString("body")
			if body == "" {
				return fmt.Errorf("--body is required")
			}

			comment, err := client.CreatePRComment(args[0], args[1], prID, body)
			if err != nil {
				return err
			}
			fmt.Printf("Added comment #%d to PR #%d\n", comment.ID, prID)
			return nil
		},
	}
	prCommentCmd.Flags().String("body", "", "Comment body (required)")
	bbPRCmd.AddCommand(prCommentCmd)

	// pr diff
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "diff <workspace> <repo-slug> <pr-id>",
		Short: "Get the diff of a pull request",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", args[2])
			}

			diff, err := client.GetPRDiff(args[0], args[1], prID)
			if err != nil {
				return err
			}
			fmt.Print(diff)
			return nil
		},
	})
}
