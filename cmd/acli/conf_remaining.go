package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// === Classification Levels ===
	listClassificationCmd := &cobra.Command{
		Use:     "list",
		Short:   "Get list of classification levels",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/classification-levels", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(listClassificationCmd)
	confClassificationCmd.AddCommand(listClassificationCmd)

	// Classification for pages
	for _, res := range []struct {
		name, pathPrefix string
	}{
		{"page", "/pages"},
		{"blogpost", "/blogposts"},
		{"database", "/databases"},
		{"whiteboard", "/whiteboards"},
	} {
		res := res
		resCmd := &cobra.Command{
			Use:   res.name,
			Short: fmt.Sprintf("Manage classification level for %s", res.name),
			RunE:  helpRunE,
		}

		getCmd := &cobra.Command{
			Use:   "get [id]",
			Short: fmt.Sprintf("Get %s classification level", res.name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := confGet(cmd, res.pathPrefix+"/"+args[0]+"/classification-level", nil)
				if err != nil {
					return err
				}
				printJSONBytes(data)
				return nil
			},
		}
		resCmd.AddCommand(getCmd)

		setCmd := &cobra.Command{
			Use:   "set [id]",
			Short: fmt.Sprintf("Update %s classification level", res.name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				body := map[string]interface{}{
					"id":     getStringFlag(cmd, "classification-id"),
					"status": getStringFlag(cmd, "status"),
				}
				data, err := confPut(cmd, res.pathPrefix+"/"+args[0]+"/classification-level", nil, body)
				if err != nil {
					return err
				}
				printJSONBytes(data)
				return nil
			},
		}
		setCmd.Flags().String("classification-id", "", "Classification level ID (required)")
		setCmd.Flags().String("status", "", "Status (required)")
		_ = setCmd.MarkFlagRequired("classification-id")
		_ = setCmd.MarkFlagRequired("status")
		resCmd.AddCommand(setCmd)

		resetCmd := &cobra.Command{
			Use:   "reset [id]",
			Short: fmt.Sprintf("Reset %s classification level", res.name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				body := map[string]interface{}{
					"status": getStringFlag(cmd, "status"),
				}
				data, err := confPost(cmd, res.pathPrefix+"/"+args[0]+"/classification-level/reset", nil, body)
				if err != nil {
					return err
				}
				printJSONBytes(data)
				return nil
			},
		}
		resetCmd.Flags().String("status", "", "Status (required)")
		_ = resetCmd.MarkFlagRequired("status")
		resCmd.AddCommand(resetCmd)

		confClassificationCmd.AddCommand(resCmd)
	}

	// Space default classification level
	spaceClassCmd := &cobra.Command{
		Use:   "space",
		Short: "Manage space default classification level",
		RunE:  helpRunE,
	}

	spaceClassGetCmd := &cobra.Command{
		Use:   "get [space-id]",
		Short: "Get space default classification level",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/spaces/"+args[0]+"/classification-level/default", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	spaceClassCmd.AddCommand(spaceClassGetCmd)

	spaceClassSetCmd := &cobra.Command{
		Use:   "set [space-id]",
		Short: "Update space default classification level",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"id": getStringFlag(cmd, "classification-id"),
			}
			data, err := confPut(cmd, "/spaces/"+args[0]+"/classification-level/default", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	spaceClassSetCmd.Flags().String("classification-id", "", "Classification level ID (required)")
	_ = spaceClassSetCmd.MarkFlagRequired("classification-id")
	spaceClassCmd.AddCommand(spaceClassSetCmd)

	spaceClassDeleteCmd := &cobra.Command{
		Use:   "delete [space-id]",
		Short: "Delete space default classification level",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/spaces/"+args[0]+"/classification-level/default", nil)
			if err != nil {
				return err
			}
			fmt.Println("Space default classification level deleted successfully.")
			return nil
		},
	}
	spaceClassCmd.AddCommand(spaceClassDeleteCmd)

	confClassificationCmd.AddCommand(spaceClassCmd)

	// === Data Policies ===
	dpMetadataCmd := &cobra.Command{
		Use:   "metadata",
		Short: "Get data policy metadata for the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/data-policies/metadata", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confDataPolicyCmd.AddCommand(dpMetadataCmd)

	dpSpacesCmd := &cobra.Command{
		Use:   "spaces",
		Short: "Get spaces with data policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if cursor := getStringFlag(cmd, "cursor"); cursor != "" {
				q.Set("cursor", cursor)
			}
			if limit := getIntFlag(cmd, "limit"); limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			if sort := getStringFlag(cmd, "sort"); sort != "" {
				q.Set("sort", sort)
			}
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
			data, err := confGet(cmd, "/data-policies/spaces", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(dpSpacesCmd)
	addSortFlag(dpSpacesCmd)
	dpSpacesCmd.Flags().StringSlice("ids", nil, "Filter by space IDs")
	dpSpacesCmd.Flags().StringSlice("keys", nil, "Filter by space keys")
	confDataPolicyCmd.AddCommand(dpSpacesCmd)

	// === Space Permissions ===
	getAvailablePermsCmd := &cobra.Command{
		Use:   "available",
		Short: "Get available space permissions",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/space-permissions", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confSpacePermissionCmd.AddCommand(getAvailablePermsCmd)

	// === Space Roles ===
	listSpaceRolesCmd := &cobra.Command{
		Use:     "list",
		Short:   "Get available space roles",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if s := getStringFlag(cmd, "space-id"); s != "" {
				q.Set("space-id", s)
			}
			if r := getStringFlag(cmd, "role-type"); r != "" {
				q.Set("role-type", r)
			}
			if p := getStringFlag(cmd, "principal-id"); p != "" {
				q.Set("principal-id", p)
			}
			if pt := getStringFlag(cmd, "principal-type"); pt != "" {
				q.Set("principal-type", pt)
			}
			if cursor := getStringFlag(cmd, "cursor"); cursor != "" {
				q.Set("cursor", cursor)
			}
			if limit := getIntFlag(cmd, "limit"); limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			data, err := confGet(cmd, "/space-roles", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(listSpaceRolesCmd)
	listSpaceRolesCmd.Flags().String("space-id", "", "Filter by space ID")
	listSpaceRolesCmd.Flags().String("role-type", "", "Filter by role type")
	listSpaceRolesCmd.Flags().String("principal-id", "", "Filter by principal ID")
	listSpaceRolesCmd.Flags().String("principal-type", "", "Filter by principal type")
	confSpaceRoleCmd.AddCommand(listSpaceRolesCmd)

	getSpaceRoleCmd := &cobra.Command{
		Use:   "get [role-id]",
		Short: "Get space role by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/space-roles/"+args[0], nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confSpaceRoleCmd.AddCommand(getSpaceRoleCmd)

	createSpaceRoleCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a space role",
		RunE: func(cmd *cobra.Command, args []string) error {
			bodyStr := getStringFlag(cmd, "body")
			if bodyStr == "" {
				return fmt.Errorf("--body is required (JSON space role definition)")
			}
			var body interface{}
			if err := parseJSONFlag(bodyStr, &body); err != nil {
				return err
			}
			data, err := confPost(cmd, "/space-roles", nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	createSpaceRoleCmd.Flags().String("body", "", "JSON space role definition")
	confSpaceRoleCmd.AddCommand(createSpaceRoleCmd)

	updateSpaceRoleCmd := &cobra.Command{
		Use:   "update [role-id]",
		Short: "Update a space role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bodyStr := getStringFlag(cmd, "body")
			if bodyStr == "" {
				return fmt.Errorf("--body is required (JSON space role definition)")
			}
			var body interface{}
			if err := parseJSONFlag(bodyStr, &body); err != nil {
				return err
			}
			data, err := confPut(cmd, "/space-roles/"+args[0], nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	updateSpaceRoleCmd.Flags().String("body", "", "JSON space role definition")
	confSpaceRoleCmd.AddCommand(updateSpaceRoleCmd)

	deleteSpaceRoleCmd := &cobra.Command{
		Use:   "delete [role-id]",
		Short: "Delete a space role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/space-roles/"+args[0], nil)
			if err != nil {
				return err
			}
			fmt.Println("Space role deleted successfully.")
			return nil
		},
	}
	confSpaceRoleCmd.AddCommand(deleteSpaceRoleCmd)

	// Space role mode
	spaceRoleModeCmd := &cobra.Command{
		Use:   "mode",
		Short: "Get space role mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/space-role-mode", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confSpaceRoleCmd.AddCommand(spaceRoleModeCmd)

	// === Forge App Properties ===
	appPropsCmd := &cobra.Command{
		Use:   "app-property",
		Short: "Manage Forge app properties",
		Aliases: []string{"ap"},
		RunE:  helpRunE,
	}

	appPropsListCmd := &cobra.Command{
		Use:     "list",
		Short:   "Get Forge app properties",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/app/properties", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	appPropsCmd.AddCommand(appPropsListCmd)

	appPropsGetCmd := &cobra.Command{
		Use:   "get [property-key]",
		Short: "Get a Forge app property by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/app/properties/"+args[0], nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	appPropsCmd.AddCommand(appPropsGetCmd)

	appPropsSetCmd := &cobra.Command{
		Use:   "set [property-key]",
		Short: "Create or update a Forge app property",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bodyStr := getStringFlag(cmd, "body")
			if bodyStr == "" {
				return fmt.Errorf("--body is required (JSON property value)")
			}
			var body interface{}
			if err := parseJSONFlag(bodyStr, &body); err != nil {
				return err
			}
			data, err := confPut(cmd, "/app/properties/"+args[0], nil, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	appPropsSetCmd.Flags().String("body", "", "JSON property value")
	appPropsCmd.AddCommand(appPropsSetCmd)

	appPropsDeleteCmd := &cobra.Command{
		Use:   "delete [property-key]",
		Short: "Delete a Forge app property",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/app/properties/"+args[0], nil)
			if err != nil {
				return err
			}
			fmt.Println("App property deleted successfully.")
			return nil
		},
	}
	appPropsCmd.AddCommand(appPropsDeleteCmd)

	confluenceCmd.AddCommand(appPropsCmd)
}
