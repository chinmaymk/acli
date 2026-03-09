package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbEnvironmentCmd = &cobra.Command{
	Use:     "environment",
	Short:   "Manage deployment environments",
	Aliases: []string{"env"},
	RunE:    helpRunE,
}

func init() {
	// environment list
	bbEnvironmentCmd.AddCommand(&cobra.Command{
		Use:     "list <workspace> <repo-slug>",
		Short:   "List deployment environments",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			envs, err := client.ListEnvironments(args[0], args[1])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "UUID\tNAME\tTYPE\tRANK")
			for _, e := range envs {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
					e.UUID, e.Name, e.EnvironmentType.Name, e.Rank)
			}
			return w.Flush()
		},
	})

	// environment get
	bbEnvironmentCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug> <environment-uuid>",
		Short: "Get environment details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			env, err := client.GetEnvironment(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			fmt.Printf("UUID:  %s\n", env.UUID)
			fmt.Printf("Name:  %s\n", env.Name)
			fmt.Printf("Type:  %s\n", env.EnvironmentType.Name)
			fmt.Printf("Rank:  %d\n", env.Rank)
			return nil
		},
	})

	// environment create
	envCreateCmd := &cobra.Command{
		Use:   "create <workspace> <repo-slug>",
		Short: "Create a deployment environment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			envType, _ := cmd.Flags().GetString("type")
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if envType == "" {
				envType = "Test"
			}

			rank := 0
			switch envType {
			case "Production":
				rank = 2
			case "Staging":
				rank = 1
			default:
				rank = 0
			}

			req := &bitbucket.CreateEnvironmentRequest{
				Name: name,
			}
			req.EnvironmentType.Name = envType
			req.EnvironmentType.Rank = rank

			env, err := client.CreateEnvironment(args[0], args[1], req)
			if err != nil {
				return err
			}

			fmt.Printf("Created environment: %s (UUID: %s)\n", env.Name, env.UUID)
			return nil
		},
	}
	envCreateCmd.Flags().String("name", "", "Environment name (required)")
	envCreateCmd.Flags().String("type", "Test", "Environment type (Test, Staging, Production)")
	bbEnvironmentCmd.AddCommand(envCreateCmd)

	// environment delete
	bbEnvironmentCmd.AddCommand(&cobra.Command{
		Use:   "delete <workspace> <repo-slug> <environment-uuid>",
		Short: "Delete a deployment environment",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			if err := client.DeleteEnvironment(args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("Deleted environment: %s\n", args[2])
			return nil
		},
	})
}
