package acli

import (
	"fmt"
	"strings"

	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

var jiraIssueCmd = &cobra.Command{
	Use:     "issue",
	Aliases: []string{"i"},
	Short:   "Manage Jira issues",
	RunE:    helpRunE,
}

// --- issue list ---

var jiraIssueListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List issues using JQL search",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		jql, _ := cmd.Flags().GetString("jql")
		project, _ := defaultProject(cmd)
		assignee, _ := cmd.Flags().GetString("assignee")
		status, _ := cmd.Flags().GetString("status")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")
		all, _ := cmd.Flags().GetBool("all")
		jsonOutput := isJSONOutput(cmd)

		if jql == "" {
			var clauses []string
			if project != "" {
				clauses = append(clauses, fmt.Sprintf("project = %s", project))
			}
			if assignee != "" {
				clauses = append(clauses, fmt.Sprintf("assignee = %q", assignee))
			}
			if status != "" {
				clauses = append(clauses, fmt.Sprintf("status = %q", status))
			}
			if len(clauses) == 0 {
				jql = "created >= -30d order by created DESC"
			} else {
				jql = strings.Join(clauses, " AND ") + " order by created DESC"
			}
		}

		fields := []string{"summary", "issuetype", "status", "priority", "assignee"}
		results, err := client.SearchJQL(jql, startAt, maxResults, fields, nil)
		if err != nil {
			return err
		}

		allIssues := results.Issues
		if all {
			for len(allIssues) < results.Total {
				next, err := client.SearchJQL(jql, startAt+len(allIssues), maxResults, fields, nil)
				if err != nil {
					return err
				}
				if len(next.Issues) == 0 {
					break
				}
				allIssues = append(allIssues, next.Issues...)
			}
			results.Issues = allIssues
		}

		if jsonOutput {
			return outputJSON(results)
		}

		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "KEY\tTYPE\tSTATUS\tPRIORITY\tASSIGNEE\tSUMMARY\n")
		for _, issue := range results.Issues {
			issueType := ""
			if issue.Fields.IssueType != nil {
				issueType = issue.Fields.IssueType.Name
			}
			status := ""
			if issue.Fields.Status != nil {
				status = issue.Fields.Status.Name
			}
			priority := ""
			if issue.Fields.Priority != nil {
				priority = issue.Fields.Priority.Name
			}
			assignee := ""
			if issue.Fields.Assignee != nil {
				assignee = issue.Fields.Assignee.DisplayName
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				issue.Key, issueType, status, priority, assignee, issue.Fields.Summary)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(results.Issues), results.Total)
		return nil
	},
}

// --- issue get ---

var jiraIssueGetCmd = &cobra.Command{
	Use:   "get <issue-key>",
	Short: "Get issue details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		issue, err := client.GetIssue(args[0], nil, nil)
		if err != nil {
			return err
		}

		if isJSONOutput(cmd) {
			return outputJSON(issue)
		}

		f := issue.Fields

		issueType := ""
		if f.IssueType != nil {
			issueType = f.IssueType.Name
		}
		status := ""
		if f.Status != nil {
			status = f.Status.Name
		}
		priority := ""
		if f.Priority != nil {
			priority = f.Priority.Name
		}
		assignee := ""
		if f.Assignee != nil {
			assignee = f.Assignee.DisplayName
		}
		reporter := ""
		if f.Reporter != nil {
			reporter = f.Reporter.DisplayName
		}

		description := ""
		if f.Description != nil {
			switch f.Description.(type) {
			case map[string]interface{}:
				description = "[Atlassian Document Format]"
			case string:
				description = f.Description.(string)
			default:
				description = fmt.Sprintf("%v", f.Description)
			}
		}

		var labels []string
		labels = append(labels, f.Labels...)

		var components []string
		for _, c := range f.Components {
			components = append(components, c.Name)
		}

		var fixVersions []string
		for _, v := range f.FixVersions {
			fixVersions = append(fixVersions, v.Name)
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Key:          %s\n", issue.Key)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Summary:      %s\n", f.Summary)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Status:       %s\n", status)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Type:         %s\n", issueType)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Priority:     %s\n", priority)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Assignee:     %s\n", assignee)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Reporter:     %s\n", reporter)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Created:      %s\n", f.Created)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Updated:      %s\n", f.Updated)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Labels:       %s\n", strings.Join(labels, ", "))
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Components:   %s\n", strings.Join(components, ", "))
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Fix Versions: %s\n", strings.Join(fixVersions, ", "))
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Description:\n%s\n", description)

		return nil
	},
}

