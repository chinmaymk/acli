package acli

import (
	"github.com/spf13/cobra"
)

func init() {
	// label list
	listLabelsCmd := &cobra.Command{
		Use:     "list",
		Short:   "List labels",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if ids := getStringSliceFlag(cmd, "label-id"); len(ids) > 0 {
				for _, id := range ids {
					q.Add("label-id", id)
				}
			}
			if prefixes := getStringSliceFlag(cmd, "prefix"); len(prefixes) > 0 {
				for _, p := range prefixes {
					q.Add("prefix", p)
				}
			}
			data, err := confGet(cmd, "/labels", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(listLabelsCmd)
	addSortFlag(listLabelsCmd)
	listLabelsCmd.Flags().StringSlice("label-id", nil, "Filter by label IDs")
	listLabelsCmd.Flags().StringSlice("prefix", nil, "Filter by prefix")
	confLabelCmd.AddCommand(listLabelsCmd)

	// label pages
	labelPagesCmd := &cobra.Command{
		Use:   "pages [label-id]",
		Short: "Get pages for label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/labels/"+args[0]+"/pages", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(labelPagesCmd)
	addSortFlag(labelPagesCmd)
	addBodyFormatFlag(labelPagesCmd)
	confLabelCmd.AddCommand(labelPagesCmd)

	// label blogposts
	labelBlogPostsCmd := &cobra.Command{
		Use:   "blogposts [label-id]",
		Short: "Get blog posts for label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/labels/"+args[0]+"/blogposts", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(labelBlogPostsCmd)
	addSortFlag(labelBlogPostsCmd)
	addBodyFormatFlag(labelBlogPostsCmd)
	confLabelCmd.AddCommand(labelBlogPostsCmd)

	// label attachments
	labelAttachmentsCmd := &cobra.Command{
		Use:   "attachments [label-id]",
		Short: "Get attachments for label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/labels/"+args[0]+"/attachments", q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	addPaginationFlags(labelAttachmentsCmd)
	addSortFlag(labelAttachmentsCmd)
	confLabelCmd.AddCommand(labelAttachmentsCmd)
}
