package acli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

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

func printJSON(data []byte) {
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

func printTable(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	fmt.Fprintln(w, strings.Repeat("-\t", len(headers)))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
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
