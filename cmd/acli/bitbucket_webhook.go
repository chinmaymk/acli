package acli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbWebhookCmd = &cobra.Command{
	Use:     "webhook",
	Short:   "Manage webhooks",
	Aliases: []string{"hook"},
	RunE:    helpRunE,
}

func init() {
	// webhook list (repo-level)
	bbWebhookCmd.AddCommand(&cobra.Command{
		Use:     "list [workspace] <repo-slug>",
		Short:   "List webhooks for a repository",
		Aliases: []string{"ls"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, err := resolveWorkspaceAndRepo(cmd, args)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			hooks, err := client.ListRepoWebhooks(workspace, repoSlug)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "UUID\tDESCRIPTION\tURL\tACTIVE\tEVENTS")
			for _, h := range hooks {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n",
					h.UUID, h.Description, h.URL, h.Active, strings.Join(h.Events, ","))
			}
			return w.Flush()
		},
	})

	// webhook get
	bbWebhookCmd.AddCommand(&cobra.Command{
		Use:   "get [workspace] <repo-slug> <webhook-uuid>",
		Short: "Get webhook details",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, hookUUID, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			hook, err := client.GetRepoWebhook(workspace, repoSlug, hookUUID)
			if err != nil {
				return err
			}

			fmt.Printf("UUID:         %s\n", hook.UUID)
			fmt.Printf("Description:  %s\n", hook.Description)
			fmt.Printf("URL:          %s\n", hook.URL)
			fmt.Printf("Active:       %v\n", hook.Active)
			fmt.Printf("Events:       %s\n", strings.Join(hook.Events, ", "))
			return nil
		},
	})

	// webhook create
	webhookCreateCmd := &cobra.Command{
		Use:   "create [workspace] <repo-slug>",
		Short: "Create a webhook for a repository",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, err := resolveWorkspaceAndRepo(cmd, args)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			webhookURL, _ := cmd.Flags().GetString("url")
			desc, _ := cmd.Flags().GetString("description")
			events, _ := cmd.Flags().GetStringSlice("events")
			active, _ := cmd.Flags().GetBool("active")

			if webhookURL == "" || len(events) == 0 {
				return fmt.Errorf("--url and --events are required")
			}

			hook, err := client.CreateRepoWebhook(workspace, repoSlug, &bitbucket.CreateWebhookRequest{
				Description: desc,
				URL:         webhookURL,
				Active:      active,
				Events:      events,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Created webhook: %s\n", hook.UUID)
			fmt.Printf("URL: %s\n", hook.URL)
			return nil
		},
	}
	webhookCreateCmd.Flags().String("url", "", "Webhook URL (required)")
	webhookCreateCmd.Flags().String("description", "", "Webhook description")
	webhookCreateCmd.Flags().StringSlice("events", nil, "Events to subscribe to (required, e.g. repo:push,pullrequest:created)")
	webhookCreateCmd.Flags().Bool("active", true, "Whether the webhook is active")
	bbWebhookCmd.AddCommand(webhookCreateCmd)

	// webhook delete
	bbWebhookCmd.AddCommand(&cobra.Command{
		Use:   "delete [workspace] <repo-slug> <webhook-uuid>",
		Short: "Delete a webhook",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, repoSlug, hookUUID, err := resolveWorkspaceRepoAndID(cmd, args)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeleteRepoWebhook(workspace, repoSlug, hookUUID); err != nil {
				return err
			}
			fmt.Printf("Deleted webhook: %s\n", hookUUID)
			return nil
		},
	})

	// webhook list-workspace
	bbWebhookCmd.AddCommand(&cobra.Command{
		Use:   "list-workspace [workspace]",
		Short: "List webhooks for a workspace",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, err := defaultWorkspace(cmd, args, 0)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			hooks, err := client.ListWorkspaceWebhooks(workspace)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "UUID\tDESCRIPTION\tURL\tACTIVE\tEVENTS")
			for _, h := range hooks {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n",
					h.UUID, h.Description, h.URL, h.Active, strings.Join(h.Events, ","))
			}
			return w.Flush()
		},
	})

	// webhook create-workspace
	wsWebhookCreateCmd := &cobra.Command{
		Use:   "create-workspace [workspace]",
		Short: "Create a webhook for a workspace",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, err := defaultWorkspace(cmd, args, 0)
			if err != nil {
				return err
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}

			webhookURL, _ := cmd.Flags().GetString("url")
			desc, _ := cmd.Flags().GetString("description")
			events, _ := cmd.Flags().GetStringSlice("events")
			active, _ := cmd.Flags().GetBool("active")

			if webhookURL == "" || len(events) == 0 {
				return fmt.Errorf("--url and --events are required")
			}

			hook, err := client.CreateWorkspaceWebhook(workspace, &bitbucket.CreateWebhookRequest{
				Description: desc,
				URL:         webhookURL,
				Active:      active,
				Events:      events,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Created workspace webhook: %s\n", hook.UUID)
			return nil
		},
	}
	wsWebhookCreateCmd.Flags().String("url", "", "Webhook URL (required)")
	wsWebhookCreateCmd.Flags().String("description", "", "Webhook description")
	wsWebhookCreateCmd.Flags().StringSlice("events", nil, "Events to subscribe to (required)")
	wsWebhookCreateCmd.Flags().Bool("active", true, "Whether the webhook is active")
	bbWebhookCmd.AddCommand(wsWebhookCreateCmd)

	// webhook delete-workspace
	bbWebhookCmd.AddCommand(&cobra.Command{
		Use:   "delete-workspace [workspace] <webhook-uuid>",
		Short: "Delete a workspace webhook",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var workspace, hookUUID string
			if len(args) >= 2 {
				workspace = args[0]
				hookUUID = args[1]
			} else {
				var err error
				workspace, err = defaultWorkspace(cmd, nil, 0)
				if err != nil {
					return err
				}
				hookUUID = args[0]
			}
			client, err := getBitbucketClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeleteWorkspaceWebhook(workspace, hookUUID); err != nil {
				return err
			}
			fmt.Printf("Deleted workspace webhook: %s\n", hookUUID)
			return nil
		},
	})
}
