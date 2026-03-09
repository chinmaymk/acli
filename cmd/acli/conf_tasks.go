package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// task list
	listTasksCmd := &cobra.Command{
		Use:     "list",
		Short:   "List tasks",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if s := getStringFlag(cmd, "task-status"); s != "" {
				q.Set("status", s)
			}
			if ids := getStringSliceFlag(cmd, "task-id"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("task-id", id)
				}
			}
			if ids := getStringSliceFlag(cmd, "space-id"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("space-id", id)
				}
			}
			if ids := getStringSliceFlag(cmd, "page-id"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("page-id", id)
				}
			}
			if ids := getStringSliceFlag(cmd, "blogpost-id"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("blogpost-id", id)
				}
			}
			if ids := getStringSliceFlag(cmd, "created-by"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("created-by", id)
				}
			}
			if ids := getStringSliceFlag(cmd, "assigned-to"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("assigned-to", id)
				}
			}
			if ids := getStringSliceFlag(cmd, "completed-by"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("completed-by", id)
				}
			}
			if getBoolFlag(cmd, "include-blank-tasks") {
				q.Set("include-blank-tasks", "true")
			}
			for _, f := range []string{"created-at-from", "created-at-to", "due-at-from", "due-at-to", "completed-at-from", "completed-at-to"} {
				if v := getStringFlag(cmd, f); v != "" {
					q.Set(f, v)
				}
			}
			data, err := confGet(cmd, "/tasks", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(listTasksCmd)
	addBodyFormatFlag(listTasksCmd)
	listTasksCmd.Flags().String("task-status", "", "Filter by task status (complete, incomplete)")
	listTasksCmd.Flags().StringSlice("task-id", nil, "Filter by task IDs")
	listTasksCmd.Flags().StringSlice("space-id", nil, "Filter by space IDs")
	listTasksCmd.Flags().StringSlice("page-id", nil, "Filter by page IDs")
	listTasksCmd.Flags().StringSlice("blogpost-id", nil, "Filter by blog post IDs")
	listTasksCmd.Flags().StringSlice("created-by", nil, "Filter by creator account IDs")
	listTasksCmd.Flags().StringSlice("assigned-to", nil, "Filter by assignee account IDs")
	listTasksCmd.Flags().StringSlice("completed-by", nil, "Filter by completer account IDs")
	listTasksCmd.Flags().Bool("include-blank-tasks", false, "Include blank tasks")
	listTasksCmd.Flags().String("created-at-from", "", "Filter by creation date start (epoch ms)")
	listTasksCmd.Flags().String("created-at-to", "", "Filter by creation date end (epoch ms)")
	listTasksCmd.Flags().String("due-at-from", "", "Filter by due date start (epoch ms)")
	listTasksCmd.Flags().String("due-at-to", "", "Filter by due date end (epoch ms)")
	listTasksCmd.Flags().String("completed-at-from", "", "Filter by completion date start (epoch ms)")
	listTasksCmd.Flags().String("completed-at-to", "", "Filter by completion date end (epoch ms)")
	confTaskCmd.AddCommand(listTasksCmd)

	// task get
	getTaskCmd := &cobra.Command{
		Use:   "get [task-id]",
		Short: "Get task by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/tasks/"+args[0], q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addBodyFormatFlag(getTaskCmd)
	confTaskCmd.AddCommand(getTaskCmd)

	// task update
	updateTaskCmd := &cobra.Command{
		Use:   "update [task-id]",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"status": getStringFlag(cmd, "task-status"),
			}
			if id := getStringFlag(cmd, "task-id-field"); id != "" {
				body["id"] = id
			}
			if a := getStringFlag(cmd, "assigned-to"); a != "" {
				body["assignedTo"] = a
			}
			if d := getStringFlag(cmd, "due-at"); d != "" {
				body["dueAt"] = d
			}
			q := getPaginationQuery(cmd)
			data, err := confPut(cmd, "/tasks/"+args[0], q, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addBodyFormatFlag(updateTaskCmd)
	updateTaskCmd.Flags().String("task-status", "", "Task status (complete, incomplete) (required)")
	updateTaskCmd.Flags().String("task-id-field", "", "Task ID field")
	updateTaskCmd.Flags().String("assigned-to", "", "Assignee account ID")
	updateTaskCmd.Flags().String("due-at", "", "Due date")
	_ = updateTaskCmd.MarkFlagRequired("task-status")
	confTaskCmd.AddCommand(updateTaskCmd)

	// content convert-ids-to-types
	convertIDsCmd := &cobra.Command{
		Use:   "convert-ids [content-ids...]",
		Short: "Convert content IDs to content types",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"contentIds": args,
			}
			data, err := confPost(cmd, "/content/convert-ids-to-types", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	// Register under the main confluence command for misc operations
	confluenceCmd.AddCommand(convertIDsCmd)

	// users-bulk
	usersBulkCmd := &cobra.Command{
		Use:   "bulk-lookup [account-ids...]",
		Short: "Create bulk user lookup using IDs",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"accountIds": args,
			}
			data, err := confPost(cmd, "/users-bulk", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confUserCmd.AddCommand(usersBulkCmd)

	// user check-access
	checkAccessCmd := &cobra.Command{
		Use:   "check-access [emails...]",
		Short: "Check site access for a list of emails",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"emails": args,
			}
			data, err := confPost(cmd, "/user/access/check-access-by-email", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confUserCmd.AddCommand(checkAccessCmd)

	// user invite
	inviteCmd := &cobra.Command{
		Use:   "invite [emails...]",
		Short: "Invite a list of emails to the site",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"emails": args,
			}
			data, err := confPost(cmd, "/user/access/invite-by-email", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confUserCmd.AddCommand(inviteCmd)

	// admin-key commands
	getAdminKeyCmd := &cobra.Command{
		Use:   "get",
		Short: "Get admin key status",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/admin-key", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confAdminKeyCmd.AddCommand(getAdminKeyCmd)

	enableAdminKeyCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable admin key",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{}
			if d := getIntFlag(cmd, "duration"); d > 0 {
				body["durationInMinutes"] = d
			}
			data, err := confPost(cmd, "/admin-key", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	enableAdminKeyCmd.Flags().Int("duration", 0, "Duration in minutes")
	confAdminKeyCmd.AddCommand(enableAdminKeyCmd)

	disableAdminKeyCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable admin key",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/admin-key", nil)
			if err != nil {
				return err
			}
			fmt.Println("Admin key disabled successfully.")
			return nil
		},
	}
	confAdminKeyCmd.AddCommand(disableAdminKeyCmd)
}
