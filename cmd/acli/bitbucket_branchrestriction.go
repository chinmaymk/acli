package acli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbBranchRestrictionCmd = &cobra.Command{
	Use:     "branch-restriction",
	Short:   "Manage branch restrictions",
	Aliases: []string{"restriction"},
	RunE:    helpRunE,
}

func init() {
	// branch-restriction list
	bbBranchRestrictionCmd.AddCommand(&cobra.Command{
		Use:     "list [workspace] <repo-slug>",
		Short:   "List branch restrictions",
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

			restrictions, err := client.ListBranchRestrictions(workspace, repoSlug)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tKIND\tPATTERN")
			for _, r := range restrictions {
				fmt.Fprintf(w, "%d\t%s\t%s\n", r.ID, r.Kind, r.Pattern)
			}
			return w.Flush()
		},
	})

	// branch-restriction get
	bbBranchRestrictionCmd.AddCommand(&cobra.Command{
		Use:   "get [workspace] <repo-slug> <restriction-id>",
		Short: "Get branch restriction details",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, idStr, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid restriction ID: %s", idStr)
			}

			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			r, err := client.GetBranchRestriction(workspace, repoSlug, id)
			if err != nil {
				return err
			}

			fmt.Printf("ID:       %d\n", r.ID)
			fmt.Printf("Kind:     %s\n", r.Kind)
			fmt.Printf("Pattern:  %s\n", r.Pattern)
			if r.Value != nil {
				fmt.Printf("Value:    %d\n", *r.Value)
			}
			if len(r.Users) > 0 {
				fmt.Println("Users:")
				for _, u := range r.Users {
					fmt.Printf("  - %s\n", u.DisplayName)
				}
			}
			if len(r.Groups) > 0 {
				fmt.Println("Groups:")
				for _, g := range r.Groups {
					fmt.Printf("  - %s (%s)\n", g.Name, g.Slug)
				}
			}
			return nil
		},
	})

	// branch-restriction create
	brCreateCmd := &cobra.Command{
		Use:   "create [workspace] <repo-slug>",
		Short: "Create a branch restriction",
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

			kind, _ := cmd.Flags().GetString("kind")
			pattern, _ := cmd.Flags().GetString("pattern")

			if kind == "" || pattern == "" {
				return fmt.Errorf("--kind and --pattern are required")
			}

			r, err := client.CreateBranchRestriction(workspace, repoSlug, &bitbucket.CreateBranchRestrictionRequest{
				Kind:    kind,
				Pattern: pattern,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Created branch restriction: %d (%s on %s)\n", r.ID, r.Kind, r.Pattern)
			return nil
		},
	}
	brCreateCmd.Flags().String("kind", "", "Restriction kind (e.g. push, force, delete, restrict_merges, require_approvals_to_merge)")
	brCreateCmd.Flags().String("pattern", "", "Branch pattern (e.g. main, release/*)")
	bbBranchRestrictionCmd.AddCommand(brCreateCmd)

	// branch-restriction delete
	bbBranchRestrictionCmd.AddCommand(&cobra.Command{
		Use:   "delete [workspace] <repo-slug> <restriction-id>",
		Short: "Delete a branch restriction",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, idStr, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid restriction ID: %s", idStr)
			}

			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			if err := client.DeleteBranchRestriction(workspace, repoSlug, id); err != nil {
				return err
			}
			fmt.Printf("Deleted branch restriction: %d\n", id)
			return nil
		},
	})
}
