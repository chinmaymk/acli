package acli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

// ============================================================================
// Roles
// ============================================================================

var jiraRoleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage project roles",
	RunE:  helpRunE,
}

var jiraRoleListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all project roles",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		roles, err := client.GetAllRoles()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
		for _, r := range roles {
			fmt.Fprintf(w, "%d\t%s\t%s\n", r.ID, r.Name, r.Description)
		}
		return w.Flush()
	},
}

var jiraRoleGetCmd = &cobra.Command{
	Use:   "get <role-id>",
	Short: "Get a project role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid role ID: %w", err)
		}
		role, err := client.GetRole(id)
		if err != nil {
			return err
		}
		return printJSON(role)
	},
}

var jiraRoleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project role",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		body := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			body["description"] = description
		}
		role, err := client.CreateRole(body)
		if err != nil {
			return err
		}
		return printJSON(role)
	},
}

var jiraRoleDeleteCmd = &cobra.Command{
	Use:   "delete <role-id>",
	Short: "Delete a project role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid role ID: %w", err)
		}
		if err := client.DeleteRole(id); err != nil {
			return err
		}
		fmt.Printf("Role %d deleted.\n", id)
		return nil
	},
}

// ============================================================================
// Issue Links
// ============================================================================

var jiraIssueLinkCmd = &cobra.Command{
	Use:     "issuelink",
	Short:   "Manage issue links",
	Aliases: []string{"il"},
	RunE:    helpRunE,
}

var jiraIssueLinkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an issue link",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		inward, _ := cmd.Flags().GetString("inward-issue")
		outward, _ := cmd.Flags().GetString("outward-issue")
		linkType, _ := cmd.Flags().GetString("type")
		link := &jira.IssueLink{
			Type:         &jira.IssueLinkType{Name: linkType},
			InwardIssue:  &jira.Issue{Key: inward},
			OutwardIssue: &jira.Issue{Key: outward},
		}
		if err := client.CreateIssueLink(link); err != nil {
			return err
		}
		fmt.Println("Issue link created.")
		return nil
	},
}

var jiraIssueLinkGetCmd = &cobra.Command{
	Use:   "get <link-id>",
	Short: "Get an issue link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		link, err := client.GetIssueLink(args[0])
		if err != nil {
			return err
		}
		return printJSON(link)
	},
}

var jiraIssueLinkDeleteCmd = &cobra.Command{
	Use:   "delete <link-id>",
	Short: "Delete an issue link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteIssueLink(args[0]); err != nil {
			return err
		}
		fmt.Printf("Issue link %s deleted.\n", args[0])
		return nil
	},
}

// ============================================================================
// Issue Link Types
// ============================================================================

var jiraIssueLinkTypeCmd = &cobra.Command{
	Use:     "issuelinktype",
	Short:   "Manage issue link types",
	Aliases: []string{"ilt"},
	RunE:    helpRunE,
}

var jiraIssueLinkTypeListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all issue link types",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		types, err := client.GetIssueLinkTypes()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tINWARD\tOUTWARD")
		for _, t := range types {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t.ID, t.Name, t.Inward, t.Outward)
		}
		return w.Flush()
	},
}

var jiraIssueLinkTypeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an issue link type",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		lt, err := client.GetIssueLinkType(args[0])
		if err != nil {
			return err
		}
		return printJSON(lt)
	},
}

var jiraIssueLinkTypeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an issue link type",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		inward, _ := cmd.Flags().GetString("inward")
		outward, _ := cmd.Flags().GetString("outward")
		lt, err := client.CreateIssueLinkType(&jira.IssueLinkType{
			Name:    name,
			Inward:  inward,
			Outward: outward,
		})
		if err != nil {
			return err
		}
		return printJSON(lt)
	},
}

var jiraIssueLinkTypeUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an issue link type",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		lt := &jira.IssueLinkType{}
		if cmd.Flags().Changed("name") {
			lt.Name, _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("inward") {
			lt.Inward, _ = cmd.Flags().GetString("inward")
		}
		if cmd.Flags().Changed("outward") {
			lt.Outward, _ = cmd.Flags().GetString("outward")
		}
		result, err := client.UpdateIssueLinkType(args[0], lt)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var jiraIssueLinkTypeDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an issue link type",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteIssueLinkType(args[0]); err != nil {
			return err
		}
		fmt.Printf("Issue link type %s deleted.\n", args[0])
		return nil
	},
}

