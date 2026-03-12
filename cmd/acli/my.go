package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
	Use:   "my",
	Short: "Show your work across Jira and Bitbucket",
	Long: `Show a dashboard of your current work: assigned Jira issues,
your open pull requests, and pull requests awaiting your review.`,
	RunE: helpRunE,
}

func init() {
	// my issues
	myIssuesCmd := &cobra.Command{
		Use:     "issues",
		Short:   "List Jira issues assigned to you",
		Aliases: []string{"tickets", "i"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")
			project, _ := cmd.Flags().GetString("project")
			maxResults, _ := cmd.Flags().GetInt("max")

			jql := "assignee = currentUser() ORDER BY updated DESC"
			if status != "" {
				jql = fmt.Sprintf("assignee = currentUser() AND status = %q ORDER BY updated DESC", status)
			}
			if project != "" {
				if status != "" {
					jql = fmt.Sprintf("assignee = currentUser() AND project = %q AND status = %q ORDER BY updated DESC", project, status)
				} else {
					jql = fmt.Sprintf("assignee = currentUser() AND project = %q ORDER BY updated DESC", project)
				}
			}

			results, err := client.SearchJQL(jql, 0, maxResults,
				[]string{"summary", "status", "priority", "issuetype", "project", "updated"}, nil)
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return outputJSON(results.Issues)
			}

			if len(results.Issues) == 0 {
				fmt.Println("No issues assigned to you.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "KEY\tTYPE\tSTATUS\tPRIORITY\tSUMMARY")
			for _, issue := range results.Issues {
				printIssueRow(w, issue)
			}
			printPaginationHint(cmd, len(results.Issues), results.Total)
			return w.Flush()
		},
	}
	myIssuesCmd.Flags().String("status", "", "Filter by status (e.g. 'In Progress', 'To Do')")
	myIssuesCmd.Flags().String("project", "", "Filter by project key")
	myIssuesCmd.Flags().Int("max", 25, "Maximum number of results")
	myCmd.AddCommand(myIssuesCmd)

	// my prs
	myPRsCmd := &cobra.Command{
		Use:     "prs",
		Short:   "List your open pull requests in Bitbucket",
		Aliases: []string{"pulls", "pr"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			workspace, _ := cmd.Flags().GetString("workspace")
			repoSlug, _ := cmd.Flags().GetString("repo")

			if workspace == "" {
				profile, err := getProfile(cmd)
				if err != nil {
					return err
				}
				workspace = profile.Defaults.Workspace
			}

			if workspace == "" {
				return fmt.Errorf("workspace is required: use --workspace or set a default with 'acli config set-defaults'")
			}

			if repoSlug != "" {
				return listMyPRsForRepo(cmd, client, workspace, repoSlug)
			}

			// List PRs across all repos the user has access to
			repos, err := client.ListRepositories(workspace, &bitbucket.ListReposOptions{
				Role: "member",
			})
			if err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				var allPRs []bitbucket.PullRequest
				for _, repo := range repos {
					prs, err := client.ListPullRequests(workspace, repo.Slug, &bitbucket.ListPRsOptions{
						State: "OPEN",
					})
					if err != nil {
						continue
					}
					// Filter to only user's PRs
					user, _ := client.GetCurrentUser()
					for _, pr := range prs {
						if user != nil && pr.Author.UUID == user.UUID {
							allPRs = append(allPRs, pr)
						}
					}
				}
				return outputJSON(allPRs)
			}

			user, err := client.GetCurrentUser()
			if err != nil {
				return fmt.Errorf("getting current user: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "REPO\tID\tTITLE\tSTATE\tSOURCE\tDESTINATION")
			count := 0
			for _, repo := range repos {
				prs, err := client.ListPullRequests(workspace, repo.Slug, &bitbucket.ListPRsOptions{
					State: "OPEN",
				})
				if err != nil {
					continue
				}
				for _, pr := range prs {
					if pr.Author.UUID == user.UUID {
						fmt.Fprintf(w, "%s\t#%d\t%s\t%s\t%s\t%s\n",
							repo.Slug, pr.ID, truncate(pr.Title, 50), pr.State,
							pr.Source.Branch.Name, pr.Destination.Branch.Name)
						count++
					}
				}
			}
			if count == 0 {
				fmt.Println("No open pull requests.")
				return nil
			}
			return w.Flush()
		},
	}
	myPRsCmd.Flags().String("workspace", "", "Bitbucket workspace (defaults to profile default)")
	myPRsCmd.Flags().String("repo", "", "Filter to a specific repository")
	myCmd.AddCommand(myPRsCmd)

	// my reviews
	myReviewsCmd := &cobra.Command{
		Use:     "reviews",
		Short:   "List pull requests awaiting your review",
		Aliases: []string{"review", "r"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			workspace, _ := cmd.Flags().GetString("workspace")
			repoSlug, _ := cmd.Flags().GetString("repo")

			if workspace == "" {
				profile, err := getProfile(cmd)
				if err != nil {
					return err
				}
				workspace = profile.Defaults.Workspace
			}

			if workspace == "" {
				return fmt.Errorf("workspace is required: use --workspace or set a default with 'acli config set-defaults'")
			}

			user, err := client.GetCurrentUser()
			if err != nil {
				return fmt.Errorf("getting current user: %w", err)
			}

			var repos []bitbucket.Repository
			if repoSlug != "" {
				repos = []bitbucket.Repository{{Slug: repoSlug}}
			} else {
				repos, err = client.ListRepositories(workspace, &bitbucket.ListReposOptions{
					Role: "member",
				})
				if err != nil {
					return err
				}
			}

			type reviewPR struct {
				Repo string
				PR   bitbucket.PullRequest
			}
			var reviewPRs []reviewPR
			for _, repo := range repos {
				prs, err := client.ListPullRequests(workspace, repo.Slug, &bitbucket.ListPRsOptions{
					State: "OPEN",
				})
				if err != nil {
					continue
				}
				for _, pr := range prs {
					for _, reviewer := range pr.Reviewers {
						if reviewer.UUID == user.UUID {
							reviewPRs = append(reviewPRs, reviewPR{Repo: repo.Slug, PR: pr})
							break
						}
					}
				}
			}

			if isJSONOutput(cmd) {
				return outputJSON(reviewPRs)
			}

			if len(reviewPRs) == 0 {
				fmt.Println("No pull requests awaiting your review.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "REPO\tID\tTITLE\tAUTHOR\tSOURCE\tDESTINATION")
			for _, rpr := range reviewPRs {
				fmt.Fprintf(w, "%s\t#%d\t%s\t%s\t%s\t%s\n",
					rpr.Repo, rpr.PR.ID, truncate(rpr.PR.Title, 50), rpr.PR.Author.DisplayName,
					rpr.PR.Source.Branch.Name, rpr.PR.Destination.Branch.Name)
			}
			return w.Flush()
		},
	}
	myReviewsCmd.Flags().String("workspace", "", "Bitbucket workspace (defaults to profile default)")
	myReviewsCmd.Flags().String("repo", "", "Filter to a specific repository")
	myCmd.AddCommand(myReviewsCmd)
}

func listMyPRsForRepo(cmd *cobra.Command, client *bitbucket.Client, workspace, repoSlug string) error {
	prs, err := client.ListPullRequests(workspace, repoSlug, &bitbucket.ListPRsOptions{
		State: "OPEN",
	})
	if err != nil {
		return err
	}

	user, err := client.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("getting current user: %w", err)
	}

	var myPRs []bitbucket.PullRequest
	for _, pr := range prs {
		if pr.Author.UUID == user.UUID {
			myPRs = append(myPRs, pr)
		}
	}

	if isJSONOutput(cmd) {
		return outputJSON(myPRs)
	}

	if len(myPRs) == 0 {
		fmt.Println("No open pull requests.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tSTATE\tSOURCE\tDESTINATION")
	for _, pr := range myPRs {
		fmt.Fprintf(w, "#%d\t%s\t%s\t%s\t%s\n",
			pr.ID, truncate(pr.Title, 50), pr.State,
			pr.Source.Branch.Name, pr.Destination.Branch.Name)
	}
	return w.Flush()
}

// myDashboardHelper is used internally to initialize both clients.
func getJiraAndBitbucketClients(cmd *cobra.Command) (*jira.Client, *bitbucket.Client, error) {
	profile, err := getProfile(cmd)
	if err != nil {
		return nil, nil, err
	}
	jiraClient, err := jira.NewClient(profile)
	if err != nil {
		return nil, nil, err
	}
	bbClient, err := bitbucket.NewClient(profile)
	if err != nil {
		return nil, nil, err
	}
	return jiraClient, bbClient, nil
}
