package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbDeploymentCmd = &cobra.Command{
	Use:     "deployment",
	Short:   "Manage deployments",
	Aliases: []string{"deploy"},
	RunE:    helpRunE,
}

func init() {
	// deployment list
	bbDeploymentCmd.AddCommand(&cobra.Command{
		Use:     "list <workspace> <repo-slug>",
		Short:   "List deployments",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			deployments, err := client.ListDeployments(args[0], args[1])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "UUID\tENVIRONMENT\tSTATUS\tRELEASE\tCOMMIT\tCREATED")
			for _, d := range deployments {
				hash := d.Release.Commit.Hash
				if len(hash) > 12 {
					hash = hash[:12]
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					d.UUID, d.Environment.Name, d.State.Status.Name,
					d.Release.Name, hash, d.Release.CreatedOn)
			}
			return w.Flush()
		},
	})

	// deployment get
	bbDeploymentCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug> <deployment-uuid>",
		Short: "Get deployment details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			d, err := client.GetDeployment(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			fmt.Printf("UUID:         %s\n", d.UUID)
			fmt.Printf("State:        %s\n", d.State.Name)
			fmt.Printf("Status:       %s\n", d.State.Status.Name)
			fmt.Printf("Environment:  %s\n", d.Environment.Name)
			fmt.Printf("Release:      %s\n", d.Release.Name)
			fmt.Printf("Commit:       %s\n", d.Release.Commit.Hash)
			if d.Release.Commit.Message != "" {
				fmt.Printf("Message:      %s\n", d.Release.Commit.Message)
			}
			fmt.Printf("Created:      %s\n", d.Release.CreatedOn)
			return nil
		},
	})
}