// --- issue create ---

var jiraIssueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		project, err := defaultProject(cmd)
		if err != nil {
			return err
		}
		if project == "" {
			return fmt.Errorf("--project is required (or set a default with 'acli config set-defaults')")
		}
		issueType, _ := cmd.Flags().GetString("type")
		summary, _ := cmd.Flags().GetString("summary")
		description, _ := cmd.Flags().GetString("description")
		assignee, _ := cmd.Flags().GetString("assignee")
		priority, _ := cmd.Flags().GetString("priority")
		labels, _ := cmd.Flags().GetStringSlice("labels")
		components, _ := cmd.Flags().GetStringSlice("components")

		fields := map[string]interface{}{
			"project":   map[string]interface{}{"key": project},
			"issuetype": map[string]interface{}{"name": issueType},
			"summary":   summary,
		}

		if description != "" {
			fields["description"] = map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": description,
							},
						},
					},
				},
			}
		}

		if assignee != "" {
			fields["assignee"] = map[string]interface{}{"accountId": assignee}
		}

		if priority != "" {
			fields["priority"] = map[string]interface{}{"name": priority}
		}

		if len(labels) > 0 {
			fields["labels"] = labels
		}

		if len(components) > 0 {
			comps := make([]map[string]interface{}, len(components))
			for i, c := range components {
				comps[i] = map[string]interface{}{"name": c}
			}
			fields["components"] = comps
		}

		details := &jira.IssueUpdateDetails{
			Fields: fields,
		}

		created, err := client.CreateIssue(details)
		if err != nil {
			return err
		}

		return outputResult(cmd, "created", created.Key, fmt.Sprintf("Created issue: %s", created.Key), created)
	},
}

// --- issue edit ---

var jiraIssueEditCmd = &cobra.Command{
	Use:   "edit <issue-key>",
	Short: "Edit an existing issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		fields := map[string]interface{}{}

		if cmd.Flags().Changed("summary") {
			v, _ := cmd.Flags().GetString("summary")
			fields["summary"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			fields["description"] = map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": v,
							},
						},
					},
				},
			}
		}
		if cmd.Flags().Changed("assignee") {
			v, _ := cmd.Flags().GetString("assignee")
			fields["assignee"] = map[string]interface{}{"accountId": v}
		}
		if cmd.Flags().Changed("priority") {
			v, _ := cmd.Flags().GetString("priority")
			fields["priority"] = map[string]interface{}{"name": v}
		}
		if cmd.Flags().Changed("labels") {
			v, _ := cmd.Flags().GetStringSlice("labels")
			fields["labels"] = v
		}
		if cmd.Flags().Changed("components") {
			v, _ := cmd.Flags().GetStringSlice("components")
			comps := make([]map[string]interface{}, len(v))
			for i, c := range v {
				comps[i] = map[string]interface{}{"name": c}
			}
			fields["components"] = comps
		}

		if len(fields) == 0 {
			return fmt.Errorf("no fields specified to update")
		}

		details := &jira.IssueUpdateDetails{
			Fields: fields,
		}

		err = client.EditIssue(args[0], details, true)
		if err != nil {
			return err
		}

		return outputResult(cmd, "updated", args[0], fmt.Sprintf("Issue %s updated successfully", args[0]), nil)
	},
}

// --- issue delete ---

var jiraIssueDeleteCmd = &cobra.Command{
	Use:   "delete <issue-key>",
	Short: "Delete an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		deleteSubtasks, _ := cmd.Flags().GetBool("delete-subtasks")

		err = client.DeleteIssue(args[0], deleteSubtasks)
		if err != nil {
			return err
		}

		return outputResult(cmd, "deleted", args[0], fmt.Sprintf("Issue %s deleted", args[0]), nil)
	},
}

