package acli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var jiraBoardCmd = &cobra.Command{
	Use:     "board",
	Aliases: []string{"b"},
	Short:   "Manage boards (Agile API)",
	RunE:    helpRunE,
}

var jiraBoardListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List boards",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		project, _ := defaultProject(cmd)
		boardType, _ := cmd.Flags().GetString("type")
		name, _ := cmd.Flags().GetString("name")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		result, err := client.GetBoards(startAt, maxResults, project, boardType, name)
		if err != nil {
			return err
		}
		if jsonFlag {
			return printJSON(result)
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tPROJECT")
		for _, b := range result.Values {
			project := ""
			if b.Location != nil {
				project = b.Location.ProjectKey
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", b.ID, b.Name, b.Type, project)
		}
		return w.Flush()
	},
}

var jiraBoardGetCmd = &cobra.Command{
	Use:   "get <board-id>",
	Short: "Get board details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		jsonFlag, _ := cmd.Flags().GetBool("json")

		board, err := client.GetBoard(id)
		if err != nil {
			return err
		}
		if jsonFlag {
			return printJSON(board)
		}
		fmt.Printf("ID:      %d\n", board.ID)
		fmt.Printf("Name:    %s\n", board.Name)
		fmt.Printf("Type:    %s\n", board.Type)
		if board.Location != nil {
			fmt.Printf("Project: %s (%s)\n", board.Location.ProjectName, board.Location.ProjectKey)
		}
		return nil
	},
}

var jiraBoardConfigCmd = &cobra.Command{
	Use:     "config <board-id>",
	Aliases: []string{"configuration"},
	Short:   "Get board configuration",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		config, err := client.GetBoardConfiguration(id)
		if err != nil {
			return err
		}
		return printJSON(config)
	},
}

var jiraBoardIssuesCmd = &cobra.Command{
	Use:   "issues <board-id>",
	Short: "List issues on a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		jql, _ := cmd.Flags().GetString("jql")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		result, err := client.GetBoardIssues(id, startAt, maxResults, jql)
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

var jiraBoardBacklogCmd = &cobra.Command{
	Use:   "backlog <board-id>",
	Short: "List backlog issues for a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		jql, _ := cmd.Flags().GetString("jql")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		result, err := client.GetBoardBacklog(id, startAt, maxResults, jql)
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

var jiraBoardSprintsCmd = &cobra.Command{
	Use:   "sprints <board-id>",
	Short: "List sprints for a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		state, _ := cmd.Flags().GetString("state")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		result, err := client.GetBoardSprints(id, startAt, maxResults, state)
		if err != nil {
			return err
		}
		if jsonFlag {
			return printJSON(result)
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tSTATE\tSTART DATE\tEND DATE")
		for _, s := range result.Values {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", s.ID, s.Name, s.State, s.StartDate, s.EndDate)
		}
		return w.Flush()
	},
}

var jiraBoardEpicsCmd = &cobra.Command{
	Use:   "epics <board-id>",
	Short: "List epics for a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		result, err := client.GetBoardEpics(id, startAt, maxResults)
		if err != nil {
			return err
		}
		if jsonFlag {
			return printJSON(result)
		}
		w := newTabWriter()
		fmt.Fprintln(w, "ID\tKEY\tNAME\tDONE\tSUMMARY")
		for _, e := range result.Values {
			fmt.Fprintf(w, "%d\t%s\t%s\t%v\t%s\n", e.ID, e.Key, e.Name, e.Done, e.Summary)
		}
		return w.Flush()
	},
}

func init() {
	// board list
	jiraBoardListCmd.Flags().String("project", "", "Filter by project key or ID (uses profile default if not set)")
	jiraBoardListCmd.Flags().String("type", "", "Filter by board type (scrum, kanban)")
	jiraBoardListCmd.Flags().String("name", "", "Filter by board name")
	jiraBoardListCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardListCmd.Flags().Int("max-results", 50, "Max results")
	jiraBoardListCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardListCmd)

	// board get
	jiraBoardGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardGetCmd)

	// board config
	jiraBoardCmd.AddCommand(jiraBoardConfigCmd)

	// board issues
	jiraBoardIssuesCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardIssuesCmd.Flags().Int("max-results", 50, "Max results")
	jiraBoardIssuesCmd.Flags().String("jql", "", "JQL filter")
	jiraBoardIssuesCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardIssuesCmd)

	// board backlog
	jiraBoardBacklogCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardBacklogCmd.Flags().Int("max-results", 50, "Max results")
	jiraBoardBacklogCmd.Flags().String("jql", "", "JQL filter")
	jiraBoardBacklogCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardBacklogCmd)

	// board sprints
	jiraBoardSprintsCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardSprintsCmd.Flags().Int("max-results", 50, "Max results")
	jiraBoardSprintsCmd.Flags().String("state", "", "Filter by state (active, closed, future)")
	jiraBoardSprintsCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardSprintsCmd)

	// board epics
	jiraBoardEpicsCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardEpicsCmd.Flags().Int("max-results", 50, "Max results")
	jiraBoardEpicsCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardEpicsCmd)
}
