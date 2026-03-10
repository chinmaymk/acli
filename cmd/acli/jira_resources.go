package acli

import (
	"fmt"

	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

// ============================================================================
// Components
// ============================================================================

var jiraComponentCmd = &cobra.Command{
	Use:     "component",
	Short:   "Manage project components",
	Aliases: []string{"comp"},
	RunE:    helpRunE,
}

var jiraComponentGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a component by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		comp, err := client.GetComponent(args[0])
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(comp)
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\t%s\n", comp.ID)
		fmt.Fprintf(w, "Name\t%s\n", comp.Name)
		fmt.Fprintf(w, "Description\t%s\n", comp.Description)
		fmt.Fprintf(w, "Project\t%s\n", comp.Project)
		fmt.Fprintf(w, "Assignee Type\t%s\n", comp.AssigneeType)
		if comp.Lead != nil {
			fmt.Fprintf(w, "Lead\t%s\n", comp.Lead.DisplayName)
		}
		w.Flush()
		return nil
	},
}

var jiraComponentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a component",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		project, _ := defaultProject(cmd)
		name, _ := cmd.Flags().GetString("name")
		body := map[string]interface{}{
			"project": project,
			"name":    name,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		if v, _ := cmd.Flags().GetString("lead"); v != "" {
			body["leadAccountId"] = v
		}
		if v, _ := cmd.Flags().GetString("assignee-type"); v != "" {
			body["assigneeType"] = v
		}
		comp, err := client.CreateComponent(body)
		if err != nil {
			return err
		}
		fmt.Printf("Component created: %s (ID: %s)\n", comp.Name, comp.ID)
		return nil
	},
}

var jiraComponentUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a component",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		body := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		if v, _ := cmd.Flags().GetString("lead"); v != "" {
			body["leadAccountId"] = v
		}
		comp, err := client.UpdateComponent(args[0], body)
		if err != nil {
			return err
		}
		fmt.Printf("Component updated: %s (ID: %s)\n", comp.Name, comp.ID)
		return nil
	},
}

var jiraComponentDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a component",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteComponent(args[0]); err != nil {
			return err
		}
		fmt.Printf("Component %s deleted\n", args[0])
		return nil
	},
}

// ============================================================================
// Versions
// ============================================================================

var jiraVersionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Manage project versions",
	Aliases: []string{"ver"},
	RunE:    helpRunE,
}

var jiraVersionGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a version by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		ver, err := client.GetVersion(args[0])
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(ver)
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\t%s\n", ver.ID)
		fmt.Fprintf(w, "Name\t%s\n", ver.Name)
		fmt.Fprintf(w, "Description\t%s\n", ver.Description)
		fmt.Fprintf(w, "Released\t%v\n", ver.Released)
		fmt.Fprintf(w, "Archived\t%v\n", ver.Archived)
		fmt.Fprintf(w, "Start Date\t%s\n", ver.StartDate)
		fmt.Fprintf(w, "Release Date\t%s\n", ver.ReleaseDate)
		w.Flush()
		return nil
	},
}

var jiraVersionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a version",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		projectID, _ := cmd.Flags().GetInt("project-id")
		name, _ := cmd.Flags().GetString("name")
		ver := &jira.Version{
			ProjectID: projectID,
			Name:      name,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			ver.Description = v
		}
		if v, _ := cmd.Flags().GetString("start-date"); v != "" {
			ver.StartDate = v
		}
		if v, _ := cmd.Flags().GetString("release-date"); v != "" {
			ver.ReleaseDate = v
		}
		if cmd.Flags().Changed("released") {
			r, _ := cmd.Flags().GetBool("released")
			ver.Released = r
		}
		if cmd.Flags().Changed("archived") {
			a, _ := cmd.Flags().GetBool("archived")
			ver.Archived = a
		}
		result, err := client.CreateVersion(ver)
		if err != nil {
			return err
		}
		fmt.Printf("Version created: %s (ID: %s)\n", result.Name, result.ID)
		return nil
	},
}

var jiraVersionUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		ver, err := client.GetVersion(args[0])
		if err != nil {
			return err
		}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			ver.Name = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			ver.Description = v
		}
		if v, _ := cmd.Flags().GetString("start-date"); v != "" {
			ver.StartDate = v
		}
		if v, _ := cmd.Flags().GetString("release-date"); v != "" {
			ver.ReleaseDate = v
		}
		if cmd.Flags().Changed("released") {
			r, _ := cmd.Flags().GetBool("released")
			ver.Released = r
		}
		if cmd.Flags().Changed("archived") {
			a, _ := cmd.Flags().GetBool("archived")
			ver.Archived = a
		}
		result, err := client.UpdateVersion(args[0], ver)
		if err != nil {
			return err
		}
		fmt.Printf("Version updated: %s (ID: %s)\n", result.Name, result.ID)
		return nil
	},
}

var jiraVersionDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteVersion(args[0]); err != nil {
			return err
		}
		fmt.Printf("Version %s deleted\n", args[0])
		return nil
	},
}

// ============================================================================
// Fields
// ============================================================================

var jiraFieldCmd = &cobra.Command{
	Use:   "field",
	Short: "Manage fields",
	RunE:  helpRunE,
}

var jiraFieldListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all fields",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		fields, err := client.GetFields()
		if err != nil {
			return err
		}
		customOnly, _ := cmd.Flags().GetBool("custom")
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if customOnly {
			var filtered []interface{}
			for i := range fields {
				if fields[i].Custom {
					filtered = append(filtered, fields[i])
				}
			}
			if jsonFlag {
				return printJSON(filtered)
			}
		} else if jsonFlag {
			return printJSON(fields)
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\tNAME\tTYPE\tCUSTOM\n")
		for _, f := range fields {
			if customOnly && !f.Custom {
				continue
			}
			fieldType := ""
			if f.Schema != nil {
				fieldType = f.Schema.Type
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", f.ID, f.Name, fieldType, f.Custom)
		}
		w.Flush()
		return nil
	},
}

var jiraFieldCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a custom field",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		fieldType, _ := cmd.Flags().GetString("type")
		body := map[string]interface{}{
			"name": name,
			"type": fieldType,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		if v, _ := cmd.Flags().GetString("search-key"); v != "" {
			body["searcherKey"] = v
		}
		field, err := client.CreateCustomField(body)
		if err != nil {
			return err
		}
		fmt.Printf("Field created: %s (ID: %s)\n", field.Name, field.ID)
		return nil
	},
}

var jiraFieldDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a custom field",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteCustomField(args[0]); err != nil {
			return err
		}
		fmt.Printf("Field %s deleted\n", args[0])
		return nil
	},
}

var jiraFieldTrashCmd = &cobra.Command{
	Use:   "trash <id>",
	Short: "Move a custom field to trash",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.TrashCustomField(args[0]); err != nil {
			return err
		}
		fmt.Printf("Field %s moved to trash\n", args[0])
		return nil
	},
}

var jiraFieldRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "Restore a custom field from trash",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.RestoreCustomField(args[0]); err != nil {
			return err
		}
		fmt.Printf("Field %s restored\n", args[0])
		return nil
	},
}

// ============================================================================
// Labels
// ============================================================================

var jiraLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "List labels",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		result, err := client.GetLabels(startAt, maxResults)
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(result)
		}
		w := newTabWriter()
		fmt.Fprintf(w, "LABEL\n")
		for _, label := range result.Values {
			fmt.Fprintf(w, "%s\n", label)
		}
		w.Flush()
		return nil
	},
}

// ============================================================================
// Issue Types
// ============================================================================

var jiraIssuetypeCmd = &cobra.Command{
	Use:     "issuetype",
	Short:   "Manage issue types",
	Aliases: []string{"it"},
	RunE:    helpRunE,
}

var jiraIssuetypeListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all issue types",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		types, err := client.GetAllIssueTypes()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\tNAME\tSUBTASK\tDESCRIPTION\n")
		for _, t := range types {
			fmt.Fprintf(w, "%s\t%s\t%v\t%s\n", t.ID, t.Name, t.Subtask, t.Description)
		}
		w.Flush()
		return nil
	},
}

var jiraIssuetypeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an issue type by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		it, err := client.GetIssueType(args[0])
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(it)
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\t%s\n", it.ID)
		fmt.Fprintf(w, "Name\t%s\n", it.Name)
		fmt.Fprintf(w, "Subtask\t%v\n", it.Subtask)
		fmt.Fprintf(w, "Description\t%s\n", it.Description)
		w.Flush()
		return nil
	},
}

var jiraIssuetypeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an issue type",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		body := map[string]interface{}{
			"name": name,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		itType, _ := cmd.Flags().GetString("type")
		body["type"] = itType
		it, err := client.CreateIssueType(body)
		if err != nil {
			return err
		}
		fmt.Printf("Issue type created: %s (ID: %s)\n", it.Name, it.ID)
		return nil
	},
}

var jiraIssuetypeUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an issue type",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		body := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		it, err := client.UpdateIssueType(args[0], body)
		if err != nil {
			return err
		}
		fmt.Printf("Issue type updated: %s (ID: %s)\n", it.Name, it.ID)
		return nil
	},
}

var jiraIssuetypeDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an issue type",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteIssueType(args[0]); err != nil {
			return err
		}
		fmt.Printf("Issue type %s deleted\n", args[0])
		return nil
	},
}

// ============================================================================
// Priorities
// ============================================================================

var jiraPriorityCmd = &cobra.Command{
	Use:     "priority",
	Short:   "Manage priorities",
	Aliases: []string{"pri"},
	RunE:    helpRunE,
}

var jiraPriorityListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all priorities",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		priorities, err := client.GetAllPriorities()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\tNAME\tDESCRIPTION\n")
		for _, p := range priorities {
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.ID, p.Name, p.Description)
		}
		w.Flush()
		return nil
	},
}

var jiraPriorityGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a priority by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		p, err := client.GetPriority(args[0])
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(p)
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\t%s\n", p.ID)
		fmt.Fprintf(w, "Name\t%s\n", p.Name)
		fmt.Fprintf(w, "Description\t%s\n", p.Description)
		fmt.Fprintf(w, "Status Color\t%s\n", p.StatusColor)
		w.Flush()
		return nil
	},
}

var jiraPriorityCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a priority",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		statusColor, _ := cmd.Flags().GetString("status-color")
		body := map[string]interface{}{
			"name":        name,
			"statusColor": statusColor,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		p, err := client.CreatePriority(body)
		if err != nil {
			return err
		}
		fmt.Printf("Priority created: %s (ID: %s)\n", p.Name, p.ID)
		return nil
	},
}

var jiraPriorityUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a priority",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		body := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		p, err := client.UpdatePriority(args[0], body)
		if err != nil {
			return err
		}
		fmt.Printf("Priority updated: %s (ID: %s)\n", p.Name, p.ID)
		return nil
	},
}

var jiraPriorityDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a priority",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeletePriority(args[0]); err != nil {
			return err
		}
		fmt.Printf("Priority %s deleted\n", args[0])
		return nil
	},
}

// ============================================================================
// Resolutions
// ============================================================================

var jiraResolutionCmd = &cobra.Command{
	Use:     "resolution",
	Short:   "Manage resolutions",
	Aliases: []string{"res"},
	RunE:    helpRunE,
}

var jiraResolutionListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all resolutions",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		resolutions, err := client.GetAllResolutions()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\tNAME\tDESCRIPTION\n")
		for _, r := range resolutions {
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.ID, r.Name, r.Description)
		}
		w.Flush()
		return nil
	},
}

var jiraResolutionGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a resolution by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		r, err := client.GetResolution(args[0])
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(r)
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\t%s\n", r.ID)
		fmt.Fprintf(w, "Name\t%s\n", r.Name)
		fmt.Fprintf(w, "Description\t%s\n", r.Description)
		w.Flush()
		return nil
	},
}

var jiraResolutionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resolution",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		body := map[string]interface{}{
			"name": name,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		r, err := client.CreateResolution(body)
		if err != nil {
			return err
		}
		fmt.Printf("Resolution created: %s (ID: %s)\n", r.Name, r.ID)
		return nil
	},
}

var jiraResolutionUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a resolution",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		body := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		r, err := client.UpdateResolution(args[0], body)
		if err != nil {
			return err
		}
		fmt.Printf("Resolution updated: %s (ID: %s)\n", r.Name, r.ID)
		return nil
	},
}

var jiraResolutionDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a resolution",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteResolution(args[0]); err != nil {
			return err
		}
		fmt.Printf("Resolution %s deleted\n", args[0])
		return nil
	},
}

// ============================================================================
// Statuses
// ============================================================================

var jiraStatusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Manage statuses",
	Aliases: []string{"st"},
	RunE:    helpRunE,
}

var jiraStatusListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all statuses",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		statuses, err := client.GetAllStatuses()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\tNAME\tCATEGORY\n")
		for _, s := range statuses {
			fmt.Fprintf(w, "%s\t%s\t%s\n", s.ID, s.Name, s.StatusCategory.Name)
		}
		w.Flush()
		return nil
	},
}