// ============================================================================
// Screens
// ============================================================================

var jiraScreenCmd = &cobra.Command{
	Use:   "screen",
	Short: "Manage screens",
	RunE:  helpRunE,
}

var jiraScreenListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List screens",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		result, err := client.GetScreens(startAt, maxResults)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
		for _, s := range result.Values {
			fmt.Fprintf(w, "%d\t%s\t%s\n", s.ID, s.Name, s.Description)
		}
		return w.Flush()
	},
}

var jiraScreenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a screen",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		body := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			body["description"] = description
		}
		screen, err := client.CreateScreen(body)
		if err != nil {
			return err
		}
		return printJSON(screen)
	},
}

var jiraScreenDeleteCmd = &cobra.Command{
	Use:   "delete <screen-id>",
	Short: "Delete a screen",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %w", err)
		}
		if err := client.DeleteScreen(id); err != nil {
			return err
		}
		fmt.Printf("Screen %d deleted.\n", id)
		return nil
	},
}

var jiraScreenTabsCmd = &cobra.Command{
	Use:   "tabs <screen-id>",
	Short: "List tabs for a screen",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		screenId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %w", err)
		}
		tabs, err := client.GetScreenTabs(screenId)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME")
		for _, t := range tabs {
			fmt.Fprintf(w, "%d\t%s\n", t.ID, t.Name)
		}
		return w.Flush()
	},
}

var jiraScreenTabFieldsCmd = &cobra.Command{
	Use:   "tab-fields <screen-id> <tab-id>",
	Short: "List fields for a screen tab",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		screenId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %w", err)
		}
		tabId, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid tab ID: %w", err)
		}
		fields, err := client.GetScreenTabFields(screenId, tabId)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME")
		for _, f := range fields {
			fmt.Fprintf(w, "%s\t%s\n", f.ID, f.Name)
		}
		return w.Flush()
	},
}

// ============================================================================
// Workflows
// ============================================================================

var jiraWorkflowCmd = &cobra.Command{
	Use:     "workflow",
	Short:   "Manage workflows",
	Aliases: []string{"wf"},
	RunE:    helpRunE,
}

var jiraWorkflowListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all workflows",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		workflows, err := client.GetWorkflows()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tDEFAULT")
		for _, wf := range workflows {
			fmt.Fprintf(w, "%s\t%s\t%v\n", wf.Name, wf.Description, wf.IsDefault)
		}
		return w.Flush()
	},
}

// ============================================================================
// Workflow Schemes
// ============================================================================

var jiraWorkflowSchemeCmd = &cobra.Command{
	Use:     "workflowscheme",
	Short:   "Manage workflow schemes",
	Aliases: []string{"wfs"},
	RunE:    helpRunE,
}

var jiraWorkflowSchemeListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List workflow schemes",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		result, err := client.GetWorkflowSchemes(startAt, maxResults)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
		for _, s := range result.Values {
			fmt.Fprintf(w, "%d\t%s\t%s\n", s.ID, s.Name, s.Description)
		}
		return w.Flush()
	},
}

var jiraWorkflowSchemeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a workflow scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid workflow scheme ID: %w", err)
		}
		scheme, err := client.GetWorkflowScheme(id)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraWorkflowSchemeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workflow scheme",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		body := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			body["description"] = description
		}
		scheme, err := client.CreateWorkflowScheme(body)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraWorkflowSchemeUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a workflow scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid workflow scheme ID: %w", err)
		}
		body := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			body["name"], _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("description") {
			body["description"], _ = cmd.Flags().GetString("description")
		}
		scheme, err := client.UpdateWorkflowScheme(id, body)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraWorkflowSchemeDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a workflow scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid workflow scheme ID: %w", err)
		}
		if err := client.DeleteWorkflowScheme(id); err != nil {
			return err
		}
		fmt.Printf("Workflow scheme %d deleted.\n", id)
		return nil
	},
}

// ============================================================================
// Permission Schemes
// ============================================================================

var jiraPermissionSchemeCmd = &cobra.Command{
	Use:     "permissionscheme",
	Short:   "Manage permission schemes",
	Aliases: []string{"ps"},
	RunE:    helpRunE,
}

var jiraPermissionSchemeListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List permission schemes",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		schemes, err := client.GetPermissionSchemes()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME")
		for _, s := range schemes {
			fmt.Fprintf(w, "%d\t%s\n", s.ID, s.Name)
		}
		return w.Flush()
	},
}

var jiraPermissionSchemeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a permission scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid permission scheme ID: %w", err)
		}
		scheme, err := client.GetPermissionScheme(id)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraPermissionSchemeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a permission scheme",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		body := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			body["description"] = description
		}
		scheme, err := client.CreatePermissionScheme(body)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraPermissionSchemeDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a permission scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid permission scheme ID: %w", err)
		}
		if err := client.DeletePermissionScheme(id); err != nil {
			return err
		}
		fmt.Printf("Permission scheme %d deleted.\n", id)
		return nil
	},
}

// ============================================================================
// Notification Schemes
// ============================================================================

var jiraNotificationSchemeCmd = &cobra.Command{
	Use:     "notificationscheme",
	Short:   "Manage notification schemes",
	Aliases: []string{"ns"},
	RunE:    helpRunE,
}

var jiraNotificationSchemeListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List notification schemes",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		result, err := client.GetNotificationSchemes(startAt, maxResults)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME")
		for _, s := range result.Values {
			fmt.Fprintf(w, "%d\t%s\n", s.ID, s.Name)
		}
		return w.Flush()
	},
}

var jiraNotificationSchemeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a notification scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid notification scheme ID: %w", err)
		}
		scheme, err := client.GetNotificationScheme(id)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraNotificationSchemeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a notification scheme",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		body := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			body["description"] = description
		}
		scheme, err := client.CreateNotificationScheme(body)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraNotificationSchemeDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a notification scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid notification scheme ID: %w", err)
		}
		if err := client.DeleteNotificationScheme(id); err != nil {
			return err
		}
		fmt.Printf("Notification scheme %d deleted.\n", id)
		return nil
	},
}

// ============================================================================
// Issue Security Schemes
// ============================================================================

var jiraIssueSecuritySchemeCmd = &cobra.Command{
	Use:     "issuesecurityscheme",
	Short:   "Manage issue security schemes",
	Aliases: []string{"iss"},
	RunE:    helpRunE,
}

var jiraIssueSecuritySchemeListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List issue security schemes",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		schemes, err := client.GetIssueSecuritySchemes()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME")
		for _, s := range schemes {
			fmt.Fprintf(w, "%d\t%s\n", s.ID, s.Name)
		}
		return w.Flush()
	},
}

var jiraIssueSecuritySchemeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an issue security scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid issue security scheme ID: %w", err)
		}
		scheme, err := client.GetIssueSecurityScheme(id)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraIssueSecuritySchemeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an issue security scheme",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		body := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			body["description"] = description
		}
		scheme, err := client.CreateIssueSecurityScheme(body)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraIssueSecuritySchemeDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an issue security scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid issue security scheme ID: %w", err)
		}
		if err := client.DeleteIssueSecurityScheme(id); err != nil {
			return err
		}
		fmt.Printf("Issue security scheme %d deleted.\n", id)
		return nil
	},
}

// ============================================================================
// Field Configs
// ============================================================================

var jiraFieldConfigCmd = &cobra.Command{
	Use:     "fieldconfig",
	Short:   "Manage field configurations",
	Aliases: []string{"fc"},
	RunE:    helpRunE,
}

var jiraFieldConfigListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List field configurations",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		result, err := client.GetFieldConfigurations(startAt, maxResults)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tDEFAULT")
		for _, fc := range result.Values {
			fmt.Fprintf(w, "%d\t%s\t%v\n", fc.ID, fc.Name, fc.IsDefault)
		}
		return w.Flush()
	},
}

var jiraFieldConfigCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a field configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		body := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			body["description"] = description
		}
		fc, err := client.CreateFieldConfiguration(body)
		if err != nil {
			return err
		}
		return printJSON(fc)
	},
}

var jiraFieldConfigDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a field configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid field configuration ID: %w", err)
		}
		if err := client.DeleteFieldConfiguration(id); err != nil {
			return err
		}
		fmt.Printf("Field configuration %d deleted.\n", id)
		return nil
	},
}

