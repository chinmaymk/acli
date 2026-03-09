package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// smart-link create
	createSmartLinkCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a smart link in the content tree",
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
			if u := getStringFlag(cmd, "embed-url"); u != "" {
				body["embedUrl"] = u
			}
			data, err := confPost(cmd, "/embeds", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	createSmartLinkCmd.Flags().String("space-id", "", "Space ID (required)")
	createSmartLinkCmd.Flags().String("title", "", "Smart link title")
	createSmartLinkCmd.Flags().String("parent-id", "", "Parent ID")
	createSmartLinkCmd.Flags().String("embed-url", "", "Embed URL")
	_ = createSmartLinkCmd.MarkFlagRequired("space-id")
	confSmartLinkCmd.AddCommand(createSmartLinkCmd)

	// smart-link get
	getSmartLinkCmd := &cobra.Command{
		Use:   "get [smart-link-id]",
		Short: "Get smart link by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			for _, flag := range []string{"include-collaborators", "include-direct-children", "include-operations", "include-properties"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/embeds/"+args[0], q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	getSmartLinkCmd.Flags().Bool("include-collaborators", false, "Include collaborators")
	getSmartLinkCmd.Flags().Bool("include-direct-children", false, "Include direct children")
	getSmartLinkCmd.Flags().Bool("include-operations", false, "Include operations")
	getSmartLinkCmd.Flags().Bool("include-properties", false, "Include properties")
	confSmartLinkCmd.AddCommand(getSmartLinkCmd)

	// smart-link delete
	deleteSmartLinkCmd := &cobra.Command{
		Use:   "delete [smart-link-id]",
		Short: "Delete a smart link",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/embeds/"+args[0], nil)
			if err != nil {
				return err
			}
			fmt.Println("Smart link deleted successfully.")
			return nil
		},
	}
	confSmartLinkCmd.AddCommand(deleteSmartLinkCmd)

	addTreeSubResources(confSmartLinkCmd, "/embeds", "smart link")
}
