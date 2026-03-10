package acli

import (
	"fmt"

	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

// printUsersTable prints a table of users.
func printUsersTable(users []jira.UserDetails) {
	w := newTabWriter()
	fmt.Fprintln(w, "ACCOUNT_ID\tDISPLAY_NAME\tEMAIL\tACTIVE")
	for _, u := range users {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", u.AccountID, u.DisplayName, u.EmailAddress, u.Active)
	}
	w.Flush()
}

// --- User commands ---

var jiraUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"u"},
	Short:   "Manage users",
	RunE:    helpRunE,
}

var jiraUserGetCmd = &cobra.Command{
	Use:   "get [account-id]",
	Short: "Get user details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		jsonOut, _ := cmd.Flags().GetBool("json")
		user, err := client.GetUser(args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(user)
		}
		fmt.Printf("Account ID:   %s\n", user.AccountID)
		fmt.Printf("Display Name: %s\n", user.DisplayName)
		fmt.Printf("Email:        %s\n", user.EmailAddress)
		fmt.Printf("Active:       %v\n", user.Active)
		fmt.Printf("Time Zone:    %s\n", user.TimeZone)
		fmt.Printf("Account Type: %s\n", user.AccountType)
		return nil
	},
}

var jiraUserSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search users",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		query, _ := cmd.Flags().GetString("query")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")
		jsonOut, _ := cmd.Flags().GetBool("json")

		users, err := client.FindUsers(query, startAt, maxResults)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(users)
		}
		printUsersTable(users)
		return nil
	},
}

var jiraUserAssignableCmd = &cobra.Command{
	Use:   "assignable",
	Short: "Find assignable users",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		query, _ := cmd.Flags().GetString("query")
		project, _ := defaultProject(cmd)
		issueKey, _ := cmd.Flags().GetString("issue-key")
		maxResults, _ := cmd.Flags().GetInt("max-results")

		users, err := client.FindUsersAssignable(query, project, issueKey, 0, maxResults)
		if err != nil {
			return err
		}
		printUsersTable(users)
		return nil
	},
}

var jiraUserMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		user, err := client.GetCurrentUser()
		if err != nil {
			return err
		}
		fmt.Printf("Account ID:   %s\n", user.AccountID)
		fmt.Printf("Display Name: %s\n", user.DisplayName)
		fmt.Printf("Email:        %s\n", user.EmailAddress)
		fmt.Printf("Active:       %v\n", user.Active)
		fmt.Printf("Time Zone:    %s\n", user.TimeZone)
		return nil
	},
}

var jiraUserListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")

		users, err := client.GetAllUsers(startAt, maxResults)
		if err != nil {
			return err
		}
		printUsersTable(users)
		return nil
	},
}

var jiraUserCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		email, _ := cmd.Flags().GetString("email")
		displayName, _ := cmd.Flags().GetString("display-name")

		body := map[string]interface{}{
			"emailAddress": email,
			"displayName":  displayName,
		}
		user, err := client.CreateUser(body)
		if err != nil {
			return err
		}
		fmt.Printf("User created: %s (Account ID: %s)\n", user.DisplayName, user.AccountID)
		return nil
	},
}

var jiraUserDeleteCmd = &cobra.Command{
	Use:   "delete [account-id]",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteUser(args[0]); err != nil {
			return err
		}
		fmt.Printf("User %s deleted\n", args[0])
		return nil
	},
}

// --- Group commands ---

var jiraGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"g"},
	Short:   "Manage groups",
	RunE:    helpRunE,
}

var jiraGroupListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")
		jsonOut, _ := cmd.Flags().GetBool("json")

		page, err := client.GetBulkGroups(startAt, maxResults)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(page)
		}
		w := newTabWriter()
		fmt.Fprintln(w, "GROUP_ID\tNAME")
		for _, g := range page.Values {
			fmt.Fprintf(w, "%s\t%s\n", g.GroupID, g.Name)
		}
		w.Flush()
		return nil
	},
}

var jiraGroupGetCmd = &cobra.Command{
	Use:   "get [group-name]",
	Short: "Get group details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		jsonOut, _ := cmd.Flags().GetBool("json")
		group, err := client.GetGroup(args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(group)
		}
		fmt.Printf("Group ID: %s\n", group.GroupID)
		fmt.Printf("Name:     %s\n", group.Name)
		return nil
	},
}

var jiraGroupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		group, err := client.CreateGroup(name)
		if err != nil {
			return err
		}
		fmt.Printf("Group created: %s (ID: %s)\n", group.Name, group.GroupID)
		return nil
	},
}