// ============================================================================
// Issue Type Schemes
// ============================================================================

var jiraIssueTypeSchemeCmd = &cobra.Command{
	Use:     "issuetypescheme",
	Short:   "Manage issue type schemes",
	Aliases: []string{"its"},
	RunE:    helpRunE,
}

var jiraIssueTypeSchemeListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List issue type schemes",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		result, err := client.GetIssueTypeSchemes(startAt, maxResults)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME")
		for _, s := range result.Values {
			fmt.Fprintf(w, "%s\t%s\n", s.ID, s.Name)
		}
		return w.Flush()
	},
}

var jiraIssueTypeSchemeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an issue type scheme",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		body := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			body["description"] = description
		}
		scheme, err := client.CreateIssueTypeScheme(body)
		if err != nil {
			return err
		}
		return printJSON(scheme)
	},
}

var jiraIssueTypeSchemeDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an issue type scheme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteIssueTypeScheme(args[0]); err != nil {
			return err
		}
		fmt.Printf("Issue type scheme %s deleted.\n", args[0])
		return nil
	},
}

// ============================================================================
// Server Info
// ============================================================================

var jiraServerInfoCmd = &cobra.Command{
	Use:     "serverinfo",
	Short:   "Show Jira server information",
	Aliases: []string{"si"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		info, err := client.GetServerInfo()
		if err != nil {
			return err
		}
		return printJSON(info)
	},
}

// ============================================================================
// Webhooks
// ============================================================================

var jiraWebhookCmd = &cobra.Command{
	Use:     "webhook",
	Short:   "Manage webhooks",
	Aliases: []string{"wh"},
	RunE:    helpRunE,
}

var jiraWebhookListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List webhooks",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		result, err := client.GetWebhooks(startAt, maxResults)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tJQL\tEVENTS")
		for _, wh := range result.Values {
			events := strings.Join(wh.Events, ", ")
			fmt.Fprintf(w, "%d\t%s\t%s\n", wh.ID, wh.JqlFilter, events)
		}
		return w.Flush()
	},
}

// ============================================================================
// Attachments
// ============================================================================

var jiraAttachmentCmd = &cobra.Command{
	Use:     "attachment",
	Short:   "Manage attachments",
	Aliases: []string{"att"},
	RunE:    helpRunE,
}

var jiraAttachmentGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an attachment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		att, err := client.GetAttachment(args[0])
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(att)
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tFILENAME\tSIZE\tMIME TYPE\tCREATED")
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", att.ID, att.Filename, att.Size, att.MimeType, att.Created)
		return w.Flush()
	},
}

var jiraAttachmentDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an attachment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteAttachment(args[0]); err != nil {
			return err
		}
		fmt.Printf("Attachment %s deleted.\n", args[0])
		return nil
	},
}

var jiraAttachmentMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Show attachment settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		meta, err := client.GetAttachmentMeta()
		if err != nil {
			return err
		}
		return printJSON(meta)
	},
}

// ============================================================================
// Audit
// ============================================================================

var jiraAuditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Show audit records",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		records, err := client.GetAuditRecords(startAt, maxResults)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tSUMMARY\tCREATED\tCATEGORY")
		for _, r := range records.Records {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", r.ID, r.Summary, r.Created, r.Category)
		}
		return w.Flush()
	},
}

// ============================================================================
// Banner
// ============================================================================

var jiraBannerCmd = &cobra.Command{
	Use:   "banner",
	Short: "Manage announcement banner",
	RunE:  helpRunE,
}

var jiraBannerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the announcement banner",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		banner, err := client.GetAnnouncementBanner()
		if err != nil {
			return err
		}
		return printJSON(banner)
	},
}

var jiraBannerSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the announcement banner",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		banner := &jira.AnnouncementBanner{}
		if cmd.Flags().Changed("message") {
			banner.Message, _ = cmd.Flags().GetString("message")
		}
		if cmd.Flags().Changed("enabled") {
			banner.IsEnabled, _ = cmd.Flags().GetBool("enabled")
		}
		if cmd.Flags().Changed("dismissible") {
			banner.IsDismissible, _ = cmd.Flags().GetBool("dismissible")
		}
		if err := client.SetAnnouncementBanner(banner); err != nil {
			return err
		}
		fmt.Println("Announcement banner updated.")
		return nil
	},
}

