package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// custom-content list
	listCustomContentCmd := &cobra.Command{
		Use:     "list",
		Short:   "List custom content by type",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if t := getStringFlag(cmd, "type"); t != "" {
				q.Set("type", t)
			}
			if ids := getStringSliceFlag(cmd, "id"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("id", id)
				}
			}
			if spaceIDs := getStringSliceFlag(cmd, "space-id"); len(spaceIDs) > 0 {
				for _, id := range spaceIDs {
					q.Add("space-id", id)
				}
			}
			data, err := confGet(cmd, "/custom-content", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(listCustomContentCmd)
	addSortFlag(listCustomContentCmd)
	addBodyFormatFlag(listCustomContentCmd)
	listCustomContentCmd.Flags().String("type", "", "Custom content type (required)")
	listCustomContentCmd.Flags().StringSlice("id", nil, "Filter by IDs")
	listCustomContentCmd.Flags().StringSlice("space-id", nil, "Filter by space IDs")
	confCustomContentCmd.AddCommand(listCustomContentCmd)

	// custom-content get
	getCustomContentCmd := &cobra.Command{
		Use:   "get [custom-content-id]",
		Short: "Get custom content by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if f := getStringFlag(cmd, "body-format"); f != "" {
				q.Set("body-format", f)
			}
			if v := getIntFlag(cmd, "version"); v > 0 {
				q.Set("version", fmt.Sprintf("%d", v))
			}
			for _, flag := range []string{"include-labels", "include-properties", "include-operations",
				"include-versions", "include-version", "include-collaborators"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/custom-content/"+args[0], q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addBodyFormatFlag(getCustomContentCmd)
	getCustomContentCmd.Flags().Int("version", 0, "Retrieve a specific version")
	getCustomContentCmd.Flags().Bool("include-labels", false, "Include labels")
	getCustomContentCmd.Flags().Bool("include-properties", false, "Include properties")
	getCustomContentCmd.Flags().Bool("include-operations", false, "Include operations")
	getCustomContentCmd.Flags().Bool("include-versions", false, "Include versions")
	getCustomContentCmd.Flags().Bool("include-version", false, "Include current version")
	getCustomContentCmd.Flags().Bool("include-collaborators", false, "Include collaborators")
	confCustomContentCmd.AddCommand(getCustomContentCmd)

	// custom-content create
	createCustomContentCmd := &cobra.Command{
		Use:   "create",
		Short: "Create custom content",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"type":  getStringFlag(cmd, "type"),
				"title": getStringFlag(cmd, "title"),
			}
			if s := getStringFlag(cmd, "status"); s != "" {
				body["status"] = s
			}
			if sid := getStringFlag(cmd, "space-id"); sid != "" {
				body["spaceId"] = sid
			}
			if pid := getStringFlag(cmd, "page-id"); pid != "" {
				body["pageId"] = pid
			}
			if bid := getStringFlag(cmd, "blogpost-id"); bid != "" {
				body["blogPostId"] = bid
			}
			if cid := getStringFlag(cmd, "custom-content-id"); cid != "" {
				body["customContentId"] = cid
			}
			if b := getStringFlag(cmd, "body"); b != "" {
				body["body"] = map[string]interface{}{
					"representation": getStringFlag(cmd, "body-format"),
					"value":          b,
				}
			}
			data, err := confPost(cmd, "/custom-content", nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	createCustomContentCmd.Flags().String("type", "", "Custom content type (required)")
	createCustomContentCmd.Flags().String("title", "", "Title (required)")
	createCustomContentCmd.Flags().String("status", "", "Status")
	createCustomContentCmd.Flags().String("space-id", "", "Space ID")
	createCustomContentCmd.Flags().String("page-id", "", "Page ID")
	createCustomContentCmd.Flags().String("blogpost-id", "", "Blog post ID")
	createCustomContentCmd.Flags().String("custom-content-id", "", "Parent custom content ID")
	createCustomContentCmd.Flags().String("body", "", "Body content")
	createCustomContentCmd.Flags().String("body-format", "storage", "Body format")
	_ = createCustomContentCmd.MarkFlagRequired("type")
	_ = createCustomContentCmd.MarkFlagRequired("title")
	confCustomContentCmd.AddCommand(createCustomContentCmd)

	// custom-content update
	updateCustomContentCmd := &cobra.Command{
		Use:   "update [custom-content-id]",
		Short: "Update custom content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"id":     args[0],
				"type":   getStringFlag(cmd, "type"),
				"status": getStringFlag(cmd, "status"),
				"title":  getStringFlag(cmd, "title"),
				"version": map[string]interface{}{
					"number":  getIntFlag(cmd, "version-number"),
					"message": getStringFlag(cmd, "version-message"),
				},
			}
			if b := getStringFlag(cmd, "body"); b != "" {
				body["body"] = map[string]interface{}{
					"representation": getStringFlag(cmd, "body-format"),
					"value":          b,
				}
			}
			if sid := getStringFlag(cmd, "space-id"); sid != "" {
				body["spaceId"] = sid
			}
			if pid := getStringFlag(cmd, "page-id"); pid != "" {
				body["pageId"] = pid
			}
			if bid := getStringFlag(cmd, "blogpost-id"); bid != "" {
				body["blogPostId"] = bid
			}
			if cid := getStringFlag(cmd, "custom-content-id"); cid != "" {
				body["customContentId"] = cid
			}
			data, err := confPut(cmd, "/custom-content/"+args[0], nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	updateCustomContentCmd.Flags().String("type", "", "Custom content type (required)")
	updateCustomContentCmd.Flags().String("title", "", "Title (required)")
	updateCustomContentCmd.Flags().String("status", "current", "Status (required)")
	updateCustomContentCmd.Flags().String("body", "", "Body content")
	updateCustomContentCmd.Flags().String("body-format", "storage", "Body format")
	updateCustomContentCmd.Flags().Int("version-number", 0, "Version number (required)")
	updateCustomContentCmd.Flags().String("version-message", "", "Version message")
	updateCustomContentCmd.Flags().String("space-id", "", "Space ID")
	updateCustomContentCmd.Flags().String("page-id", "", "Page ID")
	updateCustomContentCmd.Flags().String("blogpost-id", "", "Blog post ID")
	updateCustomContentCmd.Flags().String("custom-content-id", "", "Parent custom content ID")
	_ = updateCustomContentCmd.MarkFlagRequired("type")
	_ = updateCustomContentCmd.MarkFlagRequired("title")
	_ = updateCustomContentCmd.MarkFlagRequired("version-number")
	confCustomContentCmd.AddCommand(updateCustomContentCmd)

	// custom-content delete
	deleteCustomContentCmd := &cobra.Command{
		Use:   "delete [custom-content-id]",
		Short: "Delete custom content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if getBoolFlag(cmd, "purge") {
				q.Set("purge", "true")
			}
			_, err := confDelete(cmd, "/custom-content/"+args[0], q)
			if err != nil {
				return err
			}
			fmt.Println("Custom content deleted successfully.")
			return nil
		},
	}
	deleteCustomContentCmd.Flags().Bool("purge", false, "Purge the custom content")
	confCustomContentCmd.AddCommand(deleteCustomContentCmd)

	// custom-content sub-resources
	for _, sub := range []struct {
		use, short, path string
	}{
		{"attachments [id]", "Get attachments for custom content", "/attachments"},
		{"children [id]", "Get child custom content", "/children"},
		{"labels [id]", "Get labels for custom content", "/labels"},
		{"comments [id]", "Get footer comments for custom content", "/footer-comments"},
		{"operations [id]", "Get permitted operations", "/operations"},
		{"versions [id]", "Get custom content versions", "/versions"},
	} {
		sub := sub
		subCmd := &cobra.Command{
			Use:   sub.use,
			Short: sub.short,
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				q := getPaginationQuery(cmd)
				data, err := confGet(cmd, "/custom-content/"+args[0]+sub.path, q)
				if err != nil {
					return err
				}
				printJSON(data)
				return nil
			},
		}
		addPaginationFlags(subCmd)
		confCustomContentCmd.AddCommand(subCmd)
	}

	// custom-content version-details
	ccVersionDetailCmd := &cobra.Command{
		Use:   "version-details [custom-content-id] [version-number]",
		Short: "Get version details for custom content version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/custom-content/"+args[0]+"/versions/"+args[1], nil)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	confCustomContentCmd.AddCommand(ccVersionDetailCmd)
}
