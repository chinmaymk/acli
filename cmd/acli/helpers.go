package acli

import "github.com/spf13/cobra"

// helpRunE is used for group commands that have no action themselves.
// It prints help so they show up as proper commands (not "additional help topics").
func helpRunE(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