var jiraGroupDeleteCmd = &cobra.Command{
	Use:   "delete [group-name]",
	Short: "Delete a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteGroup(args[0]); err != nil {
			return err
		}
		fmt.Printf("Group %s deleted\n", args[0])
		return nil
	},
}

var jiraGroupMembersCmd = &cobra.Command{
	Use:   "members [group-name]",
	Short: "List group members",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")

		members, err := client.GetGroupMembers(args[0], startAt, maxResults)
		if err != nil {
			return err
		}
		printUsersTable(members.Values)
		return nil
	},
}

var jiraGroupAddUserCmd = &cobra.Command{
	Use:   "add-user [group-name]",
	Short: "Add a user to a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		accountID, _ := cmd.Flags().GetString("account-id")
		_, err = client.AddUserToGroup(args[0], accountID)
		if err != nil {
			return err
		}
		fmt.Printf("User %s added to group %s\n", accountID, args[0])
		return nil
	},
}

var jiraGroupRemoveUserCmd = &cobra.Command{
	Use:   "remove-user [group-name]",
	Short: "Remove a user from a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		accountID, _ := cmd.Flags().GetString("account-id")
		if err := client.RemoveUserFromGroup(args[0], accountID); err != nil {
			return err
		}
		fmt.Printf("User %s removed from group %s\n", accountID, args[0])
		return nil
	},
}

var jiraGroupSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		query, _ := cmd.Flags().GetString("query")
		maxResults, _ := cmd.Flags().GetInt("max-results")

		found, err := client.FindGroups(query, maxResults)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "GROUP_ID\tNAME")
		for _, g := range found.Groups {
			fmt.Fprintf(w, "%s\t%s\n", g.GroupID, g.Name)
		}
		w.Flush()
		return nil
	},
}

func init() {
	// User
	jiraUserGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraUserCmd.AddCommand(jiraUserGetCmd)

	jiraUserSearchCmd.Flags().String("query", "", "Search query (required)")
	_ = jiraUserSearchCmd.MarkFlagRequired("query")
	jiraUserSearchCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraUserSearchCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraUserSearchCmd.Flags().Bool("json", false, "Output as JSON")
	jiraUserCmd.AddCommand(jiraUserSearchCmd)

	jiraUserAssignableCmd.Flags().String("query", "", "Search query")
	jiraUserAssignableCmd.Flags().String("project", "", "Project key (uses profile default if not set)")
	jiraUserAssignableCmd.Flags().String("issue-key", "", "Issue key")
	jiraUserAssignableCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraUserCmd.AddCommand(jiraUserAssignableCmd)

	jiraUserCmd.AddCommand(jiraUserMeCmd)

	jiraUserListCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraUserListCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraUserCmd.AddCommand(jiraUserListCmd)

	jiraUserCreateCmd.Flags().String("email", "", "User email address (required)")
	_ = jiraUserCreateCmd.MarkFlagRequired("email")
	jiraUserCreateCmd.Flags().String("display-name", "", "User display name (required)")
	_ = jiraUserCreateCmd.MarkFlagRequired("display-name")
	jiraUserCmd.AddCommand(jiraUserCreateCmd)

	jiraUserCmd.AddCommand(jiraUserDeleteCmd)

	// Group
	jiraGroupListCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraGroupListCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraGroupListCmd.Flags().Bool("json", false, "Output as JSON")
	jiraGroupCmd.AddCommand(jiraGroupListCmd)

	jiraGroupGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraGroupCmd.AddCommand(jiraGroupGetCmd)

	jiraGroupCreateCmd.Flags().String("name", "", "Group name (required)")
	_ = jiraGroupCreateCmd.MarkFlagRequired("name")
	jiraGroupCmd.AddCommand(jiraGroupCreateCmd)

	jiraGroupCmd.AddCommand(jiraGroupDeleteCmd)

	jiraGroupMembersCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraGroupMembersCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraGroupCmd.AddCommand(jiraGroupMembersCmd)

	jiraGroupAddUserCmd.Flags().String("account-id", "", "User account ID (required)")
	_ = jiraGroupAddUserCmd.MarkFlagRequired("account-id")
	jiraGroupCmd.AddCommand(jiraGroupAddUserCmd)

	jiraGroupRemoveUserCmd.Flags().String("account-id", "", "User account ID (required)")
	_ = jiraGroupRemoveUserCmd.MarkFlagRequired("account-id")
	jiraGroupCmd.AddCommand(jiraGroupRemoveUserCmd)

	jiraGroupSearchCmd.Flags().String("query", "", "Search query")
	jiraGroupSearchCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraGroupCmd.AddCommand(jiraGroupSearchCmd)
}
