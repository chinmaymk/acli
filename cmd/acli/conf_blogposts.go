package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// blogpost list
	listBlogPostsCmd := &cobra.Command{
		Use:     "list",
		Short:   "List blog posts",
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
			data, err := confGet(cmd, "/blogposts", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(listBlogPostsCmd)
	addSortFlag(listBlogPostsCmd)
	addStatusFlag(listBlogPostsCmd)
	addBodyFormatFlag(listBlogPostsCmd)
	listBlogPostsCmd.Flags().StringSlice("id", nil, "Filter by blog post IDs")
	listBlogPostsCmd.Flags().StringSlice("space-id", nil, "Filter by space IDs")
	listBlogPostsCmd.Flags().String("title", "", "Filter by title")
	confBlogPostCmd.AddCommand(listBlogPostsCmd)

	// blogpost get
	getBlogPostCmd := &cobra.Command{
		Use:   "get [blogpost-id]",
		Short: "Get blog post by ID",
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
				"include-favorited-by-current-user-status", "include-collaborators"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/blogposts/"+args[0], q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addBodyFormatFlag(getBlogPostCmd)
	addStatusFlag(getBlogPostCmd)
	getBlogPostCmd.Flags().Bool("get-draft", false, "Retrieve draft version")
	getBlogPostCmd.Flags().Int("version", 0, "Retrieve a specific version")
	getBlogPostCmd.Flags().Bool("include-labels", false, "Include labels")
	getBlogPostCmd.Flags().Bool("include-properties", false, "Include properties")
	getBlogPostCmd.Flags().Bool("include-operations", false, "Include operations")
	getBlogPostCmd.Flags().Bool("include-likes", false, "Include likes")
	getBlogPostCmd.Flags().Bool("include-versions", false, "Include versions")
	getBlogPostCmd.Flags().Bool("include-version", false, "Include current version")
	getBlogPostCmd.Flags().Bool("include-favorited-by-current-user-status", false, "Include favorited status")
	getBlogPostCmd.Flags().Bool("include-collaborators", false, "Include collaborators")
	confBlogPostCmd.AddCommand(getBlogPostCmd)

	// blogpost create
	createBlogPostCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a blog post",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if getBoolFlag(cmd, "private") {
				q.Set("private", "true")
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
			if b := getStringFlag(cmd, "body"); b != "" {
				body["body"] = map[string]interface{}{
					"representation": getStringFlag(cmd, "body-format"),
					"value":          b,
				}
			}
			data, err := confPost(cmd, "/blogposts", q, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	createBlogPostCmd.Flags().String("space-id", "", "Space ID (required)")
	createBlogPostCmd.Flags().String("title", "", "Blog post title")
	createBlogPostCmd.Flags().String("status", "", "Blog post status")
	createBlogPostCmd.Flags().String("body", "", "Blog post body content")
	createBlogPostCmd.Flags().String("body-format", "storage", "Body format")
	createBlogPostCmd.Flags().Bool("private", false, "Create as private")
	_ = createBlogPostCmd.MarkFlagRequired("space-id")
	confBlogPostCmd.AddCommand(createBlogPostCmd)

	// blogpost update
	updateBlogPostCmd := &cobra.Command{
		Use:   "update [blogpost-id]",
		Short: "Update a blog post",
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
			data, err := confPut(cmd, "/blogposts/"+args[0], nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	updateBlogPostCmd.Flags().String("title", "", "Blog post title (required)")
	updateBlogPostCmd.Flags().String("status", "current", "Blog post status (required)")
	updateBlogPostCmd.Flags().String("body", "", "Blog post body content")
	updateBlogPostCmd.Flags().String("body-format", "storage", "Body format")
	updateBlogPostCmd.Flags().Int("version-number", 0, "Version number (required)")
	updateBlogPostCmd.Flags().String("version-message", "", "Version message")
	updateBlogPostCmd.Flags().String("space-id", "", "Space ID")
	_ = updateBlogPostCmd.MarkFlagRequired("title")
	_ = updateBlogPostCmd.MarkFlagRequired("version-number")
	confBlogPostCmd.AddCommand(updateBlogPostCmd)

	// blogpost delete
	deleteBlogPostCmd := &cobra.Command{
		Use:   "delete [blogpost-id]",
		Short: "Delete a blog post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if getBoolFlag(cmd, "purge") {
				q.Set("purge", "true")
			}
			if getBoolFlag(cmd, "draft") {
				q.Set("draft", "true")
			}
			_, err := confDelete(cmd, "/blogposts/"+args[0], q)
			if err != nil {
				return err
			}
			fmt.Println("Blog post deleted successfully.")
			return nil
		},
	}
	deleteBlogPostCmd.Flags().Bool("purge", false, "Purge the blog post")
	deleteBlogPostCmd.Flags().Bool("draft", false, "Delete a draft blog post")
	confBlogPostCmd.AddCommand(deleteBlogPostCmd)

	// blogpost sub-resources
	for _, sub := range []struct {
		use, short, path string
	}{
		{"attachments [blogpost-id]", "Get attachments for blog post", "/attachments"},
		{"labels [blogpost-id]", "Get labels for blog post", "/labels"},
		{"footer-comments [blogpost-id]", "Get footer comments for blog post", "/footer-comments"},
		{"inline-comments [blogpost-id]", "Get inline comments for blog post", "/inline-comments"},
		{"custom-content [blogpost-id]", "Get custom content in blog post", "/custom-content"},
		{"operations [blogpost-id]", "Get permitted operations for blog post", "/operations"},
		{"versions [blogpost-id]", "Get blog post versions", "/versions"},
	} {
		sub := sub
		subCmd := &cobra.Command{
			Use:   sub.use,
			Short: sub.short,
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				q := getPaginationQuery(cmd)
				data, err := confGet(cmd, "/blogposts/"+args[0]+sub.path, q)
				if err != nil {
					return err
				}
				printJSON(data)
				return nil
			},
		}
		addPaginationFlags(subCmd)
		addSortFlag(subCmd)
		confBlogPostCmd.AddCommand(subCmd)
	}

	// blogpost version-details
	bpVersionDetailCmd := &cobra.Command{
		Use:   "version-details [blogpost-id] [version-number]",
		Short: "Get version details for blog post version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/blogposts/"+args[0]+"/versions/"+args[1], nil)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	confBlogPostCmd.AddCommand(bpVersionDetailCmd)

	// blogpost likes
	bpLikesCountCmd := &cobra.Command{
		Use:   "likes-count [blogpost-id]",
		Short: "Get like count for blog post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/blogposts/"+args[0]+"/likes/count", nil)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	confBlogPostCmd.AddCommand(bpLikesCountCmd)

	bpLikesUsersCmd := &cobra.Command{
		Use:   "likes-users [blogpost-id]",
		Short: "Get account IDs of likes for blog post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/blogposts/"+args[0]+"/likes/users", nil)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	confBlogPostCmd.AddCommand(bpLikesUsersCmd)

	// blogpost redact
	bpRedactCmd := &cobra.Command{
		Use:   "redact [blogpost-id]",
		Short: "Redact content in a Confluence blog post",
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
			data, err := confPost(cmd, "/blogposts/"+args[0]+"/redact", nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	bpRedactCmd.Flags().String("body", "", "JSON redaction request body")
	confBlogPostCmd.AddCommand(bpRedactCmd)
}
