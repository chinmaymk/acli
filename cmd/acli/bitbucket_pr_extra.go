package acli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

func init() {
	// pr update
	prUpdateCmd := &cobra.Command{
		Use:   "update [workspace] <repo-slug> <pr-id>",
		Short: "Update a pull request's title, description, or reviewers",
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

			req := &bitbucket.UpdatePRRequest{}
			changed := false

			if cmd.Flags().Changed("title") {
				req.Title, _ = cmd.Flags().GetString("title")
				changed = true
			}
			if cmd.Flags().Changed("description") {
				req.Description, _ = cmd.Flags().GetString("description")
				changed = true
			}
			if cmd.Flags().Changed("close-source-branch") {
				v, _ := cmd.Flags().GetBool("close-source-branch")
				req.CloseSourceBranch = &v
				changed = true
			}

			if !changed {
				return fmt.Errorf("at least one of --title, --description, or --close-source-branch is required")
			}

			pr, err := client.UpdatePullRequest(workspace, repoSlug, prID, req)
			if err != nil {
				return err
			}
			return outputResult(cmd, "updated", fmt.Sprintf("%d", pr.ID), fmt.Sprintf("Updated PR #%d: %s", pr.ID, pr.Title), pr)
		},
	}
	prUpdateCmd.Flags().String("title", "", "New pull request title")
	prUpdateCmd.Flags().String("description", "", "New pull request description")
	prUpdateCmd.Flags().Bool("close-source-branch", false, "Close source branch after merge")
	bbPRCmd.AddCommand(prUpdateCmd)

	// pr reviewers (list reviewers)
	bbPRCmd.AddCommand(&cobra.Command{
		Use:   "reviewers [workspace] <repo-slug> <pr-id>",
		Short: "List reviewers on a pull request",
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
				return outputJSON(pr.Participants)
			}

			if len(pr.Participants) == 0 && len(pr.Reviewers) == 0 {
				fmt.Println("No reviewers on this pull request.")
				return nil
			}

			// Show participants with their review status
			for _, p := range pr.Participants {
				status := p.State
				if p.Approved {
					status = "APPROVED"
				}
				role := strings.ToLower(p.Role)
				fmt.Printf("%s  (%s, %s)\n", p.User.DisplayName, role, status)
			}
			return nil
		},
	})

	// pr add-reviewer
	prAddReviewerCmd := &cobra.Command{
		Use:   "add-reviewer [workspace] <repo-slug> <pr-id>",
		Short: "Add reviewers to a pull request",
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

			uuids, _ := cmd.Flags().GetStringSlice("uuid")
			if len(uuids) == 0 {
				return fmt.Errorf("--uuid is required (can be specified multiple times)")
			}

			pr, err := client.AddReviewersToPullRequest(workspace, repoSlug, prID, uuids)
			if err != nil {
				return err
			}
			return outputResult(cmd, "updated", fmt.Sprintf("%d", pr.ID),
				fmt.Sprintf("Added %d reviewer(s) to PR #%d", len(uuids), prID), pr)
		},
	}
	prAddReviewerCmd.Flags().StringSlice("uuid", nil, "Reviewer UUID (can be repeated)")
	bbPRCmd.AddCommand(prAddReviewerCmd)

	// pr remove-reviewer
	prRemoveReviewerCmd := &cobra.Command{
		Use:   "remove-reviewer [workspace] <repo-slug> <pr-id>",
		Short: "Remove reviewers from a pull request",
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

			uuids, _ := cmd.Flags().GetStringSlice("uuid")
			if len(uuids) == 0 {
				return fmt.Errorf("--uuid is required (can be repeated)")
			}

			pr, err := client.RemoveReviewersFromPullRequest(workspace, repoSlug, prID, uuids)
			if err != nil {
				return err
			}
			return outputResult(cmd, "updated", fmt.Sprintf("%d", pr.ID),
				fmt.Sprintf("Removed %d reviewer(s) from PR #%d", len(uuids), prID), pr)
		},
	}
	prRemoveReviewerCmd.Flags().StringSlice("uuid", nil, "Reviewer UUID to remove (can be repeated)")
	bbPRCmd.AddCommand(prRemoveReviewerCmd)
}