// ============================================================================
// Configuration
// ============================================================================

var jiraConfigurationCmd = &cobra.Command{
	Use:     "configuration",
	Short:   "Show Jira configuration",
	Aliases: []string{"config"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		config, err := client.GetConfiguration()
		if err != nil {
			return err
		}
		return printJSON(config)
	},
}

// ============================================================================
// Permissions
// ============================================================================

var jiraPermissionCmd = &cobra.Command{
	Use:     "permission",
	Short:   "Manage permissions",
	Aliases: []string{"perm"},
	RunE:    helpRunE,
}

var jiraPermissionMineCmd = &cobra.Command{
	Use:   "mine",
	Short: "Show my permissions",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		project, _ := defaultProject(cmd)
		issue, _ := cmd.Flags().GetString("issue")
		perms, err := client.GetMyPermissions(project, issue)
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "KEY\tNAME\tHAVE_PERMISSION")
		for _, p := range perms {
			fmt.Fprintf(w, "%s\t%s\t%v\n", p.Key, p.Name, p.HavePermission)
		}
		return w.Flush()
	},
}

var jiraPermissionAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all permissions",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		perms, err := client.GetAllPermissions()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "KEY\tNAME\tHAVE_PERMISSION")
		for _, p := range perms {
			fmt.Fprintf(w, "%s\t%s\t%v\n", p.Key, p.Name, p.HavePermission)
		}
		return w.Flush()
	},
}

// ============================================================================
// Tasks
// ============================================================================

var jiraTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage async tasks",
	RunE:  helpRunE,
}

var jiraTaskGetCmd = &cobra.Command{
	Use:   "get <task-id>",
	Short: "Get a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		task, err := client.GetTask(args[0])
		if err != nil {
			return err
		}
		return printJSON(task)
	},
}

var jiraTaskCancelCmd = &cobra.Command{
	Use:   "cancel <task-id>",
	Short: "Cancel a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.CancelTask(args[0]); err != nil {
			return err
		}
		fmt.Printf("Task %s cancelled.\n", args[0])
		return nil
	},
}

// ============================================================================
// Project Categories
// ============================================================================

var jiraProjectCategoryCmd = &cobra.Command{
	Use:     "projectcategory",
	Short:   "Manage project categories",
	Aliases: []string{"pc"},
	RunE:    helpRunE,
}

var jiraProjectCategoryListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List project categories",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		cats, err := client.GetProjectCategories()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
		for _, c := range cats {
			fmt.Fprintf(w, "%s\t%s\t%s\n", c.ID, c.Name, c.Description)
		}
		return w.Flush()
	},
}

var jiraProjectCategoryGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a project category",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		cat, err := client.GetProjectCategory(args[0])
		if err != nil {
			return err
		}
		return printJSON(cat)
	},
}

var jiraProjectCategoryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project category",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		cat := &jira.ProjectCategory{
			Name: name,
		}
		if description != "" {
			cat.Description = description
		}
		result, err := client.CreateProjectCategory(cat)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var jiraProjectCategoryUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a project category",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		cat := &jira.ProjectCategory{}
		if cmd.Flags().Changed("name") {
			cat.Name, _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("description") {
			cat.Description, _ = cmd.Flags().GetString("description")
		}
		result, err := client.UpdateProjectCategory(args[0], cat)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var jiraProjectCategoryDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a project category",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteProjectCategory(args[0]); err != nil {
			return err
		}
		fmt.Printf("Project category %s deleted.\n", args[0])
		return nil
	},
}

// ============================================================================
// init - register all admin commands
// ============================================================================