// --- issue assign ---

var jiraIssueAssignCmd = &cobra.Command{
	Use:   "assign <issue-key> <account-id>",
	Short: "Assign an issue to a user (use '-1' or 'none' to unassign)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		accountID := args[1]
		if accountID == "-1" || accountID == "none" {
			accountID = ""
		}

		err = client.AssignIssue(args[0], accountID)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Issue %s unassigned", args[0])
		if accountID != "" {
			msg = fmt.Sprintf("Issue %s assigned to %s", args[0], args[1])
		}
		return outputResult(cmd, "assigned", args[0], msg, nil)
	},
}

// --- issue transition ---

var jiraIssueTransitionCmd = &cobra.Command{
	Use:   "transition <issue-key>",
	Short: "Transition an issue to a new status",
	Long:  "Transition an issue to a new status. Provide either --id (transition ID) or --status (target status name).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		transitionID, _ := cmd.Flags().GetString("id")
		statusName, _ := cmd.Flags().GetString("status")

		if transitionID == "" && statusName == "" {
			return fmt.Errorf("either --id or --status must be provided")
		}

		if transitionID == "" {
			// Look up the transition ID by matching the status name
			resp, err := client.GetIssueTransitions(args[0])
			if err != nil {
				return fmt.Errorf("failed to get transitions: %w", err)
			}

			for _, t := range resp.Transitions {
				if strings.EqualFold(t.Name, statusName) || (t.To != nil && strings.EqualFold(t.To.Name, statusName)) {
					transitionID = t.ID
					break
				}
			}

			if transitionID == "" {
				// Build a helpful error listing available transitions
				var available []string
				for _, t := range resp.Transitions {
					name := t.Name
					if t.To != nil && t.To.Name != t.Name {
						name = fmt.Sprintf("%s (-> %s)", t.Name, t.To.Name)
					}
					available = append(available, name)
				}
				return fmt.Errorf("no transition found matching status %q; available transitions: %s", statusName, strings.Join(available, ", "))
			}
		}

		details := &jira.IssueUpdateDetails{
			Transition: &jira.IssueTransition{
				ID: transitionID,
			},
		}

		err = client.DoIssueTransition(args[0], details)
		if err != nil {
			return err
		}

		return outputResult(cmd, "transitioned", args[0], fmt.Sprintf("Issue %s transitioned successfully", args[0]), nil)
	},
}

// --- issue transitions ---

var jiraIssueTransitionsCmd = &cobra.Command{
	Use:   "transitions <issue-key>",
	Short: "List available transitions for an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		resp, err := client.GetIssueTransitions(args[0])
		if err != nil {
			return err
		}

		if isJSONOutput(cmd) {
			return outputJSON(resp)
		}

		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "ID\tNAME\tTO STATUS\n")
		for _, t := range resp.Transitions {
			toStatus := ""
			if t.To != nil {
				toStatus = t.To.Name
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", t.ID, t.Name, toStatus)
		}
		_ = w.Flush()
		return nil
	},
}

// --- issue comment (group) ---

var jiraIssueCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage issue comments",
	RunE:  helpRunE,
}

// --- issue comment list ---

var jiraIssueCommentListCmd = &cobra.Command{
	Use:     "list <issue-key>",
	Aliases: []string{"ls"},
	Short:   "List comments on an issue",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")
		all, _ := cmd.Flags().GetBool("all")

		page, err := client.GetIssueComments(args[0], startAt, maxResults)
		if err != nil {
			return err
		}

		allComments := page.Comments
		if all {
			for len(allComments) < page.Total {
				next, err := client.GetIssueComments(args[0], startAt+len(allComments), maxResults)
				if err != nil {
					return err
				}
				if len(next.Comments) == 0 {
					break
				}
				allComments = append(allComments, next.Comments...)
			}
			page.Comments = allComments
		}

		if isJSONOutput(cmd) {
			return outputJSON(page)
		}

		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "ID\tAUTHOR\tCREATED\tBODY\n")
		for _, c := range page.Comments {
			author := ""
			if c.Author != nil {
				author = c.Author.DisplayName
			}
			body := extractCommentBody(c.Body)
			if len(body) > 60 {
				body = body[:57] + "..."
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.ID, author, c.Created, body)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(page.Comments), page.Total)
		return nil
	},
}

