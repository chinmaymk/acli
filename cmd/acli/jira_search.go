package acli

import (
	"fmt"
	"strings"

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
		all, _ := cmd.Flags().GetBool("all")

		if len(fields) == 0 {
			fields = []string{"summary"}
		}

		results, err := client.SearchJQL(jql, startAt, maxResults, fields, nil)
		if err != nil {
			return err
		}

		if all {
			for len(results.Issues) < results.Total {
				next, err := client.SearchJQL(jql, startAt+len(results.Issues), maxResults, fields, nil)
				if err != nil {
					return err
				}
				if len(next.Issues) == 0 {
					break
				}
				results.Issues = append(results.Issues, next.Issues...)
			}
		}

		if isJSONOutput(cmd) {
			return outputJSON(results)
		}

		w := newTabWriter()
		headers := []string{"KEY"}
		for _, f := range fields {
			headers = append(headers, issueFieldHeader(f))
		}
		_, _ = fmt.Fprintln(w, strings.Join(headers, "\t"))
		for _, issue := range results.Issues {
			row := []string{issue.Key}
			for _, f := range fields {
				row = append(row, renderIssueField(f, issue))
			}
			_, _ = fmt.Fprintln(w, strings.Join(row, "\t"))
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(results.Issues), results.Total)
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
		all, _ := cmd.Flags().GetBool("all")

		var filters []jira.Filter
		total := 0

		if favourites {
			filters, err = client.GetFavouriteFilters()
			if err != nil {
				return err
			}
			total = len(filters)
		} else if mine {
			filters, err = client.GetMyFilters()
			if err != nil {
				return err
			}
			total = len(filters)
		} else {
			name, _ := cmd.Flags().GetString("name")
			maxResults, _ := cmd.Flags().GetInt("max-results")
			startAt, _ := cmd.Flags().GetInt("start-at")
			page, err := client.SearchFilters(name, startAt, maxResults)
			if err != nil {
				return err
			}
			filters = page.Values
			total = page.Total
			if all {
				for !page.IsLast && len(filters) < page.Total {
					next, err := client.SearchFilters(name, startAt+len(filters), maxResults)
					if err != nil {
						return err
					}
					if len(next.Values) == 0 {
						break
					}
					filters = append(filters, next.Values...)
					page.IsLast = next.IsLast
				}
			}
		}

		if isJSONOutput(cmd) {
			return outputJSON(filters)
		}

		w := newTabWriter()
		_, _ = fmt.Fprintln(w, "ID\tNAME\tOWNER\tJQL")
		for _, f := range filters {
			owner := ""
			if f.Owner != nil {
				owner = f.Owner.DisplayName
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", f.ID, f.Name, owner, truncate(f.JQL, 60))
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(filters), total)
		return nil
	},
}

var jiraFilterGetCmd = &cobra.Command{
	Use:   "get <filter-id>",
	Short: "Get filter details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		filter, err := client.GetFilter(args[0])
		if err != nil {
			return err
		}

		if isJSONOutput(cmd) {
			return outputJSON(filter)
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

		return outputResult(cmd, "created", created.ID, fmt.Sprintf("Filter created: %s (ID: %s)", created.Name, created.ID), created)
	},
}

var jiraFilterUpdateCmd = &cobra.Command{
	Use:   "update <filter-id>",
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

		return outputResult(cmd, "updated", updated.ID, fmt.Sprintf("Filter updated: %s (ID: %s)", updated.Name, updated.ID), updated)
	},
}

var jiraFilterDeleteCmd = &cobra.Command{
	Use:   "delete <filter-id>",
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

		return outputResult(cmd, "deleted", args[0], fmt.Sprintf("Filter %s deleted", args[0]), nil)
	},
}

func init() {
	// Search
	jiraSearchCmd.Flags().String("jql", "", "JQL query (required)")
	_ = jiraSearchCmd.MarkFlagRequired("jql")
	jiraSearchCmd.Flags().Int("max-results", 50, "Maximum number of results per page")
	jiraSearchCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraSearchCmd.Flags().StringSlice("fields", nil, "Fields to return")
	addAllFlag(jiraSearchCmd)
	jiraSearchCmd.Flags().Bool("json", false, "Output as JSON")

	// Filter
	jiraFilterListCmd.Flags().String("name", "", "Filter by name")
	jiraFilterListCmd.Flags().Int("max-results", 50, "Maximum number of results per page")
	jiraFilterListCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraFilterListCmd.Flags().Bool("favourites", false, "Show favourite filters")
	jiraFilterListCmd.Flags().Bool("mine", false, "Show my filters")
	addAllFlag(jiraFilterListCmd)
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
