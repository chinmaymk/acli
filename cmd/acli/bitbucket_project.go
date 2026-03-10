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
		Use:     "list [workspace]",
		Short:   "List projects in a workspace",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, err := defaultWorkspace(cmd, args, 0)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			projects, err := client.ListProjects(workspace)
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
		Use:   "get [workspace] <project-key>",
		Short: "Get project details",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var workspace, projectKey string
			if len(args) >= 2 {
				workspace = args[0]
				projectKey = args[1]
			} else {
				var err error
				workspace, err = defaultWorkspace(cmd, nil, 0)
				if err != nil {
					return err
				}
				projectKey = args[0]
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			project, err := client.GetProject(workspace, projectKey)
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
		Use:   "create [workspace]",
		Short: "Create a project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, err := defaultWorkspace(cmd, args, 0)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
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

			project, err := client.CreateProject(workspace, &bitbucket.CreateProjectRequest{
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
		Use:   "delete [workspace] <project-key>",
		Short: "Delete a project",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var workspace, projectKey string
			if len(args) >= 2 {
				workspace = args[0]
				projectKey = args[1]
			} else {
				var err error
				workspace, err = defaultWorkspace(cmd, nil, 0)
				if err != nil {
					return err
				}
				projectKey = args[0]
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeleteProject(workspace, projectKey); err != nil {
				return err
			}
			fmt.Printf("Deleted project: %s\n", projectKey)
			return nil
		},
	})
}