// extractCommentBody attempts to extract text from a comment body.
// If it's ADF (map), it tries to extract text nodes; otherwise returns the string representation.
func extractCommentBody(body interface{}) string {
	if body == nil {
		return ""
	}
	switch b := body.(type) {
	case string:
		return b
	case map[string]interface{}:
		// Try to extract text from ADF
		return extractADFText(b)
	default:
		return "ADF content"
	}
}

// extractADFText recursively extracts text from an ADF document.
func extractADFText(node map[string]interface{}) string {
	if t, ok := node["type"].(string); ok && t == "text" {
		if text, ok := node["text"].(string); ok {
			return text
		}
	}
	content, ok := node["content"].([]interface{})
	if !ok {
		return "ADF content"
	}
	var parts []string
	for _, item := range content {
		if m, ok := item.(map[string]interface{}); ok {
			text := extractADFText(m)
			if text != "" && text != "ADF content" {
				parts = append(parts, text)
			}
		}
	}
	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}
	return "ADF content"
}

// --- issue comment add ---

var jiraIssueCommentAddCmd = &cobra.Command{
	Use:   "add <issue-key>",
	Short: "Add a comment to an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		bodyText, _ := cmd.Flags().GetString("body")

		adfBody := map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []interface{}{
				map[string]interface{}{
					"type": "paragraph",
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": bodyText,
						},
					},
				},
			},
		}

		comment, err := client.AddIssueComment(args[0], adfBody, nil)
		if err != nil {
			return err
		}

		return outputResult(cmd, "created", comment.ID, fmt.Sprintf("Comment %s added to %s", comment.ID, args[0]), comment)
	},
}

// --- issue comment get ---

var jiraIssueCommentGetCmd = &cobra.Command{
	Use:   "get <issue-key> <comment-id>",
	Short: "Get a specific comment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		comment, err := client.GetIssueComment(args[0], args[1])
		if err != nil {
			return err
		}

		return outputJSON(comment)
	},
}

// --- issue comment delete ---

var jiraIssueCommentDeleteCmd = &cobra.Command{
	Use:   "delete <issue-key> <comment-id>",
	Short: "Delete a comment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		err = client.DeleteIssueComment(args[0], args[1])
		if err != nil {
			return err
		}

		return outputResult(cmd, "deleted", args[1], fmt.Sprintf("Comment %s deleted from %s", args[1], args[0]), nil)
	},
}

// --- issue worklog (group) ---

var jiraIssueWorklogCmd = &cobra.Command{
	Use:   "worklog",
	Short: "Manage issue worklogs",
	RunE:  helpRunE,
}

// --- issue worklog list ---

var jiraIssueWorklogListCmd = &cobra.Command{
	Use:     "list <issue-key>",
	Aliases: []string{"ls"},
	Short:   "List worklogs on an issue",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")
		all, _ := cmd.Flags().GetBool("all")

		page, err := client.GetIssueWorklogs(args[0], startAt, maxResults)
		if err != nil {
			return err
		}

		allWorklogs := page.Worklogs
		if all {
			for len(allWorklogs) < page.Total {
				next, err := client.GetIssueWorklogs(args[0], startAt+len(allWorklogs), maxResults)
				if err != nil {
					return err
				}
				if len(next.Worklogs) == 0 {
					break
				}
				allWorklogs = append(allWorklogs, next.Worklogs...)
			}
			page.Worklogs = allWorklogs
		}

		if isJSONOutput(cmd) {
			return outputJSON(page)
		}

		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "ID\tAUTHOR\tTIME SPENT\tSTARTED\n")
		for _, wl := range page.Worklogs {
			author := ""
			if wl.Author != nil {
				author = wl.Author.DisplayName
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", wl.ID, author, wl.TimeSpent, wl.Started)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(page.Worklogs), page.Total)
		return nil
	},
}

// --- issue worklog add ---

