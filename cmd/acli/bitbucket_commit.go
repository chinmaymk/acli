package acli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbCommitCmd = &cobra.Command{
	Use:     "commit",
	Short:   "Manage commits",
	Aliases: []string{"cm"},
	RunE:    helpRunE,
}

func init() {
	// commit list
	commitListCmd := &cobra.Command{
		Use:     "list <workspace> <repo-slug>",
		Short:   "List commits",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			include, _ := cmd.Flags().GetString("include")
			exclude, _ := cmd.Flags().GetString("exclude")

			commits, err := client.ListCommits(args[0], args[1], include, exclude)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "HASH\tDATE\tAUTHOR\tMESSAGE")
			for _, c := range commits {
				hash := c.Hash
				if len(hash) > 12 {
					hash = hash[:12]
				}
				msg := strings.Split(c.Message, "\n")[0]
				if len(msg) > 60 {
					msg = msg[:57] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					hash, c.Date, c.Author.Raw, msg)
			}
			return w.Flush()
		},
	}
	commitListCmd.Flags().String("include", "", "Include commits reachable from this ref")
	commitListCmd.Flags().String("exclude", "", "Exclude commits reachable from this ref")
	bbCommitCmd.AddCommand(commitListCmd)

	// commit get
	bbCommitCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug> <commit-hash>",
		Short: "Get commit details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			commit, err := client.GetCommit(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			fmt.Printf("Hash:     %s\n", commit.Hash)
			fmt.Printf("Date:     %s\n", commit.Date)
			fmt.Printf("Author:   %s\n", commit.Author.Raw)
			fmt.Printf("Message:  %s\n", commit.Message)
			if len(commit.Parents) > 0 {
				var parents []string
				for _, p := range commit.Parents {
					parents = append(parents, p.Hash[:12])
				}
				fmt.Printf("Parents:  %s\n", strings.Join(parents, ", "))
			}
			fmt.Printf("URL:      %s\n", commit.Links.HTML.Href)
			return nil
		},
	})

	// commit statuses
	bbCommitCmd.AddCommand(&cobra.Command{
		Use:   "statuses <workspace> <repo-slug> <commit-hash>",
		Short: "List build statuses for a commit",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			statuses, err := client.ListCommitStatuses(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "KEY\tSTATE\tNAME\tDESCRIPTION\tUPDATED")
			for _, s := range statuses {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					s.Key, s.State, s.Name, s.Description, s.UpdatedOn)
			}
			return w.Flush()
		},
	})

	// commit diff
	bbCommitCmd.AddCommand(&cobra.Command{
		Use:   "diff <workspace> <repo-slug> <spec>",
		Short: "Get diff between two commits (e.g. commit1..commit2)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			diff, err := client.GetDiff(args[0], args[1], args[2])
			if err != nil {
				return err
			}
			fmt.Print(diff)
			return nil
		},
	})
}