var jiraStatusGetCmd = &cobra.Command{
	Use:   "get <id-or-name>",
	Short: "Get a status by ID or name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		s, err := client.GetStatus(args[0])
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(s)
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\t%s\n", s.ID)
		fmt.Fprintf(w, "Name\t%s\n", s.Name)
		fmt.Fprintf(w, "Description\t%s\n", s.Description)
		fmt.Fprintf(w, "Category\t%s\n", s.StatusCategory.Name)
		w.Flush()
		return nil
	},
}

var jiraStatusCategoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List status categories",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		categories, err := client.GetStatusCategories()
		if err != nil {
			return err
		}
		w := newTabWriter()
		fmt.Fprintf(w, "ID\tKEY\tNAME\tCOLOR\n")
		for _, c := range categories {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", c.ID, c.Key, c.Name, c.ColorName)
		}
		w.Flush()
		return nil
	},
}

// ============================================================================
// init - register all subcommands and flags
// ============================================================================

func init() {
	// --- Components ---
	jiraCmd.AddCommand(jiraComponentCmd)
	jiraComponentCmd.AddCommand(jiraComponentGetCmd)
	jiraComponentCmd.AddCommand(jiraComponentCreateCmd)
	jiraComponentCmd.AddCommand(jiraComponentUpdateCmd)
	jiraComponentCmd.AddCommand(jiraComponentDeleteCmd)

	jiraComponentGetCmd.Flags().Bool("json", false, "Output as JSON")

	jiraComponentCreateCmd.Flags().String("project", "", "Project key (uses profile default if not set)")
	jiraComponentCreateCmd.Flags().String("name", "", "Component name (required)")
	jiraComponentCreateCmd.Flags().String("description", "", "Component description")
	jiraComponentCreateCmd.Flags().String("lead", "", "Lead account ID")
	jiraComponentCreateCmd.Flags().String("assignee-type", "", "Assignee type")
	_ = jiraComponentCreateCmd.MarkFlagRequired("name")

	jiraComponentUpdateCmd.Flags().String("name", "", "Component name")
	jiraComponentUpdateCmd.Flags().String("description", "", "Component description")
	jiraComponentUpdateCmd.Flags().String("lead", "", "Lead account ID")

	// --- Versions ---
	jiraCmd.AddCommand(jiraVersionCmd)
	jiraVersionCmd.AddCommand(jiraVersionGetCmd)
	jiraVersionCmd.AddCommand(jiraVersionCreateCmd)
	jiraVersionCmd.AddCommand(jiraVersionUpdateCmd)
	jiraVersionCmd.AddCommand(jiraVersionDeleteCmd)

	jiraVersionGetCmd.Flags().Bool("json", false, "Output as JSON")

	jiraVersionCreateCmd.Flags().Int("project-id", 0, "Project ID (required)")
	jiraVersionCreateCmd.Flags().String("name", "", "Version name (required)")
	jiraVersionCreateCmd.Flags().String("description", "", "Version description")
	jiraVersionCreateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	jiraVersionCreateCmd.Flags().String("release-date", "", "Release date (YYYY-MM-DD)")
	jiraVersionCreateCmd.Flags().Bool("released", false, "Whether the version is released")
	jiraVersionCreateCmd.Flags().Bool("archived", false, "Whether the version is archived")
	_ = jiraVersionCreateCmd.MarkFlagRequired("project-id")
	_ = jiraVersionCreateCmd.MarkFlagRequired("name")

	jiraVersionUpdateCmd.Flags().String("name", "", "Version name")
	jiraVersionUpdateCmd.Flags().String("description", "", "Version description")
	jiraVersionUpdateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	jiraVersionUpdateCmd.Flags().String("release-date", "", "Release date (YYYY-MM-DD)")
	jiraVersionUpdateCmd.Flags().Bool("released", false, "Whether the version is released")
	jiraVersionUpdateCmd.Flags().Bool("archived", false, "Whether the version is archived")

	// --- Fields ---
	jiraCmd.AddCommand(jiraFieldCmd)
	jiraFieldCmd.AddCommand(jiraFieldListCmd)
	jiraFieldCmd.AddCommand(jiraFieldCreateCmd)
	jiraFieldCmd.AddCommand(jiraFieldDeleteCmd)
	jiraFieldCmd.AddCommand(jiraFieldTrashCmd)
	jiraFieldCmd.AddCommand(jiraFieldRestoreCmd)

	jiraFieldListCmd.Flags().Bool("json", false, "Output as JSON")
	jiraFieldListCmd.Flags().Bool("custom", false, "Show only custom fields")

	jiraFieldCreateCmd.Flags().String("name", "", "Field name (required)")
	jiraFieldCreateCmd.Flags().String("type", "", "Field type (required)")
	jiraFieldCreateCmd.Flags().String("description", "", "Field description")
	jiraFieldCreateCmd.Flags().String("search-key", "", "Searcher key")
	_ = jiraFieldCreateCmd.MarkFlagRequired("name")
	_ = jiraFieldCreateCmd.MarkFlagRequired("type")

	// --- Labels ---
	jiraCmd.AddCommand(jiraLabelCmd)
	jiraLabelCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraLabelCmd.Flags().Int("start-at", 0, "Index of the first result")
	jiraLabelCmd.Flags().Bool("json", false, "Output as JSON")

	// --- Issue Types ---
	jiraCmd.AddCommand(jiraIssuetypeCmd)
	jiraIssuetypeCmd.AddCommand(jiraIssuetypeListCmd)
	jiraIssuetypeCmd.AddCommand(jiraIssuetypeGetCmd)
	jiraIssuetypeCmd.AddCommand(jiraIssuetypeCreateCmd)
	jiraIssuetypeCmd.AddCommand(jiraIssuetypeUpdateCmd)
	jiraIssuetypeCmd.AddCommand(jiraIssuetypeDeleteCmd)

	jiraIssuetypeGetCmd.Flags().Bool("json", false, "Output as JSON")

	jiraIssuetypeCreateCmd.Flags().String("name", "", "Issue type name (required)")
	jiraIssuetypeCreateCmd.Flags().String("description", "", "Issue type description")
	jiraIssuetypeCreateCmd.Flags().String("type", "standard", "Issue type: standard or subtask")
	_ = jiraIssuetypeCreateCmd.MarkFlagRequired("name")

	jiraIssuetypeUpdateCmd.Flags().String("name", "", "Issue type name")
	jiraIssuetypeUpdateCmd.Flags().String("description", "", "Issue type description")

	// --- Priorities ---
	jiraCmd.AddCommand(jiraPriorityCmd)
	jiraPriorityCmd.AddCommand(jiraPriorityListCmd)
	jiraPriorityCmd.AddCommand(jiraPriorityGetCmd)
	jiraPriorityCmd.AddCommand(jiraPriorityCreateCmd)
	jiraPriorityCmd.AddCommand(jiraPriorityUpdateCmd)
	jiraPriorityCmd.AddCommand(jiraPriorityDeleteCmd)

	jiraPriorityGetCmd.Flags().Bool("json", false, "Output as JSON")

	jiraPriorityCreateCmd.Flags().String("name", "", "Priority name (required)")
	jiraPriorityCreateCmd.Flags().String("description", "", "Priority description")
	jiraPriorityCreateCmd.Flags().String("status-color", "#ffffff", "Status color hex")
	_ = jiraPriorityCreateCmd.MarkFlagRequired("name")

	jiraPriorityUpdateCmd.Flags().String("name", "", "Priority name")
	jiraPriorityUpdateCmd.Flags().String("description", "", "Priority description")

	// --- Resolutions ---
	jiraCmd.AddCommand(jiraResolutionCmd)
	jiraResolutionCmd.AddCommand(jiraResolutionListCmd)
	jiraResolutionCmd.AddCommand(jiraResolutionGetCmd)
	jiraResolutionCmd.AddCommand(jiraResolutionCreateCmd)
	jiraResolutionCmd.AddCommand(jiraResolutionUpdateCmd)
	jiraResolutionCmd.AddCommand(jiraResolutionDeleteCmd)

	jiraResolutionGetCmd.Flags().Bool("json", false, "Output as JSON")

	jiraResolutionCreateCmd.Flags().String("name", "", "Resolution name (required)")
	jiraResolutionCreateCmd.Flags().String("description", "", "Resolution description")
	_ = jiraResolutionCreateCmd.MarkFlagRequired("name")

	jiraResolutionUpdateCmd.Flags().String("name", "", "Resolution name")
	jiraResolutionUpdateCmd.Flags().String("description", "", "Resolution description")

	// --- Statuses ---
	jiraCmd.AddCommand(jiraStatusCmd)
	jiraStatusCmd.AddCommand(jiraStatusListCmd)
	jiraStatusCmd.AddCommand(jiraStatusGetCmd)
	jiraStatusCmd.AddCommand(jiraStatusCategoriesCmd)

	jiraStatusGetCmd.Flags().Bool("json", false, "Output as JSON")
}
