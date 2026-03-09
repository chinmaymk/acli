package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var confluenceCmd = &cobra.Command{
	Use:     "confluence",
	Aliases: []string{"conf", "c"},
	Short:   "Interact with Confluence Cloud",
	Long:    "Manage Confluence spaces, pages, and content.",
	RunE:    helpRunE,
}

var confSpaceCmd = &cobra.Command{
	Use:     "space",
	Short:   "Manage spaces",
	Aliases: []string{"s"},
	RunE:    helpRunE,
}

var confPageCmd = &cobra.Command{
	Use:     "page",
	Short:   "Manage pages",
	Aliases: []string{"p"},
	RunE:    helpRunE,
}

func init() {
	confluenceCmd.AddCommand(confSpaceCmd)
	confluenceCmd.AddCommand(confPageCmd)

	confPageCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "List pages",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: list Confluence pages")
		},
	})
	confPageCmd.AddCommand(&cobra.Command{
		Use:   "get [page-id]",
		Short: "Get page details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("TODO: get Confluence page %s\n", args[0])
		},
	})
	confPageCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a page",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: create Confluence page")
		},
	})
}
