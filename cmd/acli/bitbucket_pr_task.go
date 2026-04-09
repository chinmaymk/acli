package acli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbPRTaskCmd = &cobra.Command{
	Use:     "task",
	Aliases: []string{"tasks"},
	Short:   "Manage tasks on a pull request",
	RunE:    helpRunE,
}

func init() {
	bbPRCmd.AddCommand(bbPRTaskCmd)

	// task list
	taskListCmd := &cobra.Command{
		Use:     "list [workspace] <repo-slug> <pr-id>",
		Short:   "List tasks on a pull request",
		Aliases: []string{"ls"},
		Args:    cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, idStr, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", idStr)
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			tasks, err := client.ListPRTasks(workspace, repoSlug, prID, getBBPaginationOpts(cmd))
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(tasks)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "ID\tSTATE\tLOCATION\tCONTENT\tCREATOR")
			for _, t := range tasks {
				_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
					t.ID, t.State, taskLocation(&t), truncate(t.Content.Raw, 60), t.Creator.DisplayName)
			}
			return w.Flush()
		},
	}
	addBBPaginationFlags(taskListCmd)
	bbPRTaskCmd.AddCommand(taskListCmd)

	// task get
	bbPRTaskCmd.AddCommand(&cobra.Command{
		Use:   "get [workspace] <repo-slug> <pr-id> <task-id>",
		Short: "Get a task on a pull request",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, prIDStr, taskIDStr, err := resolveWorkspaceRepoIDAndTaskID(cmd, args)
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(prIDStr)
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", prIDStr)
			}
			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				return fmt.Errorf("invalid task ID: %s", taskIDStr)
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			task, err := client.GetPRTask(workspace, repoSlug, prID, taskID)
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(task)
			}

			fmt.Printf("ID:        %d\n", task.ID)
			fmt.Printf("State:     %s\n", task.State)
			fmt.Printf("Content:   %s\n", task.Content.Raw)
			fmt.Printf("Creator:   %s\n", task.Creator.DisplayName)
			fmt.Printf("Created:   %s\n", task.CreatedOn)
			fmt.Printf("Updated:   %s\n", task.UpdatedOn)
			if task.ResolvedOn != "" {
				fmt.Printf("Resolved:  %s\n", task.ResolvedOn)
			}
			if task.Comment != nil {
				fmt.Printf("Comment:   #%d\n", task.Comment.ID)
				if loc := taskLocation(task); loc != "" {
					fmt.Printf("Location:  %s\n", loc)
				}
				if task.Comment.Content.Raw != "" {
					fmt.Printf("Context:   %s\n", truncate(firstLine(task.Comment.Content.Raw), 100))
				}
			}
			return nil
		},
	})

	// task create
	taskCreateCmd := &cobra.Command{
		Use:   "create [workspace] <repo-slug> <pr-id>",
		Short: "Create a task on a pull request",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, idStr, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", idStr)
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			body, _ := cmd.Flags().GetString("body")
			if body == "" {
				return fmt.Errorf("--body is required")
			}

			req := &bitbucket.CreatePRTaskRequest{
				Content: body,
			}

			filePath, _ := cmd.Flags().GetString("file")
			line, _ := cmd.Flags().GetInt("line")

			if filePath != "" && line == 0 {
				return fmt.Errorf("--line is required when --file is specified")
			}
			if line != 0 && filePath == "" {
				return fmt.Errorf("--file is required when --line is specified")
			}

			commentID, _ := cmd.Flags().GetInt("comment-id")

			if filePath != "" && cmd.Flags().Changed("comment-id") {
				return fmt.Errorf("--comment-id cannot be used with --file/--line")
			}

			if filePath != "" {
				// Create an inline comment first, then attach the task to it
				comment, err := client.CreatePRCommentInline(workspace, repoSlug, prID, body, &bitbucket.InlineCommentParams{
					Path: filePath,
					To:   line,
				})
				if err != nil {
					return fmt.Errorf("creating inline comment: %w", err)
				}
				req.CommentID = &comment.ID
			} else if cmd.Flags().Changed("comment-id") {
				req.CommentID = &commentID
			}

			task, err := client.CreatePRTask(workspace, repoSlug, prID, req)
			if err != nil {
				return err
			}
			return outputResult(cmd, "created", fmt.Sprintf("%d", task.ID), fmt.Sprintf("Created task #%d on PR #%d", task.ID, prID), task)
		},
	}
	taskCreateCmd.Flags().String("body", "", "Task content (required)")
	taskCreateCmd.Flags().Int("comment-id", 0, "Associate task with a comment")
	taskCreateCmd.Flags().String("file", "", "File path to attach the task to a specific line (creates an inline comment)")
	taskCreateCmd.Flags().Int("line", 0, "Line number in the new version of the file (requires --file)")
	bbPRTaskCmd.AddCommand(taskCreateCmd)

	// task update
	taskUpdateCmd := &cobra.Command{
		Use:   "update [workspace] <repo-slug> <pr-id> <task-id>",
		Short: "Update a task on a pull request",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, prIDStr, taskIDStr, err := resolveWorkspaceRepoIDAndTaskID(cmd, args)
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(prIDStr)
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", prIDStr)
			}
			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				return fmt.Errorf("invalid task ID: %s", taskIDStr)
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			req := &bitbucket.UpdatePRTaskRequest{}
			if cmd.Flags().Changed("body") {
				body, _ := cmd.Flags().GetString("body")
				req.Content = &body
			}
			if cmd.Flags().Changed("state") {
				req.State, _ = cmd.Flags().GetString("state")
			}

			if req.Content == nil && req.State == "" {
				return fmt.Errorf("at least one of --body or --state is required")
			}

			task, err := client.UpdatePRTask(workspace, repoSlug, prID, taskID, req)
			if err != nil {
				return err
			}
			return outputResult(cmd, "updated", fmt.Sprintf("%d", task.ID), fmt.Sprintf("Updated task #%d on PR #%d", task.ID, prID), task)
		},
	}
	taskUpdateCmd.Flags().String("body", "", "Updated task content")
	taskUpdateCmd.Flags().String("state", "", "Task state (RESOLVED, UNRESOLVED)")
	bbPRTaskCmd.AddCommand(taskUpdateCmd)

	// task resolve (shortcut for update --state RESOLVED)
	bbPRTaskCmd.AddCommand(&cobra.Command{
		Use:   "resolve [workspace] <repo-slug> <pr-id> <task-id>",
		Short: "Resolve a task on a pull request",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, prIDStr, taskIDStr, err := resolveWorkspaceRepoIDAndTaskID(cmd, args)
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(prIDStr)
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", prIDStr)
			}
			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				return fmt.Errorf("invalid task ID: %s", taskIDStr)
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			task, err := client.UpdatePRTask(workspace, repoSlug, prID, taskID, &bitbucket.UpdatePRTaskRequest{
				State: "RESOLVED",
			})
			if err != nil {
				return err
			}
			return outputResult(cmd, "resolved", fmt.Sprintf("%d", task.ID), fmt.Sprintf("Resolved task #%d on PR #%d", task.ID, prID), task)
		},
	})

	// task delete
	bbPRTaskCmd.AddCommand(&cobra.Command{
		Use:   "delete [workspace] <repo-slug> <pr-id> <task-id>",
		Short: "Delete a task on a pull request",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, prIDStr, taskIDStr, err := resolveWorkspaceRepoIDAndTaskID(cmd, args)
			if err != nil {
				return err
			}
			prID, err := strconv.Atoi(prIDStr)
			if err != nil {
				return fmt.Errorf("invalid PR ID: %s", prIDStr)
			}
			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				return fmt.Errorf("invalid task ID: %s", taskIDStr)
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			if err := client.DeletePRTask(workspace, repoSlug, prID, taskID); err != nil {
				return err
			}
			return outputResult(cmd, "deleted", fmt.Sprintf("%d", taskID), fmt.Sprintf("Deleted task #%d from PR #%d", taskID, prID), nil)
		},
	})
}

