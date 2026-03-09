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
		Use:     "list <workspace> <repo-slug>",
		Short:   "List deploy keys",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			keys, err := client.ListDeployKeys(args[0], args[1])
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
		Use:   "get <workspace> <repo-slug> <key-id>",
		Short: "Get deploy key details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			keyID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid key ID: %s", args[2])
			}

			key, err := client.GetDeployKey(args[0], args[1], keyID)
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
		Use:   "create <workspace> <repo-slug>",
		Short: "Add a deploy key to a repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			keyContent, _ := cmd.Flags().GetString("key")
			label, _ := cmd.Flags().GetString("label")

			if keyContent == "" || label == "" {
				return fmt.Errorf("--key and --label are required")
			}

			key, err := client.CreateDeployKey(args[0], args[1], &bitbucket.CreateDeployKeyRequest{
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
		Use:   "delete <workspace> <repo-slug> <key-id>",
		Short: "Delete a deploy key",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			keyID, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid key ID: %s", args[2])
			}
			if err := client.DeleteDeployKey(args[0], args[1], keyID); err != nil {
				return err
			}
			fmt.Printf("Deleted deploy key: %s\n", args[2])
			return nil
		},
	})
}
