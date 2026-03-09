package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// page list
	listPagesCmd := &cobra.Command{
		Use:     "list",
		Short:   "List pages",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
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
			if t := getStringFlag(cmd, "title"); t != "" {
				q.Set("title", t)
			}
			if s := getStringFlag(cmd, "subtype"); s != "" {
				q.Set("subtype", s)
			}
			data, err := confGet(cmd, "/pages", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(listPagesCmd)
	addSortFlag(listPagesCmd)
	addStatusFlag(listPagesCmd)
	addBodyFormatFlag(listPagesCmd)
	listPagesCmd.Flags().StringSlice("id", nil, "Filter by page IDs")
	listPagesCmd.Flags().StringSlice("space-id", nil, "Filter by space IDs")
	listPagesCmd.Flags().String("title", "", "Filter by title")
	listPagesCmd.Flags().String("subtype", "", "Filter by subtype")
	confPageCmd.AddCommand(listPagesCmd)

	// page get
	getPageCmd := &cobra.Command{
		Use:   "get [page-id]",
		Short: "Get page by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if f := getStringFlag(cmd, "body-format"); f != "" {
				q.Set("body-format", f)
			}
			if getBoolFlag(cmd, "get-draft") {
				q.Set("get-draft", "true")
			}
			if v := getIntFlag(cmd, "version"); v > 0 {
				q.Set("version", fmt.Sprintf("%d", v))
			}
			for _, flag := range []string{"include-labels", "include-properties", "include-operations",
				"include-likes", "include-versions", "include-version",
				"include-favorited-by-current-user-status", "include-webresources", "include-collaborators", "include-direct-children"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/pages/"+args[0], q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addBodyFormatFlag(getPageCmd)
	addStatusFlag(getPageCmd)
	getPageCmd.Flags().Bool("get-draft", false, "Retrieve draft version")
	getPageCmd.Flags().Int("version", 0, "Retrieve a specific version")
	getPageCmd.Flags().Bool("include-labels", false, "Include labels")
	getPageCmd.Flags().Bool("include-properties", false, "Include properties")
	getPageCmd.Flags().Bool("include-operations", false, "Include operations")
	getPageCmd.Flags().Bool("include-likes", false, "Include likes")
	getPageCmd.Flags().Bool("include-versions", false, "Include versions")
	getPageCmd.Flags().Bool("include-version", false, "Include current version")
	getPageCmd.Flags().Bool("include-favorited-by-current-user-status", false, "Include favorited status")
	getPageCmd.Flags().Bool("include-webresources", false, "Include web resources")
	getPageCmd.Flags().Bool("include-collaborators", false, "Include collaborators")
	getPageCmd.Flags().Bool("include-direct-children", false, "Include direct children")
	confPageCmd.AddCommand(getPageCmd)

	// page create
	createPageCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a page",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if getBoolFlag(cmd, "embedded") {
				q.Set("embedded", "true")
			}
			if getBoolFlag(cmd, "private") {
				q.Set("private", "true")
			}
			if getBoolFlag(cmd, "root-level") {
				q.Set("root-level", "true")
			}

			body := map[string]interface{}{
				"spaceId": getStringFlag(cmd, "space-id"),
			}
			if t := getStringFlag(cmd, "title"); t != "" {
				body["title"] = t
			}
			if s := getStringFlag(cmd, "status"); s != "" {
				body["status"] = s
			}
			if p := getStringFlag(cmd, "parent-id"); p != "" {
				body["parentId"] = p
			}
			if sub := getStringFlag(cmd, "subtype"); sub != "" {
				body["subtype"] = sub
			}
			if b := getStringFlag(cmd, "body"); b != "" {
				body["body"] = map[string]interface{}{
					"representation": getStringFlag(cmd, "body-format"),
					"value":          b,
				}
			}
			data, err := confPost(cmd, "/pages", q, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	createPageCmd.Flags().String("space-id", "", "Space ID (required)")
	createPageCmd.Flags().String("title", "", "Page title")
	createPageCmd.Flags().String("parent-id", "", "Parent page ID")
	createPageCmd.Flags().String("status", "", "Page status (current, draft)")
	createPageCmd.Flags().String("body", "", "Page body content")
	createPageCmd.Flags().String("body-format", "storage", "Body format (storage, atlas_doc_format, wiki)")
	createPageCmd.Flags().String("subtype", "", "Page subtype")
	createPageCmd.Flags().Bool("embedded", false, "Create as embedded content")
	createPageCmd.Flags().Bool("private", false, "Create as private page")
	createPageCmd.Flags().Bool("root-level", false, "Create at root level of space")
	_ = createPageCmd.MarkFlagRequired("space-id")
	confPageCmd.AddCommand(createPageCmd)

	// page update
	updatePageCmd := &cobra.Command{
		Use:   "update [page-id]",
		Short: "Update a page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"id":     args[0],
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
			if pid := getStringFlag(cmd, "parent-id"); pid != "" {
				body["parentId"] = pid
			}
			data, err := confPut(cmd, "/pages/"+args[0], nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	updatePageCmd.Flags().String("title", "", "Page title (required)")
	updatePageCmd.Flags().String("status", "current", "Page status (required)")
	updatePageCmd.Flags().String("body", "", "Page body content")
	updatePageCmd.Flags().String("body-format", "storage", "Body format")
	updatePageCmd.Flags().Int("version-number", 0, "Version number (required)")
	updatePageCmd.Flags().String("version-message", "", "Version message")
	updatePageCmd.Flags().String("space-id", "", "Space ID")
	updatePageCmd.Flags().String("parent-id", "", "Parent page ID")
	_ = updatePageCmd.MarkFlagRequired("title")
	_ = updatePageCmd.MarkFlagRequired("version-number")
	confPageCmd.AddCommand(updatePageCmd)

	// page update-title
	updatePageTitleCmd := &cobra.Command{
		Use:   "update-title [page-id]",
		Short: "Update page title only",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"title":  getStringFlag(cmd, "title"),
				"status": getStringFlag(cmd, "status"),
			}
			data, err := confPut(cmd, "/pages/"+args[0]+"/title", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	updatePageTitleCmd.Flags().String("title", "", "New title (required)")
	updatePageTitleCmd.Flags().String("status", "current", "Page status (required)")
	_ = updatePageTitleCmd.MarkFlagRequired("title")
	confPageCmd.AddCommand(updatePageTitleCmd)

	// page delete
	deletePageCmd := &cobra.Command{
		Use:   "delete [page-id]",
		Short: "Delete a page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if getBoolFlag(cmd, "purge") {
				q.Set("purge", "true")
			}
			if getBoolFlag(cmd, "draft") {
				q.Set("draft", "true")
			}
			_, err := confDelete(cmd, "/pages/"+args[0], q)
			if err != nil {
				return err
			}
			fmt.Println("Page deleted successfully.")
			return nil
		},
	}
	deletePageCmd.Flags().Bool("purge", false, "Purge the page")
	deletePageCmd.Flags().Bool("draft", false, "Delete a draft page")
	confPageCmd.AddCommand(deletePageCmd)

	// page children
	childPagesCmd := &cobra.Command{
		Use:   "children [page-id]",
		Short: "Get child pages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/pages/"+args[0]+"/children", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(childPagesCmd)
	addSortFlag(childPagesCmd)
	confPageCmd.AddCommand(childPagesCmd)

	// page direct-children
	directChildrenCmd := &cobra.Command{
		Use:   "direct-children [page-id]",
		Short: "Get direct children of a page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/pages/"+args[0]+"/direct-children", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(directChildrenCmd)
	addSortFlag(directChildrenCmd)
	confPageCmd.AddCommand(directChildrenCmd)

	// page ancestors
	ancestorsCmd := &cobra.Command{
		Use:   "ancestors [page-id]",
		Short: "Get all ancestors of page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit := getIntFlag(cmd, "limit"); limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			data, err := confGet(cmd, "/pages/"+args[0]+"/ancestors", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	ancestorsCmd.Flags().Int("limit", 25, "Maximum number of results")
	confPageCmd.AddCommand(ancestorsCmd)

	// page descendants
	descendantsCmd := &cobra.Command{
		Use:   "descendants [page-id]",
		Short: "Get descendants of page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if d := getIntFlag(cmd, "depth"); d > 0 {
				q.Set("depth", fmt.Sprintf("%d", d))
			}
			data, err := confGet(cmd, "/pages/"+args[0]+"/descendants", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(descendantsCmd)
	descendantsCmd.Flags().Int("depth", 0, "Maximum depth of descendants")
	confPageCmd.AddCommand(descendantsCmd)

	// page versions
	pageVersionsCmd := &cobra.Command{
		Use:   "versions [page-id]",
		Short: "Get page versions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/pages/"+args[0]+"/versions", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(pageVersionsCmd)
	addSortFlag(pageVersionsCmd)
	addBodyFormatFlag(pageVersionsCmd)
	confPageCmd.AddCommand(pageVersionsCmd)

	// page version-details
	pageVersionDetailCmd := &cobra.Command{
		Use:   "version-details [page-id] [version-number]",
		Short: "Get version details for page version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/pages/"+args[0]+"/versions/"+args[1], nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confPageCmd.AddCommand(pageVersionDetailCmd)

	// page labels
	pageLabelsCmd := &cobra.Command{
		Use:   "labels [page-id]",
		Short: "Get labels for page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if p := getStringFlag(cmd, "prefix"); p != "" {
				q.Set("prefix", p)
			}
			data, err := confGet(cmd, "/pages/"+args[0]+"/labels", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(pageLabelsCmd)
	addSortFlag(pageLabelsCmd)
	pageLabelsCmd.Flags().String("prefix", "", "Filter by prefix")
	confPageCmd.AddCommand(pageLabelsCmd)

	// page attachments
	pageAttachmentsCmd := &cobra.Command{
		Use:   "attachments [page-id]",
		Short: "Get attachments for page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if m := getStringFlag(cmd, "media-type"); m != "" {
				q.Set("mediaType", m)
			}
			if f := getStringFlag(cmd, "filename"); f != "" {
				q.Set("filename", f)
			}
			data, err := confGet(cmd, "/pages/"+args[0]+"/attachments", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(pageAttachmentsCmd)
	addSortFlag(pageAttachmentsCmd)
	addStatusFlag(pageAttachmentsCmd)
	pageAttachmentsCmd.Flags().String("media-type", "", "Filter by media type")
	pageAttachmentsCmd.Flags().String("filename", "", "Filter by filename")
	confPageCmd.AddCommand(pageAttachmentsCmd)

	// page footer-comments
	pageFooterCommentsCmd := &cobra.Command{
		Use:   "footer-comments [page-id]",
		Short: "Get footer comments for page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/pages/"+args[0]+"/footer-comments", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(pageFooterCommentsCmd)
	addSortFlag(pageFooterCommentsCmd)
	addStatusFlag(pageFooterCommentsCmd)
	addBodyFormatFlag(pageFooterCommentsCmd)
	confPageCmd.AddCommand(pageFooterCommentsCmd)

	// page inline-comments
	pageInlineCommentsCmd := &cobra.Command{
		Use:   "inline-comments [page-id]",
		Short: "Get inline comments for page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if r := getStringSliceFlag(cmd, "resolution-status"); len(r) > 0 {
				for _, rs := range r {
					q.Add("resolution-status", rs)
				}
			}
			data, err := confGet(cmd, "/pages/"+args[0]+"/inline-comments", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(pageInlineCommentsCmd)
	addSortFlag(pageInlineCommentsCmd)
	addStatusFlag(pageInlineCommentsCmd)
	addBodyFormatFlag(pageInlineCommentsCmd)
	pageInlineCommentsCmd.Flags().StringSlice("resolution-status", nil, "Filter by resolution status")
	confPageCmd.AddCommand(pageInlineCommentsCmd)

	// page custom-content
	pageCustomContentCmd := &cobra.Command{
		Use:   "custom-content [page-id]",
		Short: "Get custom content by type in page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if t := getStringFlag(cmd, "type"); t != "" {
				q.Set("type", t)
			}
			data, err := confGet(cmd, "/pages/"+args[0]+"/custom-content", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(pageCustomContentCmd)
	addSortFlag(pageCustomContentCmd)
	addBodyFormatFlag(pageCustomContentCmd)
	pageCustomContentCmd.Flags().String("type", "", "Custom content type")
	confPageCmd.AddCommand(pageCustomContentCmd)

	// page operations
	pageOpsCmd := &cobra.Command{
		Use:   "operations [page-id]",
		Short: "Get permitted operations for page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/pages/"+args[0]+"/operations", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confPageCmd.AddCommand(pageOpsCmd)

	// page likes count
	pageLikesCountCmd := &cobra.Command{
		Use:   "likes-count [page-id]",
		Short: "Get like count for page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/pages/"+args[0]+"/likes/count", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confPageCmd.AddCommand(pageLikesCountCmd)

	// page likes users
	pageLikesUsersCmd := &cobra.Command{
		Use:   "likes-users [page-id]",
		Short: "Get account IDs of likes for page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/pages/"+args[0]+"/likes/users", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(pageLikesUsersCmd)
	confPageCmd.AddCommand(pageLikesUsersCmd)

	// page redact
	pageRedactCmd := &cobra.Command{
		Use:   "redact [page-id]",
		Short: "Redact content in a Confluence page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bodyStr := getStringFlag(cmd, "body")
			if bodyStr == "" {
				return fmt.Errorf("--body is required (JSON redaction request)")
			}
			var body interface{}
			if err := parseJSONFlag(bodyStr, &body); err != nil {
				return err
			}
			data, err := confPost(cmd, "/pages/"+args[0]+"/redact", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	pageRedactCmd.Flags().String("body", "", "JSON redaction request body")
	confPageCmd.AddCommand(pageRedactCmd)
}
