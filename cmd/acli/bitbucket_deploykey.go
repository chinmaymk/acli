package acli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbDeployKeyCmd = &cobra.Command{
	Use:     "deploy-key",
	Short:   "Manage deploy keys",
	Aliases: []string{"dk"},
	RunE:    helpRunE,
}

func init() {
	// deploy-key list
	bbDeployKeyCmd.AddCommand(&cobra.Command{
		Use:     "list [workspace] <repo-slug>",
		Short:   "List deploy keys",
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

			keys, err := client.ListDeployKeys(workspace, repoSlug)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tLABEL\tCOMMENT\tCREATED")
			for _, k := range keys {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
					k.ID, k.Label, k.Comment, k.CreatedOn)
			}
			return w.Flush()
		},
	})

	// deploy-key get
	bbDeployKeyCmd.AddCommand(&cobra.Command{
		Use:   "get [workspace] <repo-slug> <key-id>",
		Short: "Get deploy key details",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, idStr, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}

			keyID, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid key ID: %s", idStr)
			}

			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			key, err := client.GetDeployKey(workspace, repoSlug, keyID)
			if err != nil {
				return err
			}

			fmt.Printf("ID:       %d\n", key.ID)
			fmt.Printf("Label:    %s\n", key.Label)
			fmt.Printf("Comment:  %s\n", key.Comment)
			fmt.Printf("Created:  %s\n", key.CreatedOn)
			fmt.Printf("Key:      %s\n", key.Key)
			return nil
		},
	})

	// deploy-key create
	dkCreateCmd := &cobra.Command{
		Use:   "create [workspace] <repo-slug>",
		Short: "Add a deploy key to a repository",
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

			keyContent, _ := cmd.Flags().GetString("key")
			label, _ := cmd.Flags().GetString("label")

			if keyContent == "" || label == "" {
				return fmt.Errorf("--key and --label are required")
			}

			key, err := client.CreateDeployKey(workspace, repoSlug, &bitbucket.CreateDeployKeyRequest{
				Key:   keyContent,
				Label: label,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Created deploy key: %s (ID: %d)\n", key.Label, key.ID)
			return nil
		},
	}
	dkCreateCmd.Flags().String("key", "", "SSH public key content (required)")
	dkCreateCmd.Flags().String("label", "", "Label for the key (required)")
	bbDeployKeyCmd.AddCommand(dkCreateCmd)

	// deploy-key delete
	bbDeployKeyCmd.AddCommand(&cobra.Command{
		Use:   "delete [workspace] <repo-slug> <key-id>",
		Short: "Delete a deploy key",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, idStr, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}
			keyID, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid key ID: %s", idStr)
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeleteDeployKey(workspace, repoSlug, keyID); err != nil {
				return err
			}
			fmt.Printf("Deleted deploy key: %s\n", idStr)
			return nil
		},
	})
}
