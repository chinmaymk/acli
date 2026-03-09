package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// database create
	createDatabaseCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if getBoolFlag(cmd, "private") {
				q.Set("private", "true")
			}
			body := map[string]interface{}{
				"spaceId": getStringFlag(cmd, "space-id"),
			}
			if t := getStringFlag(cmd, "title"); t != "" {
				body["title"] = t
			}
			if p := getStringFlag(cmd, "parent-id"); p != "" {
				body["parentId"] = p
			}
			data, err := confPost(cmd, "/databases", q, body)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	createDatabaseCmd.Flags().String("space-id", "", "Space ID (required)")
	createDatabaseCmd.Flags().String("title", "", "Database title")
	createDatabaseCmd.Flags().String("parent-id", "", "Parent ID")
	createDatabaseCmd.Flags().Bool("private", false, "Create as private")
	_ = createDatabaseCmd.MarkFlagRequired("space-id")
	confDatabaseCmd.AddCommand(createDatabaseCmd)

	// database get
	getDatabaseCmd := &cobra.Command{
		Use:   "get [database-id]",
		Short: "Get database by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			for _, flag := range []string{"include-collaborators", "include-direct-children", "include-operations", "include-properties"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/databases/"+args[0], q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	getDatabaseCmd.Flags().Bool("include-collaborators", false, "Include collaborators")
	getDatabaseCmd.Flags().Bool("include-direct-children", false, "Include direct children")
	getDatabaseCmd.Flags().Bool("include-operations", false, "Include operations")
	getDatabaseCmd.Flags().Bool("include-properties", false, "Include properties")
	confDatabaseCmd.AddCommand(getDatabaseCmd)

	// database delete
	deleteDatabaseCmd := &cobra.Command{
		Use:   "delete [database-id]",
		Short: "Delete a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := confDelete(cmd, "/databases/"+args[0], nil)
			if err != nil {
				return err
			}
			fmt.Println("Database deleted successfully.")
			return nil
		},
	}
	confDatabaseCmd.AddCommand(deleteDatabaseCmd)

	addTreeSubResources(confDatabaseCmd, "/databases", "database")
}
