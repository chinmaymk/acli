package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbRepoCmd = &cobra.Command{
	Use:     "repo",
	Short:   "Manage repositories",
	Aliases: []string{"r"},
	RunE:    helpRunE,
}

func init() {
	// repo list
	repoListCmd := &cobra.Command{
		Use:     "list <workspace>",
		Short:   "List repositories in a workspace",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			role, _ := cmd.Flags().GetString("role")
			q, _ := cmd.Flags().GetString("query")
			sort, _ := cmd.Flags().GetString("sort")
			repos, err := client.ListRepositories(args[0], &bitbucket.ListReposOptions{
				Role: role,
				Q:    q,
				Sort: sort,
			})
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSLUG\tLANGUAGE\tPRIVATE\tUPDATED")
			for _, r := range repos {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n",
					r.FullName, r.Slug, r.Language, r.IsPrivate, r.UpdatedOn)
			}
			return w.Flush()
		},
	}
	repoListCmd.Flags().String("role", "", "Filter by role (admin, contributor, member, owner)")
	repoListCmd.Flags().String("query", "", "Filter with query (Bitbucket query syntax)")
	repoListCmd.Flags().String("sort", "", "Sort field (e.g. -updated_on)")
	bbRepoCmd.AddCommand(repoListCmd)

	// repo get
	bbRepoCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug>",
		Short: "Get repository details",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			repo, err := client.GetRepository(args[0], args[1])
			if err != nil {
				return err
			}

			mainBranch := "N/A"
			if repo.MainBranch != nil {
				mainBranch = repo.MainBranch.Name
			}

			fmt.Printf("Name:         %s\n", repo.FullName)
			fmt.Printf("Description:  %s\n", repo.Description)
			fmt.Printf("Language:     %s\n", repo.Language)
			fmt.Printf("SCM:          %s\n", repo.SCM)
			fmt.Printf("Private:      %v\n", repo.IsPrivate)
			fmt.Printf("Main Branch:  %s\n", mainBranch)
			fmt.Printf("Created:      %s\n", repo.CreatedOn)
			fmt.Printf("Updated:      %s\n", repo.UpdatedOn)
			fmt.Printf("URL:          %s\n", repo.Links.HTML.Href)
			return nil
		},
	})

	// repo create
	repoCreateCmd := &cobra.Command{
		Use:   "create <workspace>",
		Short: "Create a new repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				return fmt.Errorf("--name is required")
			}

			desc, _ := cmd.Flags().GetString("description")
			language, _ := cmd.Flags().GetString("language")
			isPrivate, _ := cmd.Flags().GetBool("private")
			projectKey, _ := cmd.Flags().GetString("project")

			req := &bitbucket.CreateRepoRequest{
				SCM:         "git",
				Name:        name,
				IsPrivate:   isPrivate,
				Description: desc,
				Language:    language,
			}
			if projectKey != "" {
				req.Project = &struct {
					Key string `json:"key"`
				}{Key: projectKey}
			}

			repo, err := client.CreateRepository(args[0], req)
			if err != nil {
				return err
			}

			fmt.Printf("Created repository: %s\n", repo.FullName)
			fmt.Printf("URL: %s\n", repo.Links.HTML.Href)
			return nil
		},
	}
	repoCreateCmd.Flags().String("name", "", "Repository name (required)")
	repoCreateCmd.Flags().String("description", "", "Repository description")
	repoCreateCmd.Flags().String("language", "", "Programming language")
	repoCreateCmd.Flags().Bool("private", true, "Make repository private")
	repoCreateCmd.Flags().String("project", "", "Project key to assign to")
	bbRepoCmd.AddCommand(repoCreateCmd)

	// repo delete
	bbRepoCmd.AddCommand(&cobra.Command{
		Use:   "delete <workspace> <repo-slug>",
		Short: "Delete a repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			if err := client.DeleteRepository(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Deleted repository: %s/%s\n", args[0], args[1])
			return nil
		},
	})

	// repo fork
	repoForkCmd := &cobra.Command{
		Use:   "fork <workspace> <repo-slug>",
		Short: "Fork a repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			targetWorkspace, _ := cmd.Flags().GetString("workspace")

			req := &bitbucket.ForkRepoRequest{
				Name: name,
			}
			if targetWorkspace != "" {
				req.Workspace = &struct {
					Slug string `json:"slug"`
				}{Slug: targetWorkspace}
			}

			repo, err := client.ForkRepository(args[0], args[1], req)
			if err != nil {
				return err
			}

			fmt.Printf("Forked repository: %s\n", repo.FullName)
			fmt.Printf("URL: %s\n", repo.Links.HTML.Href)
			return nil
		},
	}
	repoForkCmd.Flags().String("name", "", "Name for the forked repo")
	repoForkCmd.Flags().String("workspace", "", "Target workspace for the fork")
	bbRepoCmd.AddCommand(repoForkCmd)

	// repo forks (list forks)
	bbRepoCmd.AddCommand(&cobra.Command{
		Use:   "forks <workspace> <repo-slug>",
		Short: "List forks of a repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			forks, err := client.ListForks(args[0], args[1])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSLUG\tOWNER\tPRIVATE")
			for _, r := range forks {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
					r.FullName, r.Slug, r.Owner.DisplayName, r.IsPrivate)
			}
			return w.Flush()
		},
	})
}
