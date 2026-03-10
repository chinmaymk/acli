package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var bbDownloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Manage repository downloads",
	Aliases: []string{"dl"},
	RunE:    helpRunE,
}

func init() {
	// download list
	bbDownloadCmd.AddCommand(&cobra.Command{
		Use:     "list [workspace] <repo-slug>",
		Short:   "List downloads for a repository",
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

			downloads, err := client.ListDownloads(workspace, repoSlug)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSIZE\tDOWNLOADS\tCREATED")
			for _, d := range downloads {
				fmt.Fprintf(w, "%s\t%d\t%d\t%s\n",
					d.Name, d.Size, d.Downloads, d.CreatedOn)
			}
			return w.Flush()
		},
	})

	// download delete
	bbDownloadCmd.AddCommand(&cobra.Command{
		Use:   "delete [workspace] <repo-slug> <filename>",
		Short: "Delete a download artifact",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, filename, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeleteDownload(workspace, repoSlug, filename); err != nil {
				return err
			}
			fmt.Printf("Deleted download: %s\n", filename)
			return nil
		},
	})
}