// resolveWorkspaceRepoIDAndTaskID resolves 4 positional args with optional workspace default.
// With 4 args: workspace=args[0], repo=args[1], id=args[2], taskID=args[3].
// With 3 args: workspace from profile default, repo=args[0], id=args[1], taskID=args[2].
func resolveWorkspaceRepoIDAndTaskID(cmd *cobra.Command, args []string) (string, string, string, string, error) {
	if len(args) >= 4 {
		return args[0], args[1], args[2], args[3], nil
	}
	workspace, err := defaultWorkspace(cmd, nil, 0)
	if err != nil {
		return "", "", "", "", err
	}
	return workspace, args[0], args[1], args[2], nil
}

// taskLocation returns a human-readable location string (e.g. "path/to/file.go:42")
// if the task is attached to an inline comment, or an empty string otherwise.
// Multi-line comments are rendered as "path:start-end".
func taskLocation(t *bitbucket.PRTask) string {
	if t == nil || t.Comment == nil || t.Comment.Inline == nil {
		return ""
	}
	inline := t.Comment.Inline
	if inline.Path == "" {
		return ""
	}
	line := 0
	if inline.To != nil {
		line = *inline.To
	}
	start := 0
	if inline.StartTo != nil {
		start = *inline.StartTo
	}
	switch {
	case start > 0 && line > 0 && start != line:
		return fmt.Sprintf("%s:%d-%d", inline.Path, start, line)
	case line > 0:
		return fmt.Sprintf("%s:%d", inline.Path, line)
	default:
		return inline.Path
	}
}
