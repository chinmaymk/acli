package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbWorkspaceCmd = &cobra.Command{
	Use:     "workspace",
	Short:   "Manage workspaces",
	Aliases: []string{"ws"},
	RunE:    helpRunE,
}

func init() {
	// workspace list
	bbWorkspaceCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "List workspaces",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			workspaces, err := client.ListWorkspaces()
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SLUG\tNAME\tUUID")
			for _, ws := range workspaces {
				fmt.Fprintf(w, "%s\t%s\t%s\n", ws.Slug, ws.Name, ws.UUID)
			}
			return w.Flush()
		},
	})

	// workspace get
	bbWorkspaceCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace>",
		Short: "Get workspace details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			ws, err := client.GetWorkspace(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Name:     %s\n", ws.Name)
			fmt.Printf("Slug:     %s\n", ws.Slug)
			fmt.Printf("UUID:     %s\n", ws.UUID)
			fmt.Printf("Private:  %v\n", ws.IsPrivate)
			fmt.Printf("Created:  %s\n", ws.CreatedOn)
			fmt.Printf("URL:      %s\n", ws.Links.HTML.Href)
			return nil
		},
	})

	// workspace members
	bbWorkspaceCmd.AddCommand(&cobra.Command{
		Use:   "members <workspace>",
		Short: "List workspace members",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			members, err := client.ListWorkspaceMembers(args[0])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "DISPLAY NAME\tNICKNAME\tUUID")
			for _, m := range members {
				fmt.Fprintf(w, "%s\t%s\t%s\n",
					m.User.DisplayName, m.User.Nickname, m.User.UUID)
			}
			return w.Flush()
		},
	})

	// workspace permissions
	bbWorkspaceCmd.AddCommand(&cobra.Command{
		Use:   "permissions <workspace>",
		Short: "List user permissions in workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			perms, err := client.ListWorkspacePermissions(args[0])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "USER\tPERMISSION")
			for _, p := range perms {
				fmt.Fprintf(w, "%s\t%s\n", p.User.DisplayName, p.Permission)
			}
			return w.Flush()
		},
	})
}
