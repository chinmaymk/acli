package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbProjectCmd = &cobra.Command{
	Use:     "project",
	Short:   "Manage projects",
	Aliases: []string{"proj"},
	RunE:    helpRunE,
}

func init() {
	// project list
	bbProjectCmd.AddCommand(&cobra.Command{
		Use:     "list <workspace>",
		Short:   "List projects in a workspace",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			projects, err := client.ListProjects(args[0])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "KEY\tNAME\tPRIVATE\tUPDATED")
			for _, p := range projects {
				fmt.Fprintf(w, "%s\t%s\t%v\t%s\n",
					p.Key, p.Name, p.IsPrivate, p.UpdatedOn)
			}
			return w.Flush()
		},
	})

	// project get
	bbProjectCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <project-key>",
		Short: "Get project details",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			project, err := client.GetProject(args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Printf("Key:          %s\n", project.Key)
			fmt.Printf("Name:         %s\n", project.Name)
			fmt.Printf("Description:  %s\n", project.Description)
			fmt.Printf("Private:      %v\n", project.IsPrivate)
			fmt.Printf("Created:      %s\n", project.CreatedOn)
			fmt.Printf("Updated:      %s\n", project.UpdatedOn)
			fmt.Printf("URL:          %s\n", project.Links.HTML.Href)
			return nil
		},
	})

	// project create
	projectCreateCmd := &cobra.Command{
		Use:   "create <workspace>",
		Short: "Create a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			key, _ := cmd.Flags().GetString("key")
			desc, _ := cmd.Flags().GetString("description")
			isPrivate, _ := cmd.Flags().GetBool("private")

			if name == "" || key == "" {
				return fmt.Errorf("--name and --key are required")
			}

			project, err := client.CreateProject(args[0], &bitbucket.CreateProjectRequest{
				Name:        name,
				Key:         key,
				Description: desc,
				IsPrivate:   isPrivate,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Created project: %s (%s)\n", project.Name, project.Key)
			fmt.Printf("URL: %s\n", project.Links.HTML.Href)
			return nil
		},
	}
	projectCreateCmd.Flags().String("name", "", "Project name (required)")
	projectCreateCmd.Flags().String("key", "", "Project key (required)")
	projectCreateCmd.Flags().String("description", "", "Project description")
	projectCreateCmd.Flags().Bool("private", true, "Make project private")
	bbProjectCmd.AddCommand(projectCreateCmd)

	// project delete
	bbProjectCmd.AddCommand(&cobra.Command{
		Use:   "delete <workspace> <project-key>",
		Short: "Delete a project",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			if err := client.DeleteProject(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Deleted project: %s\n", args[1])
			return nil
		},
	})
}
