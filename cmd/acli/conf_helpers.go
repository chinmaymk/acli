package acli

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chinmaymk/acli/internal/api"
	"github.com/chinmaymk/acli/internal/config"
	"github.com/spf13/cobra"
)

func newConfluenceClient(cmd *cobra.Command) (*api.Client, error) {
	profileName, _ := cmd.Flags().GetString("profile")
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return nil, err
	}

	if profile.AtlassianURL == "" {
		return nil, fmt.Errorf("no Atlassian URL configured: run 'acli config setup' to set one")
	}
	if profile.APIToken == "" {
		return nil, fmt.Errorf("no API token configured: run 'acli config setup' to set one")
	}

	return api.NewClient(profile.AtlassianURL, profile.Email, profile.APIToken), nil
}

func confGet(cmd *cobra.Command, path string, query url.Values) ([]byte, error) {
	client, err := newConfluenceClient(cmd)
	if err != nil {
		return nil, err
	}
	return client.ConfluenceV2("GET", path, query, nil)
}

func confPost(cmd *cobra.Command, path string, query url.Values, body interface{}) ([]byte, error) {
	client, err := newConfluenceClient(cmd)
	if err != nil {
		return nil, err
	}
	return client.ConfluenceV2("POST", path, query, body)
}

func confPut(cmd *cobra.Command, path string, query url.Values, body interface{}) ([]byte, error) {
	client, err := newConfluenceClient(cmd)
	if err != nil {
		return nil, err
	}
	return client.ConfluenceV2("PUT", path, query, body)
}

func confDelete(cmd *cobra.Command, path string, query url.Values) ([]byte, error) {
	client, err := newConfluenceClient(cmd)
	if err != nil {
		return nil, err
	}
	return client.ConfluenceV2("DELETE", path, query, nil)
}

func printJSONBytes(data []byte) {
	var out interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		fmt.Println(string(data))
		return
	}
	pretty, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		fmt.Println(string(data))
		return
	}
	fmt.Println(string(pretty))
}

func addPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().Int("limit", 25, "Maximum number of results to return")
	cmd.Flags().String("cursor", "", "Pagination cursor")
}

func addSortFlag(cmd *cobra.Command) {
	cmd.Flags().String("sort", "", "Sort field")
}

func addBodyFormatFlag(cmd *cobra.Command) {
	cmd.Flags().String("body-format", "", "Body format (storage, atlas_doc_format, view, export_view, anonymous_export_view, styled_view, editor)")
}

func addStatusFlag(cmd *cobra.Command) {
	cmd.Flags().StringSlice("status", nil, "Filter by status")
}

func getPaginationQuery(cmd *cobra.Command) url.Values {
	q := url.Values{}
	if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor, _ := cmd.Flags().GetString("cursor"); cursor != "" {
		q.Set("cursor", cursor)
	}
	if sort, _ := cmd.Flags().GetString("sort"); sort != "" {
		q.Set("sort", sort)
	}
	if bodyFormat, _ := cmd.Flags().GetString("body-format"); bodyFormat != "" {
		q.Set("body-format", bodyFormat)
	}
	if statuses, _ := cmd.Flags().GetStringSlice("status"); len(statuses) > 0 {
		for _, s := range statuses {
			q.Add("status", s)
		}
	}
	return q
}

func getStringFlag(cmd *cobra.Command, name string) string {
	val, _ := cmd.Flags().GetString(name)
	return val
}

func getBoolFlag(cmd *cobra.Command, name string) bool {
	val, _ := cmd.Flags().GetBool(name)
	return val
}

func getIntFlag(cmd *cobra.Command, name string) int {
	val, _ := cmd.Flags().GetInt(name)
	return val
}

func getStringSliceFlag(cmd *cobra.Command, name string) []string {
	val, _ := cmd.Flags().GetStringSlice(name)
	return val
}

func parseJSONFlag(s string, v interface{}) error {
	if err := json.Unmarshal([]byte(s), v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// addTreeSubResources adds ancestors, descendants, direct-children, operations, and properties
// sub-commands to a parent command for tree-like resources (whiteboards, databases, folders, smart links).
func addTreeSubResources(parentCmd *cobra.Command, pathPrefix, resourceName string) {
	// ancestors
	ancestorsCmd := &cobra.Command{
		Use:   "ancestors [id]",
		Short: fmt.Sprintf("Get all ancestors of %s", resourceName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit := getIntFlag(cmd, "limit"); limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			data, err := confGet(cmd, pathPrefix+"/"+args[0]+"/ancestors", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	ancestorsCmd.Flags().Int("limit", 25, "Maximum number of results")
	parentCmd.AddCommand(ancestorsCmd)

	// descendants
	descendantsCmd := &cobra.Command{
		Use:   "descendants [id]",
		Short: fmt.Sprintf("Get descendants of a %s", resourceName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit := getIntFlag(cmd, "limit"); limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			if depth := getIntFlag(cmd, "depth"); depth > 0 {
				q.Set("depth", fmt.Sprintf("%d", depth))
			}
			if cursor := getStringFlag(cmd, "cursor"); cursor != "" {
				q.Set("cursor", cursor)
			}
			data, err := confGet(cmd, pathPrefix+"/"+args[0]+"/descendants", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	descendantsCmd.Flags().Int("limit", 25, "Maximum number of results")
	descendantsCmd.Flags().Int("depth", 0, "Maximum depth of descendants")
	descendantsCmd.Flags().String("cursor", "", "Pagination cursor")
	parentCmd.AddCommand(descendantsCmd)

	// direct-children
	directChildrenCmd := &cobra.Command{
		Use:   "direct-children [id]",
		Short: fmt.Sprintf("Get direct children of a %s", resourceName),
		Args:  cobra.ExactArgs(1),
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
			data, err := confGet(cmd, pathPrefix+"/"+args[0]+"/direct-children", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(directChildrenCmd)
	addSortFlag(directChildrenCmd)
	parentCmd.AddCommand(directChildrenCmd)

	// operations
	operationsCmd := &cobra.Command{
		Use:   "operations [id]",
		Short: fmt.Sprintf("Get permitted operations for %s", resourceName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, pathPrefix+"/"+args[0]+"/operations", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	parentCmd.AddCommand(operationsCmd)

	// properties (list only - CRUD is in conf_properties.go)
	propertiesCmd := &cobra.Command{
		Use:   "properties [id]",
		Short: fmt.Sprintf("Get content properties for %s", resourceName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if k := getStringFlag(cmd, "key"); k != "" {
				q.Set("key", k)
			}
			data, err := confGet(cmd, pathPrefix+"/"+args[0]+"/properties", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(propertiesCmd)
	addSortFlag(propertiesCmd)
	propertiesCmd.Flags().String("key", "", "Filter by property key")
	parentCmd.AddCommand(propertiesCmd)
}