var jiraIssueWorklogAddCmd = &cobra.Command{
	Use:   "add <issue-key>",
	Short: "Add a worklog entry to an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		timeSpent, _ := cmd.Flags().GetString("time-spent")
		started, _ := cmd.Flags().GetString("started")

		worklog := &jira.Worklog{
			TimeSpent: timeSpent,
		}
		if started != "" {
			worklog.Started = started
		}

		result, err := client.AddIssueWorklog(args[0], worklog)
		if err != nil {
			return err
		}

		return outputResult(cmd, "created", result.ID, fmt.Sprintf("Worklog %s added to %s", result.ID, args[0]), result)
	},
}

// --- issue worklog delete ---

var jiraIssueWorklogDeleteCmd = &cobra.Command{
	Use:   "delete <issue-key> <worklog-id>",
	Short: "Delete a worklog entry",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		err = client.DeleteIssueWorklog(args[0], args[1])
		if err != nil {
			return err
		}

		return outputResult(cmd, "deleted", args[1], fmt.Sprintf("Worklog %s deleted from %s", args[1], args[0]), nil)
	},
}

// --- issue attach ---

var jiraIssueAttachCmd = &cobra.Command{
	Use:   "attach <issue-key> <file-path>",
	Short: "Upload an attachment to an issue",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		attachments, err := client.AddIssueAttachment(args[0], args[1])
		if err != nil {
			return err
		}

		if isJSONOutput(cmd) {
			return outputJSON(attachments)
		}
		for _, a := range attachments {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Attached: %s (id: %s)\n", a.Filename, a.ID)
		}
		return nil
	},
}

// --- issue vote ---

var jiraIssueVoteCmd = &cobra.Command{
	Use:   "vote <issue-key>",
	Short: "Add your vote to an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		err = client.AddIssueVote(args[0])
		if err != nil {
			return err
		}

		return outputResult(cmd, "voted", args[0], fmt.Sprintf("Vote added to %s", args[0]), nil)
	},
}

// --- issue unvote ---

var jiraIssueUnvoteCmd = &cobra.Command{
	Use:   "unvote <issue-key>",
	Short: "Remove your vote from an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		err = client.RemoveIssueVote(args[0])
		if err != nil {
			return err
		}

		return outputResult(cmd, "unvoted", args[0], fmt.Sprintf("Vote removed from %s", args[0]), nil)
	},
}

// --- issue watch ---

var jiraIssueWatchCmd = &cobra.Command{
	Use:   "watch <issue-key>",
	Short: "Add a watcher to an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		accountID, _ := cmd.Flags().GetString("account-id")

		err = client.AddIssueWatcher(args[0], accountID)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("You are now watching %s", args[0])
		if accountID != "" {
			msg = fmt.Sprintf("Watcher %s added to %s", accountID, args[0])
		}
		return outputResult(cmd, "watch_added", args[0], msg, nil)
	},
}

// --- issue unwatch ---

var jiraIssueUnwatchCmd = &cobra.Command{
	Use:   "unwatch <issue-key>",
	Short: "Remove a watcher from an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		accountID, _ := cmd.Flags().GetString("account-id")

		err = client.RemoveIssueWatcher(args[0], accountID)
		if err != nil {
			return err
		}

		return outputResult(cmd, "watch_removed", args[0], fmt.Sprintf("Watcher %s removed from %s", accountID, args[0]), nil)
	},
}

// --- issue watchers ---

var jiraIssueWatchersCmd = &cobra.Command{
	Use:   "watchers <issue-key>",
	Short: "List watchers of an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		watches, err := client.GetIssueWatchers(args[0])
		if err != nil {
			return err
		}

		if isJSONOutput(cmd) {
			return outputJSON(watches)
		}

		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "ACCOUNT ID\tDISPLAY NAME\n")
		for _, watcher := range watches.Watchers {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", watcher.AccountID, watcher.DisplayName)
		}
		_ = w.Flush()
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal watchers: %d\n", watches.WatchCount)
		return nil
	},
}

// --- issue changelog ---

