package acli

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

// confSearchCmd implements `acli confluence search` (CQL).
// The Confluence v1 REST API (/wiki/rest/api/content/search?cql=...) is the
// only supported endpoint for CQL content search; the v2 API does not expose
// it yet, so we use ConfluenceV1 rather than the standard confGet helper.
var confSearchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Search content using CQL",
	Long: `Search Confluence content using the Confluence Query Language (CQL).

Examples:
  acli confluence search --cql 'type=page AND space=DEV'
  acli confluence search --cql 'title~"Meeting Notes"' --limit 20
  acli confluence search --cql 'label=documentation' --start 25 --all
  acli confluence search --cql 'text~"important concept"' --expand "body.view" --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient(cmd)
		if err != nil {
			return err
		}

		cql, _ := cmd.Flags().GetString("cql")
		limit, _ := cmd.Flags().GetInt("limit")
		start, _ := cmd.Flags().GetInt("start")
		expand, _ := cmd.Flags().GetString("expand")
		all, _ := cmd.Flags().GetBool("all")

		query := url.Values{}
		query.Set("cql", cql)
		query.Set("limit", fmt.Sprintf("%d", limit))
		query.Set("start", fmt.Sprintf("%d", start))
		if expand != "" {
			query.Set("expand", expand)
		}

		raw, err := client.ConfluenceV1("GET", "/content/search", query, nil)
		if err != nil {
			return err
		}

		if !all {
			if isJSONOutput(cmd) {
				printJSONBytes(raw)
				return nil
			}
			return printSearchTable(raw)
		}

		type pageResults struct {
			Results []json.RawMessage `json:"results"`
			Size    int               `json:"size"`
			Total   int               `json:"totalSize"`
		}

		var firstPage pageResults
		if err := json.Unmarshal(raw, &firstPage); err != nil {
			return fmt.Errorf("parsing search response: %w", err)
		}

		combined := firstPage.Results
		nextStart := start + firstPage.Size
		for firstPage.Size > 0 && len(combined) < firstPage.Total {
			query.Set("start", fmt.Sprintf("%d", nextStart))
			nextRaw, err := client.ConfluenceV1("GET", "/content/search", query, nil)
			if err != nil {
				return err
			}
			var nextPage pageResults
			if err := json.Unmarshal(nextRaw, &nextPage); err != nil {
				return fmt.Errorf("parsing paginated search response at start=%d: %w", nextStart, err)
			}
			if len(nextPage.Results) == 0 {
				break
			}
			combined = append(combined, nextPage.Results...)
			nextStart += nextPage.Size
		}

		if isJSONOutput(cmd) {
			return outputJSON(combined)
		}

		combinedRaw, err := json.Marshal(map[string]interface{}{
			"results":   combined,
			"totalSize": firstPage.Total,
			"size":      len(combined),
		})
		if err != nil {
			return outputJSON(combined)
		}
		return printSearchTable(combinedRaw)
	},
}

func printSearchTable(raw []byte) error {
	type contentItem struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Title string `json:"title"`
		Space struct {
			Key string `json:"key"`
		} `json:"space"`
		Links struct {
			WebUI string `json:"webui"`
		} `json:"_links"`
	}
	type searchPage struct {
		Results []contentItem `json:"results"`
		Size    int           `json:"size"`
		Total   int           `json:"totalSize"`
	}

	var page searchPage
	if err := json.Unmarshal(raw, &page); err != nil {
		printJSONBytes(raw)
		return nil
	}

	w := newTabWriter()
	_, _ = fmt.Fprintln(w, "ID\tTYPE\tSPACE\tTITLE")
	for _, item := range page.Results {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			item.ID, item.Type, item.Space.Key, truncate(item.Title, 60))
	}
	_ = w.Flush()
	if page.Total > page.Size {
		fmt.Printf("\nShowing %d of %d results. Use --start to paginate or --all to fetch everything.\n",
			page.Size, page.Total)
	}
	return nil
}

func init() {
	confSearchCmd.Flags().String("cql", "", "CQL query (required)")
	_ = confSearchCmd.MarkFlagRequired("cql")
	confSearchCmd.Flags().Int("limit", 25, "Maximum number of results per page")
	confSearchCmd.Flags().Int("start", 0, "Index of the first result (for pagination)")
	confSearchCmd.Flags().String("expand", "", "Comma-separated list of properties to expand (e.g. body.view,space)")
	confSearchCmd.Flags().Bool("all", false, "Fetch all results, overriding --limit and --start")
	confSearchCmd.Flags().Bool("json", false, "Output as JSON")
}
