package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbSnippetCmd = &cobra.Command{
	Use:     "snippet",
	Short:   "Manage snippets",
	Aliases: []string{"snip"},
	RunE:    helpRunE,
}

func init() {
	// snippet list
	bbSnippetCmd.AddCommand(&cobra.Command{
		Use:     "list <workspace>",
		Short:   "List snippets in a workspace",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			snippets, err := client.ListSnippets(args[0])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tPRIVATE\tCREATED\tOWNER")
			for _, s := range snippets {
				fmt.Fprintf(w, "%d\t%s\t%v\t%s\t%s\n",
					s.ID, s.Title, s.IsPrivate, s.CreatedOn, s.Owner.DisplayName)
			}
			return w.Flush()
		},
	})

	// snippet get
	bbSnippetCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <snippet-id>",
		Short: "Get snippet details",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			snippet, err := client.GetSnippet(args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Printf("ID:       %d\n", snippet.ID)
			fmt.Printf("Title:    %s\n", snippet.Title)
			fmt.Printf("Private:  %v\n", snippet.IsPrivate)
			fmt.Printf("Owner:    %s\n", snippet.Owner.DisplayName)
			fmt.Printf("Created:  %s\n", snippet.CreatedOn)
			fmt.Printf("Updated:  %s\n", snippet.UpdatedOn)
			fmt.Printf("URL:      %s\n", snippet.Links.HTML.Href)
			return nil
		},
	})

	// snippet create
	snippetCreateCmd := &cobra.Command{
		Use:   "create <workspace>",
		Short: "Create a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			title, _ := cmd.Flags().GetString("title")
			isPrivate, _ := cmd.Flags().GetBool("private")

			if title == "" {
				return fmt.Errorf("--title is required")
			}

			snippet, err := client.CreateSnippet(args[0], &bitbucket.CreateSnippetRequest{
				Title:     title,
				IsPrivate: isPrivate,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Created snippet: %s (ID: %d)\n", snippet.Title, snippet.ID)
			fmt.Printf("URL: %s\n", snippet.Links.HTML.Href)
			return nil
		},
	}
	snippetCreateCmd.Flags().String("title", "", "Snippet title (required)")
	snippetCreateCmd.Flags().Bool("private", false, "Make snippet private")
	bbSnippetCmd.AddCommand(snippetCreateCmd)

	// snippet delete
	bbSnippetCmd.AddCommand(&cobra.Command{
		Use:   "delete <workspace> <snippet-id>",
		Short: "Delete a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			if err := client.DeleteSnippet(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Deleted snippet: %s\n", args[1])
			return nil
		},
	})
}
