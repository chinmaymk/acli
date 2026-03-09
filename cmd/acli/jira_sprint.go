package acli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var jiraSprintCmd = &cobra.Command{
	Use:     "sprint",
	Aliases: []string{"sp"},
	Short:   "Manage sprints (Agile API)",
	RunE:    helpRunE,
}

var jiraSprintGetCmd = &cobra.Command{
	Use:   "get <sprint-id>",
	Short: "Get sprint details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid sprint ID: %w", err)
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")

		sprint, err := client.GetSprint(id)
		if err != nil {
			return err
		}
		if jsonFlag {
			return printJSON(sprint)
		}
		fmt.Printf("ID:         %d\n", sprint.ID)
		fmt.Printf("Name:       %s\n", sprint.Name)
		fmt.Printf("State:      %s\n", sprint.State)
		fmt.Printf("Start Date: %s\n", sprint.StartDate)
		fmt.Printf("End Date:   %s\n", sprint.EndDate)
		if sprint.CompleteDate != "" {
			fmt.Printf("Completed:  %s\n", sprint.CompleteDate)
		}
		fmt.Printf("Board ID:   %d\n", sprint.OriginBoardID)
		if sprint.Goal != "" {
			fmt.Printf("Goal:       %s\n", sprint.Goal)
		}
		return nil
	},
}

var jiraSprintCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a sprint",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		boardID, _ := cmd.Flags().GetInt("board-id")

		body := map[string]interface{}{
			"name":          name,
			"originBoardId": boardID,
		}
		if v, _ := cmd.Flags().GetString("start-date"); v != "" {
			body["startDate"] = v
		}
		if v, _ := cmd.Flags().GetString("end-date"); v != "" {
			body["endDate"] = v
		}
		if v, _ := cmd.Flags().GetString("goal"); v != "" {
			body["goal"] = v
		}

		sprint, err := client.CreateSprint(body)
		if err != nil {
			return err
		}
		fmt.Printf("Sprint created: %s (ID: %d)\n", sprint.Name, sprint.ID)
		return nil
	},
}

var jiraSprintUpdateCmd = &cobra.Command{
	Use:   "update <sprint-id>",
	Short: "Update a sprint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid sprint ID: %w", err)
		}

		body := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			body["name"] = v
		}
		if cmd.Flags().Changed("state") {
			v, _ := cmd.Flags().GetString("state")
			body["state"] = v
		}
		if cmd.Flags().Changed("start-date") {
			v, _ := cmd.Flags().GetString("start-date")
			body["startDate"] = v
		}
		if cmd.Flags().Changed("end-date") {
			v, _ := cmd.Flags().GetString("end-date")
			body["endDate"] = v
		}
		if cmd.Flags().Changed("goal") {
			v, _ := cmd.Flags().GetString("goal")
			body["goal"] = v
		}

		sprint, err := client.UpdateSprint(id, body)
		if err != nil {
			return err
		}
		fmt.Printf("Sprint updated: %s (ID: %d)\n", sprint.Name, sprint.ID)
		return nil
	},
}

var jiraSprintDeleteCmd = &cobra.Command{
	Use:   "delete <sprint-id>",
	Short: "Delete a sprint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid sprint ID: %w", err)
		}
		if err := client.DeleteSprint(id); err != nil {
			return err
		}
		fmt.Printf("Sprint %d deleted.\n", id)
		return nil
	},
}

var jiraSprintIssuesCmd = &cobra.Command{
	Use:   "issues <sprint-id>",
	Short: "List issues in a sprint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid sprint ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		jql, _ := cmd.Flags().GetString("jql")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		result, err := client.GetSprintIssues(id, startAt, maxResults, jql)
		if err != nil {
			return err
		}
		if jsonFlag {
			return printJSON(result)
		}
		w := newTabWriter()
		fmt.Fprintln(w, "KEY\tTYPE\tSTATUS\tPRIORITY\tASSIGNEE\tSUMMARY")
		for _, issue := range result.Issues {
			printIssueRow(w, issue)
		}
		return w.Flush()
	},
}

var jiraSprintMoveCmd = &cobra.Command{
	Use:   "move <sprint-id> <issue-keys...>",
	Short: "Move issues to a sprint",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid sprint ID: %w", err)
		}
		issueKeys := args[1:]
		if err := client.MoveIssuesToSprint(id, issueKeys); err != nil {
			return err
		}
		fmt.Printf("Moved %s to sprint %d\n", strings.Join(issueKeys, ", "), id)
		return nil
	},
}

// --- Epic commands ---

var jiraEpicCmd = &cobra.Command{
	Use:     "epic",
	Aliases: []string{"e"},
	Short:   "Manage epics (Agile API)",
	RunE:    helpRunE,
}

