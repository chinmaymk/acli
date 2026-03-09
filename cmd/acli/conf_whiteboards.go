package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// whiteboard create
	createWhiteboardCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a whiteboard",
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
			if p := getStringFlag(cmd, "parent-id"); p != "" {
				body["parentId"] = p
			}
			if tk := getStringFlag(cmd, "template-key"); tk != "" {
				body["templateKey"] = tk
			}
			if l := getStringFlag(cmd, "locale"); l != "" {
				body["locale"] = l
			}
			data, err := confPost(cmd, "/whiteboards", q, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	createWhiteboardCmd.Flags().String("space-id", "", "Space ID (required)")
	createWhiteboardCmd.Flags().String("title", "", "Whiteboard title")
	createWhiteboardCmd.Flags().String("parent-id", "", "Parent ID")
	createWhiteboardCmd.Flags().String("template-key", "", "Template key")
	createWhiteboardCmd.Flags().String("locale", "", "Locale")
	createWhiteboardCmd.Flags().Bool("private", false, "Create as private")
	_ = createWhiteboardCmd.MarkFlagRequired("space-id")
	confWhiteboardCmd.AddCommand(createWhiteboardCmd)

	// whiteboard get
	getWhiteboardCmd := &cobra.Command{
		Use:   "get [whiteboard-id]",
		Short: "Get whiteboard by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			for _, flag := range []string{"include-collaborators", "include-direct-children", "include-operations", "include-properties"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/whiteboards/"+args[0], q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	getWhiteboardCmd.Flags().Bool("include-collaborators", false, "Include collaborators")
	getWhiteboardCmd.Flags().Bool("include-direct-children", false, "Include direct children")
	getWhiteboardCmd.Flags().Bool("include-operations", false, "Include operations")
	getWhiteboardCmd.Flags().Bool("include-properties", false, "Include properties")
	confWhiteboardCmd.AddCommand(getWhiteboardCmd)

	// whiteboard delete
	deleteWhiteboardCmd := &cobra.Command{
		Use:   "delete [whiteboard-id]",
		Short: "Delete a whiteboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/whiteboards/"+args[0], nil)
			if err != nil {
				return err
			}
			fmt.Println("Whiteboard deleted successfully.")
			return nil
		},
	}
	confWhiteboardCmd.AddCommand(deleteWhiteboardCmd)

	// whiteboard sub-resources
	for _, sub := range []struct {
		use, short, path string
	}{
		{"ancestors [id]", "Get all ancestors of whiteboard", "/ancestors"},
		{"descendants [id]", "Get descendants of a whiteboard", "/descendants"},
		{"direct-children [id]", "Get direct children of a whiteboard", "/direct-children"},
		{"operations [id]", "Get permitted operations", "/operations"},
		{"properties [id]", "Get content properties", "/properties"},
	} {
		sub := sub
		subCmd := &cobra.Command{
			Use:   sub.use,
			Short: sub.short,
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				q := getPaginationQuery(cmd)
				data, err := confGet(cmd, "/whiteboards/"+args[0]+sub.path, q)
				if err != nil {
					return err
				}
				printJSON(data)
				return nil
			},
		}
		addPaginationFlags(subCmd)
		confWhiteboardCmd.AddCommand(subCmd)
	}
}