var jiraIssueChangelogCmd = &cobra.Command{
	Use:   "changelog <issue-key>",
	Short: "List changelog for an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		maxResults, _ := cmd.Flags().GetInt("max-results")
		startAt, _ := cmd.Flags().GetInt("start-at")
		all, _ := cmd.Flags().GetBool("all")

		page, err := client.GetIssueChangelog(args[0], startAt, maxResults)
		if err != nil {
			return err
		}

		if all {
			allValues := page.Values
			for len(allValues) < page.Total {
				next, err := client.GetIssueChangelog(args[0], startAt+len(allValues), maxResults)
				if err != nil {
					return err
				}
				if len(next.Values) == 0 {
					break
				}
				allValues = append(allValues, next.Values...)
			}
			page.Values = allValues
		}

		if isJSONOutput(cmd) {
			return outputJSON(page)
		}

		// The changelog API may return values in "Values" or "Histories"
		histories := page.Values
		if len(histories) == 0 {
			histories = page.Histories
		}

		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "DATE\tAUTHOR\tFIELD\tFROM\tTO\n")
		for _, entry := range histories {
			author := ""
			if entry.Author != nil {
				author = entry.Author.DisplayName
			}
			for _, item := range entry.Items {
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					entry.Created, author, item.Field, item.FromString, item.ToString)
			}
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(histories), page.Total)
		return nil
	},
}

// --- issue link (group for remote links) ---

var jiraIssueRemoteLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Manage remote links on issues",
	RunE:  helpRunE,
}

// --- issue link list ---

var jiraIssueRemoteLinkListCmd = &cobra.Command{
	Use:     "list <issue-key>",
	Aliases: []string{"ls"},
	Short:   "List remote links on an issue",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		links, err := client.GetIssueRemoteLinks(args[0])
		if err != nil {
			return err
		}

		if isJSONOutput(cmd) {
			return outputJSON(links)
		}

		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "ID\tTITLE\tURL\n")
		for _, link := range links {
			title := ""
			url := ""
			if link.Object != nil {
				title = link.Object.Title
				url = link.Object.URL
			}
			_, _ = fmt.Fprintf(w, "%d\t%s\t%s\n", link.ID, title, url)
		}
		_ = w.Flush()
		return nil
	},
}

// --- issue link create ---

var jiraIssueRemoteLinkCreateCmd = &cobra.Command{
	Use:   "create <issue-key>",
	Short: "Create a remote link on an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		linkURL, _ := cmd.Flags().GetString("url")
		title, _ := cmd.Flags().GetString("title")

		link := &jira.RemoteIssueLink{
			Object: &jira.RemoteObject{
				URL:   linkURL,
				Title: title,
			},
		}

		created, err := client.CreateIssueRemoteLink(args[0], link)
		if err != nil {
			return err
		}

		return outputResult(cmd, "created", fmt.Sprintf("%d", created.ID), fmt.Sprintf("Remote link %d created on %s", created.ID, args[0]), created)
	},
}

// --- issue link delete ---

var jiraIssueRemoteLinkDeleteCmd = &cobra.Command{
	Use:   "delete <issue-key> <link-id>",
	Short: "Delete a remote link from an issue",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		err = client.DeleteIssueRemoteLink(args[0], args[1])
		if err != nil {
			return err
		}

		return outputResult(cmd, "deleted", args[1], fmt.Sprintf("Remote link %s deleted from %s", args[1], args[0]), nil)
	},
}

// --- issue notify ---

var jiraIssueNotifyCmd = &cobra.Command{
	Use:   "notify <issue-key>",
	Short: "Send a notification for an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		subject, _ := cmd.Flags().GetString("subject")
		textBody, _ := cmd.Flags().GetString("text-body")

		notify := &jira.IssueNotifyRequest{
			Subject:  subject,
			TextBody: textBody,
		}

		err = client.NotifyIssue(args[0], notify)
		if err != nil {
			return err
		}

		return outputResult(cmd, "notified", args[0], fmt.Sprintf("Notification sent for %s", args[0]), nil)
	},
}

// --- issue createmeta ---

