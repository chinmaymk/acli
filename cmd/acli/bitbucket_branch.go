package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbBranchCmd = &cobra.Command{
	Use:     "branch",
	Short:   "Manage branches",
	Aliases: []string{"br"},
	RunE:    helpRunE,
}

var bbTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
	RunE:  helpRunE,
}

func init() {
	// branch list
	branchListCmd := &cobra.Command{
		Use:     "list <workspace> <repo-slug>",
		Short:   "List branches",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			q, _ := cmd.Flags().GetString("query")
			branches, err := client.ListBranches(args[0], args[1], q)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tHASH\tDATE\tAUTHOR")
			for _, b := range branches {
				hash := b.Target.Hash
				if len(hash) > 12 {
					hash = hash[:12]
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					b.Name, hash, b.Target.Date, b.Target.Author.Raw)
			}
			return w.Flush()
		},
	}
	branchListCmd.Flags().String("query", "", "Filter branches (e.g. name ~ \"feature\")")
	bbBranchCmd.AddCommand(branchListCmd)

	// branch get
	bbBranchCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug> <branch-name>",
		Short: "Get branch details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			branch, err := client.GetBranch(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			fmt.Printf("Name:     %s\n", branch.Name)
			fmt.Printf("Hash:     %s\n", branch.Target.Hash)
			fmt.Printf("Date:     %s\n", branch.Target.Date)
			fmt.Printf("Author:   %s\n", branch.Target.Author.Raw)
			fmt.Printf("Message:  %s\n", branch.Target.Message)
			return nil
		},
	})

	// branch create
	branchCreateCmd := &cobra.Command{
		Use:   "create <workspace> <repo-slug>",
		Short: "Create a branch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			target, _ := cmd.Flags().GetString("target")
			if name == "" || target == "" {
				return fmt.Errorf("--name and --target are required")
			}

			req := &bitbucket.CreateBranchRequest{
				Name: name,
			}
			req.Target.Hash = target

			branch, err := client.CreateBranch(args[0], args[1], req)
			if err != nil {
				return err
			}
			fmt.Printf("Created branch: %s (%s)\n", branch.Name, branch.Target.Hash)
			return nil
		},
	}
	branchCreateCmd.Flags().String("name", "", "Branch name (required)")
	branchCreateCmd.Flags().String("target", "", "Target commit hash (required)")
	bbBranchCmd.AddCommand(branchCreateCmd)

	// branch delete
	bbBranchCmd.AddCommand(&cobra.Command{
		Use:   "delete <workspace> <repo-slug> <branch-name>",
		Short: "Delete a branch",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			if err := client.DeleteBranch(args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("Deleted branch: %s\n", args[2])
			return nil
		},
	})

	// tag list
	tagListCmd := &cobra.Command{
		Use:     "list <workspace> <repo-slug>",
		Short:   "List tags",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			q, _ := cmd.Flags().GetString("query")
			tags, err := client.ListTags(args[0], args[1], q)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tHASH\tDATE\tMESSAGE")
			for _, t := range tags {
				hash := t.Target.Hash
				if len(hash) > 12 {
					hash = hash[:12]
				}
				msg := t.Message
				if len(msg) > 60 {
					msg = msg[:57] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					t.Name, hash, t.Target.Date, msg)
			}
			return w.Flush()
		},
	}
	tagListCmd.Flags().String("query", "", "Filter tags (e.g. name ~ \"v1\")")
	bbTagCmd.AddCommand(tagListCmd)

	// tag get
	bbTagCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug> <tag-name>",
		Short: "Get tag details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			tag, err := client.GetTag(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			fmt.Printf("Name:     %s\n", tag.Name)
			fmt.Printf("Hash:     %s\n", tag.Target.Hash)
			fmt.Printf("Date:     %s\n", tag.Target.Date)
			if tag.Message != "" {
				fmt.Printf("Message:  %s\n", tag.Message)
			}
			return nil
		},
	})

	// tag create
	tagCreateCmd := &cobra.Command{
		Use:   "create <workspace> <repo-slug>",
		Short: "Create a tag",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			target, _ := cmd.Flags().GetString("target")
			message, _ := cmd.Flags().GetString("message")
			if name == "" || target == "" {
				return fmt.Errorf("--name and --target are required")
			}

			req := &bitbucket.CreateTagRequest{
				Name:    name,
				Message: message,
			}
			req.Target.Hash = target

			tag, err := client.CreateTag(args[0], args[1], req)
			if err != nil {
				return err
			}
			fmt.Printf("Created tag: %s (%s)\n", tag.Name, tag.Target.Hash)
			return nil
		},
	}
	tagCreateCmd.Flags().String("name", "", "Tag name (required)")
	tagCreateCmd.Flags().String("target", "", "Target commit hash (required)")
	tagCreateCmd.Flags().String("message", "", "Tag message")
	bbTagCmd.AddCommand(tagCreateCmd)

	// tag delete
	bbTagCmd.AddCommand(&cobra.Command{
		Use:   "delete <workspace> <repo-slug> <tag-name>",
		Short: "Delete a tag",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			if err := client.DeleteTag(args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("Deleted tag: %s\n", args[2])
			return nil
		},
	})
}
