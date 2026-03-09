package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// Content properties are managed per-resource. This command provides a unified
	// interface for managing properties across pages, blogposts, comments, attachments,
	// custom content, whiteboards, databases, folders, smart links, and spaces.

	for _, res := range []struct {
		name, pathPrefix, idParam string
	}{
		{"page", "/pages", "page-id"},
		{"blogpost", "/blogposts", "blogpost-id"},
		{"comment", "/comments", "comment-id"},
		{"attachment", "/attachments", "attachment-id"},
		{"custom-content", "/custom-content", "custom-content-id"},
		{"whiteboard", "/whiteboards", "whiteboard-id"},
		{"database", "/databases", "database-id"},
		{"folder", "/folders", "folder-id"},
		{"smart-link", "/embeds", "embed-id"},
	} {
		res := res
		resCmd := &cobra.Command{
			Use:   res.name,
			Short: fmt.Sprintf("Manage content properties for %s", res.name),
			RunE:  helpRunE,
		}

		// list properties
		listCmd := &cobra.Command{
			Use:     "list [" + res.idParam + "]",
			Short:   fmt.Sprintf("List content properties for %s", res.name),
			Aliases: []string{"ls"},
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				q := getPaginationQuery(cmd)
				if k := getStringFlag(cmd, "key"); k != "" {
					q.Set("key", k)
				}
				data, err := confGet(cmd, res.pathPrefix+"/"+args[0]+"/properties", q)
				if err != nil {
					return err
				}
				printJSONBytes(data)
				return nil
			},
		}
		addPaginationFlags(listCmd)
		addSortFlag(listCmd)
		listCmd.Flags().String("key", "", "Filter by property key")
		resCmd.AddCommand(listCmd)

		// get property
		getCmd := &cobra.Command{
			Use:   "get [" + res.idParam + "] [property-id]",
			Short: fmt.Sprintf("Get content property for %s by ID", res.name),
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := confGet(cmd, res.pathPrefix+"/"+args[0]+"/properties/"+args[1], nil)
				if err != nil {
					return err
				}
				printJSONBytes(data)
				return nil
			},
		}
		resCmd.AddCommand(getCmd)

		// create property
		createCmd := &cobra.Command{
			Use:   "create [" + res.idParam + "]",
			Short: fmt.Sprintf("Create content property for %s", res.name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				body := map[string]interface{}{
					"key": getStringFlag(cmd, "key"),
				}
				valStr := getStringFlag(cmd, "value")
				if valStr != "" {
					var val interface{}
					if err := parseJSONFlag(valStr, &val); err != nil {
						// Treat as string value if not valid JSON
						body["value"] = valStr
					} else {
						body["value"] = val
					}
				}
				data, err := confPost(cmd, res.pathPrefix+"/"+args[0]+"/properties", nil, body)
				if err != nil {
					return err
				}
				printJSONBytes(data)
				return nil
			},
		}
		createCmd.Flags().String("key", "", "Property key (required)")
		createCmd.Flags().String("value", "", "Property value (JSON)")
		_ = createCmd.MarkFlagRequired("key")
		resCmd.AddCommand(createCmd)

		// update property
		updateCmd := &cobra.Command{
			Use:   "update [" + res.idParam + "] [property-id]",
			Short: fmt.Sprintf("Update content property for %s", res.name),
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				body := map[string]interface{}{
					"key": getStringFlag(cmd, "key"),
					"version": map[string]interface{}{
						"number":  getIntFlag(cmd, "version-number"),
						"message": getStringFlag(cmd, "version-message"),
					},
				}
				valStr := getStringFlag(cmd, "value")
				if valStr != "" {
					var val interface{}
					if err := parseJSONFlag(valStr, &val); err != nil {
						body["value"] = valStr
					} else {
						body["value"] = val
					}
				}
				data, err := confPut(cmd, res.pathPrefix+"/"+args[0]+"/properties/"+args[1], nil, body)
				if err != nil {
					return err
				}
				printJSONBytes(data)
				return nil
			},
		}
		updateCmd.Flags().String("key", "", "Property key (required)")
		updateCmd.Flags().String("value", "", "Property value (JSON)")
		updateCmd.Flags().Int("version-number", 0, "Version number (required)")
		updateCmd.Flags().String("version-message", "", "Version message")
		_ = updateCmd.MarkFlagRequired("key")
		_ = updateCmd.MarkFlagRequired("version-number")
		resCmd.AddCommand(updateCmd)

		// delete property
		deleteCmd := &cobra.Command{
			Use:   "delete [" + res.idParam + "] [property-id]",
			Short: fmt.Sprintf("Delete content property for %s", res.name),
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				_, err := confDelete(cmd, res.pathPrefix+"/"+args[0]+"/properties/"+args[1], nil)
				if err != nil {
					return err
				}
				fmt.Println("Property deleted successfully.")
				return nil
			},
		}
		resCmd.AddCommand(deleteCmd)

		confPropertyCmd.AddCommand(resCmd)
	}

	// Space properties (different path pattern: /spaces/{space-id}/properties)
	spacePropsCmd := &cobra.Command{
		Use:   "space",
		Short: "Manage space properties",
		RunE:  helpRunE,
	}

	spListCmd := &cobra.Command{
		Use:     "list [space-id]",
		Short:   "List space properties",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if k := getStringFlag(cmd, "key"); k != "" {
				q.Set("key", k)
			}
			data, err := confGet(cmd, "/spaces/"+args[0]+"/properties", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(spListCmd)
	addSortFlag(spListCmd)
	spListCmd.Flags().String("key", "", "Filter by property key")
	spacePropsCmd.AddCommand(spListCmd)

	spGetCmd := &cobra.Command{
		Use:   "get [space-id] [property-id]",
		Short: "Get space property by ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/spaces/"+args[0]+"/properties/"+args[1], nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	spacePropsCmd.AddCommand(spGetCmd)

	spCreateCmd := &cobra.Command{
		Use:   "create [space-id]",
		Short: "Create space property",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"key": getStringFlag(cmd, "key"),
			}
			valStr := getStringFlag(cmd, "value")
			if valStr != "" {
				var val interface{}
				if err := parseJSONFlag(valStr, &val); err != nil {
					body["value"] = valStr
				} else {
					body["value"] = val
				}
			}
			data, err := confPost(cmd, "/spaces/"+args[0]+"/properties", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	spCreateCmd.Flags().String("key", "", "Property key (required)")
	spCreateCmd.Flags().String("value", "", "Property value (JSON)")
	_ = spCreateCmd.MarkFlagRequired("key")
	spacePropsCmd.AddCommand(spCreateCmd)

	spUpdateCmd := &cobra.Command{
		Use:   "update [space-id] [property-id]",
		Short: "Update space property",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"key": getStringFlag(cmd, "key"),
				"version": map[string]interface{}{
					"number":  getIntFlag(cmd, "version-number"),
					"message": getStringFlag(cmd, "version-message"),
				},
			}
			valStr := getStringFlag(cmd, "value")
			if valStr != "" {
				var val interface{}
				if err := parseJSONFlag(valStr, &val); err != nil {
					body["value"] = valStr
				} else {
					body["value"] = val
				}
			}
			data, err := confPut(cmd, "/spaces/"+args[0]+"/properties/"+args[1], nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	spUpdateCmd.Flags().String("key", "", "Property key (required)")
	spUpdateCmd.Flags().String("value", "", "Property value (JSON)")
	spUpdateCmd.Flags().Int("version-number", 0, "Version number (required)")
	spUpdateCmd.Flags().String("version-message", "", "Version message")
	_ = spUpdateCmd.MarkFlagRequired("key")
	_ = spUpdateCmd.MarkFlagRequired("version-number")
	spacePropsCmd.AddCommand(spUpdateCmd)

	spDeleteCmd := &cobra.Command{
		Use:   "delete [space-id] [property-id]",
		Short: "Delete space property",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/spaces/"+args[0]+"/properties/"+args[1], nil)
			if err != nil {
				return err
			}
			fmt.Println("Space property deleted successfully.")
			return nil
		},
	}
	spacePropsCmd.AddCommand(spDeleteCmd)

	confPropertyCmd.AddCommand(spacePropsCmd)
}
