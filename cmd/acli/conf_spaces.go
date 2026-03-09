package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// space list
	listSpacesCmd := &cobra.Command{
		Use:     "list",
		Short:   "List spaces",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if ids := getStringSliceFlag(cmd, "ids"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("ids", id)
				}
			}
			if keys := getStringSliceFlag(cmd, "keys"); len(keys) > 0 {
				for _, k := range keys {
					q.Add("keys", k)
				}
			}
			if t := getStringFlag(cmd, "type"); t != "" {
				q.Set("type", t)
			}
			if l := getStringSliceFlag(cmd, "labels"); len(l) > 0 {
				for _, label := range l {
					q.Add("labels", label)
				}
			}
			if f := getStringFlag(cmd, "favorited-by"); f != "" {
				q.Set("favorited-by", f)
			}
			if f := getStringFlag(cmd, "not-favorited-by"); f != "" {
				q.Set("not-favorited-by", f)
			}
			if f := getStringFlag(cmd, "description-format"); f != "" {
				q.Set("description-format", f)
			}
			if getBoolFlag(cmd, "include-icon") {
				q.Set("include-icon", "true")
			}
			data, err := confGet(cmd, "/spaces", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(listSpacesCmd)
	addSortFlag(listSpacesCmd)
	addStatusFlag(listSpacesCmd)
	listSpacesCmd.Flags().StringSlice("ids", nil, "Filter by space IDs")
	listSpacesCmd.Flags().StringSlice("keys", nil, "Filter by space keys")
	listSpacesCmd.Flags().String("type", "", "Filter by type (global, personal)")
	listSpacesCmd.Flags().StringSlice("labels", nil, "Filter by labels")
	listSpacesCmd.Flags().String("favorited-by", "", "Filter by favorited-by user account ID")
	listSpacesCmd.Flags().String("not-favorited-by", "", "Filter by not-favorited-by user account ID")
	listSpacesCmd.Flags().String("description-format", "", "Description format (plain, view)")
	listSpacesCmd.Flags().Bool("include-icon", false, "Include space icon")
	confSpaceCmd.AddCommand(listSpacesCmd)

	// space get
	getSpaceCmd := &cobra.Command{
		Use:   "get [space-id]",
		Short: "Get space by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if f := getStringFlag(cmd, "description-format"); f != "" {
				q.Set("description-format", f)
			}
			if getBoolFlag(cmd, "include-icon") {
				q.Set("include-icon", "true")
			}
			if getBoolFlag(cmd, "include-operations") {
				q.Set("include-operations", "true")
			}
			if getBoolFlag(cmd, "include-properties") {
				q.Set("include-properties", "true")
			}
			if getBoolFlag(cmd, "include-permissions") {
				q.Set("include-permissions", "true")
			}
			if getBoolFlag(cmd, "include-role-assignments") {
				q.Set("include-role-assignments", "true")
			}
			if getBoolFlag(cmd, "include-labels") {
				q.Set("include-labels", "true")
			}
			data, err := confGet(cmd, "/spaces/"+args[0], q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	getSpaceCmd.Flags().String("description-format", "", "Description format")
	getSpaceCmd.Flags().Bool("include-icon", false, "Include icon")
	getSpaceCmd.Flags().Bool("include-operations", false, "Include operations")
	getSpaceCmd.Flags().Bool("include-properties", false, "Include properties")
	getSpaceCmd.Flags().Bool("include-permissions", false, "Include permissions")
	getSpaceCmd.Flags().Bool("include-role-assignments", false, "Include role assignments")
	getSpaceCmd.Flags().Bool("include-labels", false, "Include labels")
	confSpaceCmd.AddCommand(getSpaceCmd)

	// space create
	createSpaceCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a space",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"name": getStringFlag(cmd, "name"),
			}
			if k := getStringFlag(cmd, "key"); k != "" {
				body["key"] = k
			}
			if a := getStringFlag(cmd, "alias"); a != "" {
				body["alias"] = a
			}
			if desc := getStringFlag(cmd, "description"); desc != "" {
				body["description"] = map[string]interface{}{
					"representation": "plain",
					"value":          desc,
				}
			}
			if getBoolFlag(cmd, "private") {
				body["createPrivateSpace"] = true
			}
			if tmpl := getStringFlag(cmd, "template-key"); tmpl != "" {
				body["templateKey"] = tmpl
			}
			data, err := confPost(cmd, "/spaces", nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	createSpaceCmd.Flags().String("name", "", "Space name (required)")
	createSpaceCmd.Flags().String("key", "", "Space key")
	createSpaceCmd.Flags().String("alias", "", "Space alias")
	createSpaceCmd.Flags().String("description", "", "Space description")
	createSpaceCmd.Flags().Bool("private", false, "Create as private space")
	createSpaceCmd.Flags().String("template-key", "", "Template key")
	_ = createSpaceCmd.MarkFlagRequired("name")
	confSpaceCmd.AddCommand(createSpaceCmd)

	// space pages
	spacePagesCmd := &cobra.Command{
		Use:   "pages [space-id]",
		Short: "List pages in a space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if d := getStringFlag(cmd, "depth"); d != "" {
				q.Set("depth", d)
			}
			if t := getStringFlag(cmd, "title"); t != "" {
				q.Set("title", t)
			}
			data, err := confGet(cmd, "/spaces/"+args[0]+"/pages", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(spacePagesCmd)
	addSortFlag(spacePagesCmd)
	addStatusFlag(spacePagesCmd)
	addBodyFormatFlag(spacePagesCmd)
	spacePagesCmd.Flags().String("depth", "", "Filter by depth (root, all)")
	spacePagesCmd.Flags().String("title", "", "Filter by title")
	confSpaceCmd.AddCommand(spacePagesCmd)

	// space blogposts
	spaceBlogPostsCmd := &cobra.Command{
		Use:   "blogposts [space-id]",
		Short: "List blog posts in a space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if t := getStringFlag(cmd, "title"); t != "" {
				q.Set("title", t)
			}
			data, err := confGet(cmd, "/spaces/"+args[0]+"/blogposts", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(spaceBlogPostsCmd)
	addSortFlag(spaceBlogPostsCmd)
	addStatusFlag(spaceBlogPostsCmd)
	addBodyFormatFlag(spaceBlogPostsCmd)
	spaceBlogPostsCmd.Flags().String("title", "", "Filter by title")
	confSpaceCmd.AddCommand(spaceBlogPostsCmd)

	// space labels
	spaceLabelsCmd := &cobra.Command{
		Use:   "labels [space-id]",
		Short: "Get labels for a space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/spaces/"+args[0]+"/labels", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(spaceLabelsCmd)
	confSpaceCmd.AddCommand(spaceLabelsCmd)

	// space content-labels
	spaceContentLabelsCmd := &cobra.Command{
		Use:   "content-labels [space-id]",
		Short: "Get labels for space content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if p := getStringFlag(cmd, "prefix"); p != "" {
				q.Set("prefix", p)
			}
			data, err := confGet(cmd, "/spaces/"+args[0]+"/content/labels", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(spaceContentLabelsCmd)
	spaceContentLabelsCmd.Flags().String("prefix", "", "Filter by prefix")
	confSpaceCmd.AddCommand(spaceContentLabelsCmd)

	// space custom-content
	spaceCustomContentCmd := &cobra.Command{
		Use:   "custom-content [space-id]",
		Short: "Get custom content by type in space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if t := getStringFlag(cmd, "type"); t != "" {
				q.Set("type", t)
			}
			data, err := confGet(cmd, "/spaces/"+args[0]+"/custom-content", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(spaceCustomContentCmd)
	addSortFlag(spaceCustomContentCmd)
	addBodyFormatFlag(spaceCustomContentCmd)
	spaceCustomContentCmd.Flags().String("type", "", "Custom content type")
	confSpaceCmd.AddCommand(spaceCustomContentCmd)

	// space operations
	spaceOpsCmd := &cobra.Command{
		Use:   "operations [space-id]",
		Short: "Get permitted operations for space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/spaces/"+args[0]+"/operations", nil)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	confSpaceCmd.AddCommand(spaceOpsCmd)

	// space permissions
	spacePermsCmd := &cobra.Command{
		Use:   "permissions [space-id]",
		Short: "Get space permissions assignments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if cursor := getStringFlag(cmd, "cursor"); cursor != "" {
				q.Set("cursor", cursor)
			}
			if limit := getIntFlag(cmd, "limit"); limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			data, err := confGet(cmd, "/spaces/"+args[0]+"/permissions", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(spacePermsCmd)
	confSpaceCmd.AddCommand(spacePermsCmd)

	// space role-assignments
	spaceRoleAssignCmd := &cobra.Command{
		Use:   "role-assignments [space-id]",
		Short: "Get space role assignments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if cursor := getStringFlag(cmd, "cursor"); cursor != "" {
				q.Set("cursor", cursor)
			}
			if limit := getIntFlag(cmd, "limit"); limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			data, err := confGet(cmd, "/spaces/"+args[0]+"/role-assignments", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(spaceRoleAssignCmd)
	confSpaceCmd.AddCommand(spaceRoleAssignCmd)

	// space set-role-assignments
	setSpaceRoleAssignCmd := &cobra.Command{
		Use:   "set-role-assignments [space-id]",
		Short: "Set space role assignments (provide JSON body via stdin or --body flag)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bodyStr := getStringFlag(cmd, "body")
			if bodyStr == "" {
				return fmt.Errorf("--body is required (JSON role assignments)")
			}
			var body interface{}
			if err := parseJSONFlag(bodyStr, &body); err != nil {
				return err
			}
			data, err := confPost(cmd, "/spaces/"+args[0]+"/role-assignments", nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	setSpaceRoleAssignCmd.Flags().String("body", "", "JSON body for role assignments")
	confSpaceCmd.AddCommand(setSpaceRoleAssignCmd)
}
