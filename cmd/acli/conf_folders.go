package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// folder create
	createFolderCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"spaceId": getStringFlag(cmd, "space-id"),
			}
			if t := getStringFlag(cmd, "title"); t != "" {
				body["title"] = t
			}
			if p := getStringFlag(cmd, "parent-id"); p != "" {
				body["parentId"] = p
			}
			data, err := confPost(cmd, "/folders", nil, body)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	createFolderCmd.Flags().String("space-id", "", "Space ID (required)")
	createFolderCmd.Flags().String("title", "", "Folder title")
	createFolderCmd.Flags().String("parent-id", "", "Parent ID")
	_ = createFolderCmd.MarkFlagRequired("space-id")
	confFolderCmd.AddCommand(createFolderCmd)

	// folder get
	getFolderCmd := &cobra.Command{
		Use:   "get [folder-id]",
		Short: "Get folder by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			for _, flag := range []string{"include-collaborators", "include-direct-children", "include-operations", "include-properties"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/folders/"+args[0], q)
			if err != nil {
				return err
			}
			printJSON(data)
			return nil
		},
	}
	getFolderCmd.Flags().Bool("include-collaborators", false, "Include collaborators")
	getFolderCmd.Flags().Bool("include-direct-children", false, "Include direct children")
	getFolderCmd.Flags().Bool("include-operations", false, "Include operations")
	getFolderCmd.Flags().Bool("include-properties", false, "Include properties")
	confFolderCmd.AddCommand(getFolderCmd)

	// folder delete
	deleteFolderCmd := &cobra.Command{
		Use:   "delete [folder-id]",
		Short: "Delete a folder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/folders/"+args[0], nil)
			if err != nil {
				return err
			}
			fmt.Println("Folder deleted successfully.")
			return nil
		},
	}
	confFolderCmd.AddCommand(deleteFolderCmd)

	// folder sub-resources
	for _, sub := range []struct {
		use, short, path string
	}{
		{"ancestors [id]", "Get all ancestors of folder", "/ancestors"},
		{"descendants [id]", "Get descendants of folder", "/descendants"},
		{"direct-children [id]", "Get direct children of a folder", "/direct-children"},
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
				data, err := confGet(cmd, "/folders/"+args[0]+sub.path, q)
				if err != nil {
					return err
				}
				printJSON(data)
				return nil
			},
		}
		addPaginationFlags(subCmd)
		confFolderCmd.AddCommand(subCmd)
	}
}