var jiraIssueCreateMetaCmd = &cobra.Command{
	Use:   "createmeta",
	Short: "Get issue create metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		projects, _ := cmd.Flags().GetStringSlice("project")

		meta, err := client.GetCreateMeta(projects, []string{"projects.issuetypes.fields"})
		if err != nil {
			return err
		}

		return outputJSON(meta)
	},
}

// --- issue editmeta ---

var jiraIssueEditMetaCmd = &cobra.Command{
	Use:   "editmeta <issue-key>",
	Short: "Get issue edit metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}

		meta, err := client.GetIssueEditMeta(args[0])
		if err != nil {
			return err
		}

		return outputJSON(meta)
	},
}

// --- init: register all subcommands ---

func init() {
	// issue list flags
	jiraIssueListCmd.Flags().String("jql", "", "JQL query string (overrides convenience flags)")
	jiraIssueListCmd.Flags().String("project", "", "Filter by project key (uses profile default if not set)")
	jiraIssueListCmd.Flags().String("assignee", "", "Filter by assignee")
	jiraIssueListCmd.Flags().String("status", "", "Filter by status")
	jiraIssueListCmd.Flags().Int("max-results", 50, "Maximum number of results per page")
	jiraIssueListCmd.Flags().Int("start-at", 0, "Index of the first result")
	addAllFlag(jiraIssueListCmd)
	jiraIssueListCmd.Flags().Bool("json", false, "Output as JSON")

	// issue get flags
	jiraIssueGetCmd.Flags().Bool("json", false, "Output as JSON")

	// issue create flags
	jiraIssueCreateCmd.Flags().String("project", "", "Project key (uses profile default if not set)")
	jiraIssueCreateCmd.Flags().String("type", "Task", "Issue type")
	jiraIssueCreateCmd.Flags().String("summary", "", "Issue summary (required)")
	jiraIssueCreateCmd.Flags().String("description", "", "Issue description")
	jiraIssueCreateCmd.Flags().String("assignee", "", "Assignee account ID")
	jiraIssueCreateCmd.Flags().String("priority", "", "Priority name")
	jiraIssueCreateCmd.Flags().StringSlice("labels", nil, "Labels")
	jiraIssueCreateCmd.Flags().StringSlice("components", nil, "Component names")
	_ = jiraIssueCreateCmd.MarkFlagRequired("summary")

	// issue edit flags
	jiraIssueEditCmd.Flags().String("summary", "", "Issue summary")
	jiraIssueEditCmd.Flags().String("description", "", "Issue description")
	jiraIssueEditCmd.Flags().String("assignee", "", "Assignee account ID")
	jiraIssueEditCmd.Flags().String("priority", "", "Priority name")
	jiraIssueEditCmd.Flags().StringSlice("labels", nil, "Labels")
	jiraIssueEditCmd.Flags().StringSlice("components", nil, "Component names")

	// issue delete flags
	jiraIssueDeleteCmd.Flags().Bool("delete-subtasks", false, "Also delete subtasks")

	// issue transition flags
	jiraIssueTransitionCmd.Flags().String("id", "", "Transition ID")
	jiraIssueTransitionCmd.Flags().String("status", "", "Target status name (case-insensitive)")

	// comment list flags
	jiraIssueCommentListCmd.Flags().Int("max-results", 50, "Maximum number of results per page")
	jiraIssueCommentListCmd.Flags().Int("start-at", 0, "Index of the first result")
	addAllFlag(jiraIssueCommentListCmd)
	jiraIssueCommentListCmd.Flags().Bool("json", false, "Output as JSON")

	// comment add flags
	jiraIssueCommentAddCmd.Flags().String("body", "", "Comment body text (required)")
	_ = jiraIssueCommentAddCmd.MarkFlagRequired("body")

	// worklog list flags
	jiraIssueWorklogListCmd.Flags().Int("max-results", 50, "Maximum number of results per page")
	jiraIssueWorklogListCmd.Flags().Int("start-at", 0, "Index of the first result")
	addAllFlag(jiraIssueWorklogListCmd)

	// worklog add flags
	jiraIssueWorklogAddCmd.Flags().String("time-spent", "", "Time spent (e.g. '2h', '30m') (required)")
	jiraIssueWorklogAddCmd.Flags().String("started", "", "Start time (ISO datetime)")
	_ = jiraIssueWorklogAddCmd.MarkFlagRequired("time-spent")

	// watch flags
	jiraIssueWatchCmd.Flags().String("account-id", "", "Account ID of user to add as watcher (default: self)")

	// unwatch flags
	jiraIssueUnwatchCmd.Flags().String("account-id", "", "Account ID of user to remove as watcher (required)")
	_ = jiraIssueUnwatchCmd.MarkFlagRequired("account-id")

	// changelog flags
	jiraIssueChangelogCmd.Flags().Int("max-results", 50, "Maximum number of results per page")
	jiraIssueChangelogCmd.Flags().Int("start-at", 0, "Index of the first result")
	addAllFlag(jiraIssueChangelogCmd)

	// remote link create flags
	jiraIssueRemoteLinkCreateCmd.Flags().String("url", "", "URL of the remote link")
	jiraIssueRemoteLinkCreateCmd.Flags().String("title", "", "Title of the remote link")

	// notify flags
	jiraIssueNotifyCmd.Flags().String("subject", "", "Notification subject")
	jiraIssueNotifyCmd.Flags().String("text-body", "", "Notification text body")

	// createmeta flags
	jiraIssueCreateMetaCmd.Flags().StringSlice("project", nil, "Project keys to filter")

	// Wire comment subcommands
	jiraIssueCommentCmd.AddCommand(jiraIssueCommentListCmd)
	jiraIssueCommentCmd.AddCommand(jiraIssueCommentAddCmd)
	jiraIssueCommentCmd.AddCommand(jiraIssueCommentGetCmd)
	jiraIssueCommentCmd.AddCommand(jiraIssueCommentDeleteCmd)

	// Wire worklog subcommands
	jiraIssueWorklogCmd.AddCommand(jiraIssueWorklogListCmd)
	jiraIssueWorklogCmd.AddCommand(jiraIssueWorklogAddCmd)
	jiraIssueWorklogCmd.AddCommand(jiraIssueWorklogDeleteCmd)

	// Wire remote link subcommands
	jiraIssueRemoteLinkCmd.AddCommand(jiraIssueRemoteLinkListCmd)
	jiraIssueRemoteLinkCmd.AddCommand(jiraIssueRemoteLinkCreateCmd)
	jiraIssueRemoteLinkCmd.AddCommand(jiraIssueRemoteLinkDeleteCmd)

	// Wire all issue subcommands
	jiraIssueCmd.AddCommand(jiraIssueListCmd)
	jiraIssueCmd.AddCommand(jiraIssueGetCmd)
	jiraIssueCmd.AddCommand(jiraIssueCreateCmd)
	jiraIssueCmd.AddCommand(jiraIssueEditCmd)
	jiraIssueCmd.AddCommand(jiraIssueDeleteCmd)
	jiraIssueCmd.AddCommand(jiraIssueAssignCmd)
	jiraIssueCmd.AddCommand(jiraIssueTransitionCmd)
	jiraIssueCmd.AddCommand(jiraIssueTransitionsCmd)
	jiraIssueCmd.AddCommand(jiraIssueCommentCmd)
	jiraIssueCmd.AddCommand(jiraIssueWorklogCmd)
	jiraIssueCmd.AddCommand(jiraIssueAttachCmd)
	jiraIssueCmd.AddCommand(jiraIssueVoteCmd)
	jiraIssueCmd.AddCommand(jiraIssueUnvoteCmd)
	jiraIssueCmd.AddCommand(jiraIssueWatchCmd)
	jiraIssueCmd.AddCommand(jiraIssueUnwatchCmd)
	jiraIssueCmd.AddCommand(jiraIssueWatchersCmd)
	jiraIssueCmd.AddCommand(jiraIssueChangelogCmd)
	jiraIssueCmd.AddCommand(jiraIssueRemoteLinkCmd)
	jiraIssueCmd.AddCommand(jiraIssueNotifyCmd)
	jiraIssueCmd.AddCommand(jiraIssueCreateMetaCmd)
	jiraIssueCmd.AddCommand(jiraIssueEditMetaCmd)
}
