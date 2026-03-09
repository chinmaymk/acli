package acli

import (
	"fmt"

	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

// --- Search command ---

var jiraSearchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Search issues using JQL",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		jql, _ := cmd.Flags().GetString("jql")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")
		fields, _ := cmd.Flags().GetStringSlice("fields")
		jsonOut, _ := cmd.Flags().GetBool("json")

		results, err := client.SearchJQL(jql, startAt, maxResults, fields, nil)
		if err != nil {
			return err
		}

		if jsonOut {
			return printJSON(results)
		}

		w := newTabWriter()
		fmt.Fprintln(w, "KEY\tTYPE\tSTATUS\tPRIORITY\tASSIGNEE\tSUMMARY")
		for _, issue := range results.Issues {
			printIssueRow(w, issue)
		}
		w.Flush()
		return nil
	},
}

// --- Filter commands ---

var jiraFilterCmd = &cobra.Command{
	Use:     "filter",
	Aliases: []string{"f"},
	Short:   "Manage filters",
	RunE:    helpRunE,
}

var jiraFilterListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List or search filters",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		favourites, _ := cmd.Flags().GetBool("favourites")
		mine, _ := cmd.Flags().GetBool("mine")

		var filters []jira.Filter

		if favourites {
			filters, err = client.GetFavouriteFilters()
			if err != nil {
				return err
			}
		} else if mine {
			filters, err = client.GetMyFilters()
			if err != nil {
				return err
			}
		} else {
			name, _ := cmd.Flags().GetString("name")
			maxResults, _ := cmd.Flags().GetInt("max-results")
			startAt, _ := cmd.Flags().GetInt("start-at")
			page, err := client.SearchFilters(name, startAt, maxResults)
			if err != nil {
				return err
			}
			filters = page.Values
		}

		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tOWNER\tJQL")
		for _, f := range filters {
			owner := ""
			if f.Owner != nil {
				owner = f.Owner.DisplayName
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", f.ID, f.Name, owner, truncate(f.JQL, 60))
		}
		w.Flush()
		return nil
	},
}

var jiraFilterGetCmd = &cobra.Command{
	Use:   "get [filter-id]",
	Short: "Get filter details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		jsonOut, _ := cmd.Flags().GetBool("json")

		filter, err := client.GetFilter(args[0])
		if err != nil {
			return err
		}

		if jsonOut {
			return printJSON(filter)
		}

		fmt.Printf("ID:          %s\n", filter.ID)
		fmt.Printf("Name:        %s\n", filter.Name)
		fmt.Printf("Description: %s\n", filter.Description)
		if filter.Owner != nil {
			fmt.Printf("Owner:       %s\n", filter.Owner.DisplayName)
		}
		fmt.Printf("JQL:         %s\n", filter.JQL)
		fmt.Printf("Favourite:   %v\n", filter.Favourite)
		fmt.Printf("View URL:    %s\n", filter.ViewURL)
		return nil
	},
}

var jiraFilterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a filter",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		jql, _ := cmd.Flags().GetString("jql")
		description, _ := cmd.Flags().GetString("description")
		favourite, _ := cmd.Flags().GetBool("favourite")

		filter := &jira.Filter{
			Name:        name,
			JQL:         jql,
			Description: description,
			Favourite:   favourite,
		}

		created, err := client.CreateFilter(filter)
		if err != nil {
			return err
		}

		fmt.Printf("Filter created: %s (ID: %s)\n", created.Name, created.ID)
		return nil
	},
}

var jiraFilterUpdateCmd = &cobra.Command{
	Use:   "update [filter-id]",
	Short: "Update a filter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		filter := &jira.Filter{}
		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			filter.Name = name
		}
		if cmd.Flags().Changed("jql") {
			jql, _ := cmd.Flags().GetString("jql")
			filter.JQL = jql
		}
		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			filter.Description = description
		}

		updated, err := client.UpdateFilter(args[0], filter)
		if err != nil {
			return err
		}

		fmt.Printf("Filter updated: %s (ID: %s)\n", updated.Name, updated.ID)
		return nil
	},
}

var jiraFilterDeleteCmd = &cobra.Command{
	Use:   "delete [filter-id]",
	Short: "Delete a filter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		if err := client.DeleteFilter(args[0]); err != nil {
			return err
		}

		fmt.Printf("Filter %s deleted\n", args[0])
		return nil
	},
}

func init() {
	// Search
	jiraSearchCmd.Flags().String("jql", "", "JQL query (required)")
	_ = jiraSearchCmd.MarkFlagRequired("jql")
	jiraSearchCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraSearchCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraSearchCmd.Flags().StringSlice("fields", nil, "Fields to return")
	jiraSearchCmd.Flags().Bool("json", false, "Output as JSON")

	// Filter
	jiraFilterListCmd.Flags().String("name", "", "Filter by name")
	jiraFilterListCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraFilterListCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraFilterListCmd.Flags().Bool("favourites", false, "Show favourite filters")
	jiraFilterListCmd.Flags().Bool("mine", false, "Show my filters")
	jiraFilterCmd.AddCommand(jiraFilterListCmd)

	jiraFilterGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraFilterCmd.AddCommand(jiraFilterGetCmd)

	jiraFilterCreateCmd.Flags().String("name", "", "Filter name (required)")
	_ = jiraFilterCreateCmd.MarkFlagRequired("name")
	jiraFilterCreateCmd.Flags().String("jql", "", "JQL query (required)")
	_ = jiraFilterCreateCmd.MarkFlagRequired("jql")
	jiraFilterCreateCmd.Flags().String("description", "", "Filter description")
	jiraFilterCreateCmd.Flags().Bool("favourite", false, "Mark as favourite")
	jiraFilterCmd.AddCommand(jiraFilterCreateCmd)

	jiraFilterUpdateCmd.Flags().String("name", "", "Filter name")
	jiraFilterUpdateCmd.Flags().String("jql", "", "JQL query")
	jiraFilterUpdateCmd.Flags().String("description", "", "Filter description")
	jiraFilterCmd.AddCommand(jiraFilterUpdateCmd)

	jiraFilterCmd.AddCommand(jiraFilterDeleteCmd)
}
