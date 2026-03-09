package acli

import (
	"fmt"
	"strings"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search code in a workspace",
	RunE:  helpRunE,
}

func init() {
	// search code
	searchCodeCmd := &cobra.Command{
		Use:   "code <workspace>",
		Short: "Search for code in a workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			query, _ := cmd.Flags().GetString("query")
			if query == "" {
				return fmt.Errorf("--query is required")
			}

			results, err := client.SearchCode(args[0], query)
			if err != nil {
				return err
			}

			fmt.Printf("Found %d results\n\n", results.Size)
			for _, r := range results.Values {
				fmt.Printf("File: %s (%d matches)\n", r.File.Path, r.ContentMatchCount)
				for _, m := range r.ContentMatches {
					for _, line := range m.Lines {
						var parts []string
						for _, seg := range line.Segments {
							if seg.Match {
								parts = append(parts, fmt.Sprintf("[%s]", seg.Text))
							} else {
								parts = append(parts, seg.Text)
							}
						}
						fmt.Printf("  %d: %s\n", line.Line, strings.Join(parts, ""))
					}
				}
				fmt.Println()
			}
			return nil
		},
	}
	searchCodeCmd.Flags().String("query", "", "Search query (required)")
	bbSearchCmd.AddCommand(searchCodeCmd)
}