var jiraEpicGetCmd = &cobra.Command{
	Use:   "get <epic-id-or-key>",
	Short: "Get epic details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		epic, err := client.GetEpic(args[0])
		if err != nil {
			return err
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			return printJSON(epic)
		}
		fmt.Printf("ID:      %d\n", epic.ID)
		fmt.Printf("Key:     %s\n", epic.Key)
		fmt.Printf("Name:    %s\n", epic.Name)
		fmt.Printf("Summary: %s\n", epic.Summary)
		fmt.Printf("Done:    %v\n", epic.Done)
		return nil
	},
}

var jiraEpicIssuesCmd = &cobra.Command{
	Use:   "issues <epic-id-or-key>",
	Short: "List issues in an epic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		jql, _ := cmd.Flags().GetString("jql")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		result, err := client.GetEpicIssues(args[0], startAt, maxResults, jql)
		if err != nil {
			return err
		}
		if jsonFlag {
			return printJSON(result)
		}
		w := newTabWriter()
		fmt.Fprintln(w, "KEY\tTYPE\tSTATUS\tPRIORITY\tASSIGNEE\tSUMMARY")
		for _, issue := range result.Issues {
			printIssueRow(w, issue)
		}
		return w.Flush()
	},
}

var jiraEpicMoveCmd = &cobra.Command{
	Use:   "move <epic-id-or-key> <issue-keys...>",
	Short: "Move issues to an epic",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.MoveIssuesToEpic(args[0], args[1:]); err != nil {
			return err
		}
		fmt.Printf("Moved %s to epic %s\n", strings.Join(args[1:], ", "), args[0])
		return nil
	},
}

// --- Backlog command ---

var jiraBacklogMoveCmd = &cobra.Command{
	Use:   "backlog <issue-keys...>",
	Short: "Move issues to the backlog",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.MoveIssuesToBacklog(args); err != nil {
			return err
		}
		fmt.Printf("Moved %s to backlog\n", strings.Join(args, ", "))
		return nil
	},
}

func init() {
	// sprint get
	jiraSprintGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraSprintCmd.AddCommand(jiraSprintGetCmd)

	// sprint create
	jiraSprintCreateCmd.Flags().String("name", "", "Sprint name (required)")
	_ = jiraSprintCreateCmd.MarkFlagRequired("name")
	jiraSprintCreateCmd.Flags().Int("board-id", 0, "Origin board ID (required)")
	_ = jiraSprintCreateCmd.MarkFlagRequired("board-id")
	jiraSprintCreateCmd.Flags().String("start-date", "", "Start date (ISO 8601)")
	jiraSprintCreateCmd.Flags().String("end-date", "", "End date (ISO 8601)")
	jiraSprintCreateCmd.Flags().String("goal", "", "Sprint goal")
	jiraSprintCmd.AddCommand(jiraSprintCreateCmd)

	// sprint update
	jiraSprintUpdateCmd.Flags().String("name", "", "Sprint name")
	jiraSprintUpdateCmd.Flags().String("state", "", "Sprint state (active, closed, future)")
	jiraSprintUpdateCmd.Flags().String("start-date", "", "Start date (ISO 8601)")
	jiraSprintUpdateCmd.Flags().String("end-date", "", "End date (ISO 8601)")
	jiraSprintUpdateCmd.Flags().String("goal", "", "Sprint goal")
	jiraSprintCmd.AddCommand(jiraSprintUpdateCmd)

	// sprint delete
	jiraSprintCmd.AddCommand(jiraSprintDeleteCmd)

	// sprint issues
	jiraSprintIssuesCmd.Flags().Int("start-at", 0, "Start index")
	jiraSprintIssuesCmd.Flags().Int("max-results", 50, "Max results")
	jiraSprintIssuesCmd.Flags().String("jql", "", "JQL filter")
	jiraSprintIssuesCmd.Flags().Bool("json", false, "Output as JSON")
	jiraSprintCmd.AddCommand(jiraSprintIssuesCmd)

	// sprint move
	jiraSprintCmd.AddCommand(jiraSprintMoveCmd)

	// epic get
	jiraEpicGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraEpicCmd.AddCommand(jiraEpicGetCmd)

	// epic issues
	jiraEpicIssuesCmd.Flags().Int("start-at", 0, "Start index")
	jiraEpicIssuesCmd.Flags().Int("max-results", 50, "Max results")
	jiraEpicIssuesCmd.Flags().String("jql", "", "JQL filter")
	jiraEpicIssuesCmd.Flags().Bool("json", false, "Output as JSON")
	jiraEpicCmd.AddCommand(jiraEpicIssuesCmd)

	// epic move
	jiraEpicCmd.AddCommand(jiraEpicMoveCmd)
}
