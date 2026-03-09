package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// === Footer Comments ===

	// footer list
	listFooterCmd := &cobra.Command{
		Use:     "list",
		Short:   "List footer comments",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/footer-comments", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(listFooterCmd)
	addSortFlag(listFooterCmd)
	addBodyFormatFlag(listFooterCmd)
	confFooterCommentCmd.AddCommand(listFooterCmd)

	// footer get
	getFooterCmd := &cobra.Command{
		Use:   "get [comment-id]",
		Short: "Get footer comment by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if f := getStringFlag(cmd, "body-format"); f != "" {
				q.Set("body-format", f)
			}
			if v := getIntFlag(cmd, "version"); v > 0 {
				q.Set("version", fmt.Sprintf("%d", v))
			}
			for _, flag := range []string{"include-properties", "include-operations", "include-likes", "include-versions", "include-version"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/footer-comments/"+args[0], q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addBodyFormatFlag(getFooterCmd)
	getFooterCmd.Flags().Int("version", 0, "Retrieve a specific version")
	getFooterCmd.Flags().Bool("include-properties", false, "Include properties")
	getFooterCmd.Flags().Bool("include-operations", false, "Include operations")
	getFooterCmd.Flags().Bool("include-likes", false, "Include likes")
	getFooterCmd.Flags().Bool("include-versions", false, "Include versions")
	getFooterCmd.Flags().Bool("include-version", false, "Include current version")
	confFooterCommentCmd.AddCommand(getFooterCmd)

	// footer create
	createFooterCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a footer comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{}
			if pid := getStringFlag(cmd, "page-id"); pid != "" {
				body["pageId"] = pid
			}
			if bid := getStringFlag(cmd, "blogpost-id"); bid != "" {
				body["blogPostId"] = bid
			}
			if aid := getStringFlag(cmd, "attachment-id"); aid != "" {
				body["attachmentId"] = aid
			}
			if cid := getStringFlag(cmd, "custom-content-id"); cid != "" {
				body["customContentId"] = cid
			}
			if pcid := getStringFlag(cmd, "parent-comment-id"); pcid != "" {
				body["parentCommentId"] = pcid
			}
			if b := getStringFlag(cmd, "body"); b != "" {
				body["body"] = map[string]interface{}{
					"representation": getStringFlag(cmd, "body-format"),
					"value":          b,
				}
			}
			data, err := confPost(cmd, "/footer-comments", nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	createFooterCmd.Flags().String("page-id", "", "Page ID to comment on")
	createFooterCmd.Flags().String("blogpost-id", "", "Blog post ID to comment on")
	createFooterCmd.Flags().String("attachment-id", "", "Attachment ID to comment on")
	createFooterCmd.Flags().String("custom-content-id", "", "Custom content ID to comment on")
	createFooterCmd.Flags().String("parent-comment-id", "", "Parent comment ID (for replies)")
	createFooterCmd.Flags().String("body", "", "Comment body content")
	createFooterCmd.Flags().String("body-format", "storage", "Body format")
	confFooterCommentCmd.AddCommand(createFooterCmd)

	// footer update
	updateFooterCmd := &cobra.Command{
		Use:   "update [comment-id]",
		Short: "Update a footer comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
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
			data, err := confPut(cmd, "/footer-comments/"+args[0], nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	updateFooterCmd.Flags().Int("version-number", 0, "Version number (required)")
	updateFooterCmd.Flags().String("version-message", "", "Version message")
	updateFooterCmd.Flags().String("body", "", "Comment body content")
	updateFooterCmd.Flags().String("body-format", "storage", "Body format")
	_ = updateFooterCmd.MarkFlagRequired("version-number")
	confFooterCommentCmd.AddCommand(updateFooterCmd)

	// footer delete
	deleteFooterCmd := &cobra.Command{
		Use:   "delete [comment-id]",
		Short: "Delete a footer comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/footer-comments/"+args[0], nil)
			if err != nil {
				return err
			}
			fmt.Println("Footer comment deleted successfully.")
			return nil
		},
	}
	confFooterCommentCmd.AddCommand(deleteFooterCmd)

	// footer children
	footerChildrenCmd := &cobra.Command{
		Use:   "children [comment-id]",
		Short: "Get children footer comments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/footer-comments/"+args[0]+"/children", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(footerChildrenCmd)
	confFooterCommentCmd.AddCommand(footerChildrenCmd)

	// footer sub-resources
	for _, sub := range []struct {
		use, short, path string
	}{
		{"operations [comment-id]", "Get permitted operations", "/operations"},
		{"versions [comment-id]", "Get footer comment versions", "/versions"},
		{"likes-count [comment-id]", "Get like count", "/likes/count"},
		{"likes-users [comment-id]", "Get like users", "/likes/users"},
	} {
		sub := sub
		subCmd := &cobra.Command{
			Use:   sub.use,
			Short: sub.short,
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := confGet(cmd, "/footer-comments/"+args[0]+sub.path, nil)
				if err != nil {
					return err
				}
				printJSON(data)
				return nil
			},
		}
		confFooterCommentCmd.AddCommand(subCmd)
	}

	// footer version-details
	footerVersionDetailCmd := &cobra.Command{
		Use:   "version-details [comment-id] [version-number]",
		Short: "Get version details for footer comment version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/footer-comments/"+args[0]+"/versions/"+args[1], nil)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	confFooterCommentCmd.AddCommand(footerVersionDetailCmd)

	// === Inline Comments ===

	// inline list
	listInlineCmd := &cobra.Command{
		Use:     "list",
		Short:   "List inline comments",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/inline-comments", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(listInlineCmd)
	addSortFlag(listInlineCmd)
	addBodyFormatFlag(listInlineCmd)
	confInlineCommentCmd.AddCommand(listInlineCmd)

	// inline get
	getInlineCmd := &cobra.Command{
		Use:   "get [comment-id]",
		Short: "Get inline comment by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if f := getStringFlag(cmd, "body-format"); f != "" {
				q.Set("body-format", f)
			}
			if v := getIntFlag(cmd, "version"); v > 0 {
				q.Set("version", fmt.Sprintf("%d", v))
			}
			for _, flag := range []string{"include-properties", "include-operations", "include-likes", "include-versions", "include-version"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/inline-comments/"+args[0], q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addBodyFormatFlag(getInlineCmd)
	getInlineCmd.Flags().Int("version", 0, "Retrieve a specific version")
	getInlineCmd.Flags().Bool("include-properties", false, "Include properties")
	getInlineCmd.Flags().Bool("include-operations", false, "Include operations")
	getInlineCmd.Flags().Bool("include-likes", false, "Include likes")
	getInlineCmd.Flags().Bool("include-versions", false, "Include versions")
	getInlineCmd.Flags().Bool("include-version", false, "Include current version")
	confInlineCommentCmd.AddCommand(getInlineCmd)

	// inline create
	createInlineCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an inline comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{}
			if pid := getStringFlag(cmd, "page-id"); pid != "" {
				body["pageId"] = pid
			}
			if bid := getStringFlag(cmd, "blogpost-id"); bid != "" {
				body["blogPostId"] = bid
			}
			if pcid := getStringFlag(cmd, "parent-comment-id"); pcid != "" {
				body["parentCommentId"] = pcid
			}
			if b := getStringFlag(cmd, "body"); b != "" {
				body["body"] = map[string]interface{}{
					"representation": getStringFlag(cmd, "body-format"),
					"value":          b,
				}
			}
			if props := getStringFlag(cmd, "inline-comment-properties"); props != "" {
				var p interface{}
				if err := parseJSONFlag(props, &p); err != nil {
					return err
				}
				body["inlineCommentProperties"] = p
			}
			data, err := confPost(cmd, "/inline-comments", nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	createInlineCmd.Flags().String("page-id", "", "Page ID to comment on")
	createInlineCmd.Flags().String("blogpost-id", "", "Blog post ID to comment on")
	createInlineCmd.Flags().String("parent-comment-id", "", "Parent comment ID (for replies)")
	createInlineCmd.Flags().String("body", "", "Comment body content")
	createInlineCmd.Flags().String("body-format", "storage", "Body format")
	createInlineCmd.Flags().String("inline-comment-properties", "", "JSON inline comment properties")
	confInlineCommentCmd.AddCommand(createInlineCmd)

	// inline update
	updateInlineCmd := &cobra.Command{
		Use:   "update [comment-id]",
		Short: "Update an inline comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
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
			if cmd.Flags().Changed("resolved") {
				body["resolved"] = getBoolFlag(cmd, "resolved")
			}
			data, err := confPut(cmd, "/inline-comments/"+args[0], nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	updateInlineCmd.Flags().Int("version-number", 0, "Version number (required)")
	updateInlineCmd.Flags().String("version-message", "", "Version message")
	updateInlineCmd.Flags().String("body", "", "Comment body content")
	updateInlineCmd.Flags().String("body-format", "storage", "Body format")
	updateInlineCmd.Flags().Bool("resolved", false, "Resolved state")
	_ = updateInlineCmd.MarkFlagRequired("version-number")
	confInlineCommentCmd.AddCommand(updateInlineCmd)

	// inline delete
	deleteInlineCmd := &cobra.Command{
		Use:   "delete [comment-id]",
		Short: "Delete an inline comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/inline-comments/"+args[0], nil)
			if err != nil {
				return err
			}
			fmt.Println("Inline comment deleted successfully.")
			return nil
		},
	}
	confInlineCommentCmd.AddCommand(deleteInlineCmd)

	// inline children
	inlineChildrenCmd := &cobra.Command{
		Use:   "children [comment-id]",
		Short: "Get children inline comments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/inline-comments/"+args[0]+"/children", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(inlineChildrenCmd)
	confInlineCommentCmd.AddCommand(inlineChildrenCmd)

	// inline sub-resources
	for _, sub := range []struct {
		use, short, path string
	}{
		{"operations [comment-id]", "Get permitted operations", "/operations"},
		{"versions [comment-id]", "Get inline comment versions", "/versions"},
		{"likes-count [comment-id]", "Get like count", "/likes/count"},
		{"likes-users [comment-id]", "Get like users", "/likes/users"},
	} {
		sub := sub
		subCmd := &cobra.Command{
			Use:   sub.use,
			Short: sub.short,
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := confGet(cmd, "/inline-comments/"+args[0]+sub.path, nil)
				if err != nil {
					return err
				}
				printJSON(data)
				return nil
			},
		}
		confInlineCommentCmd.AddCommand(subCmd)
	}

	// inline version-details
	inlineVersionDetailCmd := &cobra.Command{
		Use:   "version-details [comment-id] [version-number]",
		Short: "Get version details for inline comment version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/inline-comments/"+args[0]+"/versions/"+args[1], nil)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	confInlineCommentCmd.AddCommand(inlineVersionDetailCmd)
}