func init() {
	// --- Roles ---
	jiraRoleCreateCmd.Flags().String("name", "", "Role name (required)")
	jiraRoleCreateCmd.Flags().String("description", "", "Role description")
	_ = jiraRoleCreateCmd.MarkFlagRequired("name")
	jiraRoleCmd.AddCommand(jiraRoleListCmd)
	jiraRoleCmd.AddCommand(jiraRoleGetCmd)
	jiraRoleCmd.AddCommand(jiraRoleCreateCmd)
	jiraRoleCmd.AddCommand(jiraRoleDeleteCmd)
	jiraCmd.AddCommand(jiraRoleCmd)

	// --- Issue Links ---
	jiraIssueLinkCreateCmd.Flags().String("inward-issue", "", "Inward issue key (required)")
	jiraIssueLinkCreateCmd.Flags().String("outward-issue", "", "Outward issue key (required)")
	jiraIssueLinkCreateCmd.Flags().String("type", "", "Link type name (required)")
	_ = jiraIssueLinkCreateCmd.MarkFlagRequired("inward-issue")
	_ = jiraIssueLinkCreateCmd.MarkFlagRequired("outward-issue")
	_ = jiraIssueLinkCreateCmd.MarkFlagRequired("type")
	jiraIssueLinkCmd.AddCommand(jiraIssueLinkCreateCmd)
	jiraIssueLinkCmd.AddCommand(jiraIssueLinkGetCmd)
	jiraIssueLinkCmd.AddCommand(jiraIssueLinkDeleteCmd)
	jiraCmd.AddCommand(jiraIssueLinkCmd)

	// --- Issue Link Types ---
	jiraIssueLinkTypeCreateCmd.Flags().String("name", "", "Link type name (required)")
	jiraIssueLinkTypeCreateCmd.Flags().String("inward", "", "Inward description (required)")
	jiraIssueLinkTypeCreateCmd.Flags().String("outward", "", "Outward description (required)")
	_ = jiraIssueLinkTypeCreateCmd.MarkFlagRequired("name")
	_ = jiraIssueLinkTypeCreateCmd.MarkFlagRequired("inward")
	_ = jiraIssueLinkTypeCreateCmd.MarkFlagRequired("outward")
	jiraIssueLinkTypeUpdateCmd.Flags().String("name", "", "Link type name")
	jiraIssueLinkTypeUpdateCmd.Flags().String("inward", "", "Inward description")
	jiraIssueLinkTypeUpdateCmd.Flags().String("outward", "", "Outward description")
	jiraIssueLinkTypeCmd.AddCommand(jiraIssueLinkTypeListCmd)
	jiraIssueLinkTypeCmd.AddCommand(jiraIssueLinkTypeGetCmd)
	jiraIssueLinkTypeCmd.AddCommand(jiraIssueLinkTypeCreateCmd)
	jiraIssueLinkTypeCmd.AddCommand(jiraIssueLinkTypeUpdateCmd)
	jiraIssueLinkTypeCmd.AddCommand(jiraIssueLinkTypeDeleteCmd)
	jiraCmd.AddCommand(jiraIssueLinkTypeCmd)

	// --- Screens ---
	jiraScreenListCmd.Flags().Int("start-at", 0, "Start index")
	jiraScreenListCmd.Flags().Int("max-results", 50, "Max results")
	jiraScreenCreateCmd.Flags().String("name", "", "Screen name (required)")
	jiraScreenCreateCmd.Flags().String("description", "", "Screen description")
	_ = jiraScreenCreateCmd.MarkFlagRequired("name")
	jiraScreenCmd.AddCommand(jiraScreenListCmd)
	jiraScreenCmd.AddCommand(jiraScreenCreateCmd)
	jiraScreenCmd.AddCommand(jiraScreenDeleteCmd)
	jiraScreenCmd.AddCommand(jiraScreenTabsCmd)
	jiraScreenCmd.AddCommand(jiraScreenTabFieldsCmd)
	jiraCmd.AddCommand(jiraScreenCmd)

	// --- Workflows ---
	jiraWorkflowCmd.AddCommand(jiraWorkflowListCmd)
	jiraCmd.AddCommand(jiraWorkflowCmd)

	// --- Workflow Schemes ---
	jiraWorkflowSchemeListCmd.Flags().Int("start-at", 0, "Start index")
	jiraWorkflowSchemeListCmd.Flags().Int("max-results", 50, "Max results")
	jiraWorkflowSchemeCreateCmd.Flags().String("name", "", "Scheme name (required)")
	jiraWorkflowSchemeCreateCmd.Flags().String("description", "", "Scheme description")
	_ = jiraWorkflowSchemeCreateCmd.MarkFlagRequired("name")
	jiraWorkflowSchemeUpdateCmd.Flags().String("name", "", "Scheme name")
	jiraWorkflowSchemeUpdateCmd.Flags().String("description", "", "Scheme description")
	jiraWorkflowSchemeCmd.AddCommand(jiraWorkflowSchemeListCmd)
	jiraWorkflowSchemeCmd.AddCommand(jiraWorkflowSchemeGetCmd)
	jiraWorkflowSchemeCmd.AddCommand(jiraWorkflowSchemeCreateCmd)
	jiraWorkflowSchemeCmd.AddCommand(jiraWorkflowSchemeUpdateCmd)
	jiraWorkflowSchemeCmd.AddCommand(jiraWorkflowSchemeDeleteCmd)
	jiraCmd.AddCommand(jiraWorkflowSchemeCmd)

	// --- Permission Schemes ---
	jiraPermissionSchemeCreateCmd.Flags().String("name", "", "Scheme name (required)")
	jiraPermissionSchemeCreateCmd.Flags().String("description", "", "Scheme description")
	_ = jiraPermissionSchemeCreateCmd.MarkFlagRequired("name")
	jiraPermissionSchemeCmd.AddCommand(jiraPermissionSchemeListCmd)
	jiraPermissionSchemeCmd.AddCommand(jiraPermissionSchemeGetCmd)
	jiraPermissionSchemeCmd.AddCommand(jiraPermissionSchemeCreateCmd)
	jiraPermissionSchemeCmd.AddCommand(jiraPermissionSchemeDeleteCmd)
	jiraCmd.AddCommand(jiraPermissionSchemeCmd)

	// --- Notification Schemes ---
	jiraNotificationSchemeListCmd.Flags().Int("start-at", 0, "Start index")
	jiraNotificationSchemeListCmd.Flags().Int("max-results", 50, "Max results")
	jiraNotificationSchemeCreateCmd.Flags().String("name", "", "Scheme name (required)")
	jiraNotificationSchemeCreateCmd.Flags().String("description", "", "Scheme description")
	_ = jiraNotificationSchemeCreateCmd.MarkFlagRequired("name")
	jiraNotificationSchemeCmd.AddCommand(jiraNotificationSchemeListCmd)
	jiraNotificationSchemeCmd.AddCommand(jiraNotificationSchemeGetCmd)
	jiraNotificationSchemeCmd.AddCommand(jiraNotificationSchemeCreateCmd)
	jiraNotificationSchemeCmd.AddCommand(jiraNotificationSchemeDeleteCmd)
	jiraCmd.AddCommand(jiraNotificationSchemeCmd)

	// --- Issue Security Schemes ---
	jiraIssueSecuritySchemeCreateCmd.Flags().String("name", "", "Scheme name (required)")
	jiraIssueSecuritySchemeCreateCmd.Flags().String("description", "", "Scheme description")
	_ = jiraIssueSecuritySchemeCreateCmd.MarkFlagRequired("name")
	jiraIssueSecuritySchemeCmd.AddCommand(jiraIssueSecuritySchemeListCmd)
	jiraIssueSecuritySchemeCmd.AddCommand(jiraIssueSecuritySchemeGetCmd)
	jiraIssueSecuritySchemeCmd.AddCommand(jiraIssueSecuritySchemeCreateCmd)
	jiraIssueSecuritySchemeCmd.AddCommand(jiraIssueSecuritySchemeDeleteCmd)
	jiraCmd.AddCommand(jiraIssueSecuritySchemeCmd)

	// --- Field Configs ---
	jiraFieldConfigListCmd.Flags().Int("start-at", 0, "Start index")
	jiraFieldConfigListCmd.Flags().Int("max-results", 50, "Max results")
	jiraFieldConfigCreateCmd.Flags().String("name", "", "Configuration name (required)")
	jiraFieldConfigCreateCmd.Flags().String("description", "", "Configuration description")
	_ = jiraFieldConfigCreateCmd.MarkFlagRequired("name")
	jiraFieldConfigCmd.AddCommand(jiraFieldConfigListCmd)
	jiraFieldConfigCmd.AddCommand(jiraFieldConfigCreateCmd)
	jiraFieldConfigCmd.AddCommand(jiraFieldConfigDeleteCmd)
	jiraCmd.AddCommand(jiraFieldConfigCmd)

	// --- Issue Type Schemes ---
	jiraIssueTypeSchemeListCmd.Flags().Int("start-at", 0, "Start index")
	jiraIssueTypeSchemeListCmd.Flags().Int("max-results", 50, "Max results")
	jiraIssueTypeSchemeCreateCmd.Flags().String("name", "", "Scheme name (required)")
	jiraIssueTypeSchemeCreateCmd.Flags().String("description", "", "Scheme description")
	_ = jiraIssueTypeSchemeCreateCmd.MarkFlagRequired("name")
	jiraIssueTypeSchemeCmd.AddCommand(jiraIssueTypeSchemeListCmd)
	jiraIssueTypeSchemeCmd.AddCommand(jiraIssueTypeSchemeCreateCmd)
	jiraIssueTypeSchemeCmd.AddCommand(jiraIssueTypeSchemeDeleteCmd)
	jiraCmd.AddCommand(jiraIssueTypeSchemeCmd)

	// --- Server Info ---
	jiraCmd.AddCommand(jiraServerInfoCmd)

	// --- Webhooks ---
	jiraWebhookListCmd.Flags().Int("start-at", 0, "Start index")
	jiraWebhookListCmd.Flags().Int("max-results", 50, "Max results")
	jiraWebhookCmd.AddCommand(jiraWebhookListCmd)
	jiraCmd.AddCommand(jiraWebhookCmd)

	// --- Attachments ---
	jiraAttachmentGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraAttachmentCmd.AddCommand(jiraAttachmentGetCmd)
	jiraAttachmentCmd.AddCommand(jiraAttachmentDeleteCmd)
	jiraAttachmentCmd.AddCommand(jiraAttachmentMetaCmd)
	jiraCmd.AddCommand(jiraAttachmentCmd)

	// --- Audit ---
	jiraAuditCmd.Flags().Int("start-at", 0, "Start index")
	jiraAuditCmd.Flags().Int("max-results", 50, "Max results")
	jiraCmd.AddCommand(jiraAuditCmd)

	// --- Banner ---
	jiraBannerSetCmd.Flags().String("message", "", "Banner message")
	jiraBannerSetCmd.Flags().Bool("enabled", false, "Enable the banner")
	jiraBannerSetCmd.Flags().Bool("dismissible", false, "Allow dismissing the banner")
	jiraBannerCmd.AddCommand(jiraBannerGetCmd)
	jiraBannerCmd.AddCommand(jiraBannerSetCmd)
	jiraCmd.AddCommand(jiraBannerCmd)

	// --- Configuration ---
	jiraCmd.AddCommand(jiraConfigurationCmd)

	// --- Permissions ---
	jiraPermissionMineCmd.Flags().String("project", "", "Project key (uses profile default if not set)")
	jiraPermissionMineCmd.Flags().String("issue", "", "Issue key")
	jiraPermissionCmd.AddCommand(jiraPermissionMineCmd)
	jiraPermissionCmd.AddCommand(jiraPermissionAllCmd)
	jiraCmd.AddCommand(jiraPermissionCmd)

	// --- Tasks ---
	jiraTaskCmd.AddCommand(jiraTaskGetCmd)
	jiraTaskCmd.AddCommand(jiraTaskCancelCmd)
	jiraCmd.AddCommand(jiraTaskCmd)

	// --- Project Categories ---
	jiraProjectCategoryCreateCmd.Flags().String("name", "", "Category name (required)")
	jiraProjectCategoryCreateCmd.Flags().String("description", "", "Category description")
	_ = jiraProjectCategoryCreateCmd.MarkFlagRequired("name")
	jiraProjectCategoryUpdateCmd.Flags().String("name", "", "Category name")
	jiraProjectCategoryUpdateCmd.Flags().String("description", "", "Category description")
	jiraProjectCategoryCmd.AddCommand(jiraProjectCategoryListCmd)
	jiraProjectCategoryCmd.AddCommand(jiraProjectCategoryGetCmd)
	jiraProjectCategoryCmd.AddCommand(jiraProjectCategoryCreateCmd)
	jiraProjectCategoryCmd.AddCommand(jiraProjectCategoryUpdateCmd)
	jiraProjectCategoryCmd.AddCommand(jiraProjectCategoryDeleteCmd)
	jiraCmd.AddCommand(jiraProjectCategoryCmd)
}
