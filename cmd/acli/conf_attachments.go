package acli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	// attachment list
	listAttachmentsCmd := &cobra.Command{
		Use:     "list",
		Short:   "List attachments",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if m := getStringFlag(cmd, "media-type"); m != "" {
				q.Set("mediaType", m)
			}
			if f := getStringFlag(cmd, "filename"); f != "" {
				q.Set("filename", f)
			}
			data, err := confGet(cmd, "/attachments", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(listAttachmentsCmd)
	addSortFlag(listAttachmentsCmd)
	addStatusFlag(listAttachmentsCmd)
	listAttachmentsCmd.Flags().String("media-type", "", "Filter by media type")
	listAttachmentsCmd.Flags().String("filename", "", "Filter by filename")
	confAttachmentCmd.AddCommand(listAttachmentsCmd)

	// attachment get
	getAttachmentCmd := &cobra.Command{
		Use:   "get [attachment-id]",
		Short: "Get attachment by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if v := getIntFlag(cmd, "version"); v > 0 {
				q.Set("version", fmt.Sprintf("%d", v))
			}
			for _, flag := range []string{"include-labels", "include-properties", "include-operations",
				"include-versions", "include-version", "include-collaborators"} {
				if getBoolFlag(cmd, flag) {
					q.Set(flag, "true")
				}
			}
			data, err := confGet(cmd, "/attachments/"+args[0], q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	getAttachmentCmd.Flags().Int("version", 0, "Retrieve a specific version")
	getAttachmentCmd.Flags().Bool("include-labels", false, "Include labels")
	getAttachmentCmd.Flags().Bool("include-properties", false, "Include properties")
	getAttachmentCmd.Flags().Bool("include-operations", false, "Include operations")
	getAttachmentCmd.Flags().Bool("include-versions", false, "Include versions")
	getAttachmentCmd.Flags().Bool("include-version", false, "Include current version")
	getAttachmentCmd.Flags().Bool("include-collaborators", false, "Include collaborators")
	confAttachmentCmd.AddCommand(getAttachmentCmd)

	// attachment delete
	deleteAttachmentCmd := &cobra.Command{
		Use:   "delete [attachment-id]",
		Short: "Delete attachment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if getBoolFlag(cmd, "purge") {
				q.Set("purge", "true")
			}
			_, err := confDelete(cmd, "/attachments/"+args[0], q)
			if err != nil {
				return err
			}
			fmt.Println("Attachment deleted successfully.")
			return nil
		},
	}
	deleteAttachmentCmd.Flags().Bool("purge", false, "Purge the attachment")
	confAttachmentCmd.AddCommand(deleteAttachmentCmd)

	// attachment labels
	attLabelsCmd := &cobra.Command{
		Use:   "labels [attachment-id]",
		Short: "Get labels for attachment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if p := getStringFlag(cmd, "prefix"); p != "" {
				q.Set("prefix", p)
			}
			data, err := confGet(cmd, "/attachments/"+args[0]+"/labels", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(attLabelsCmd)
	addSortFlag(attLabelsCmd)
	attLabelsCmd.Flags().String("prefix", "", "Filter by prefix")
	confAttachmentCmd.AddCommand(attLabelsCmd)

	// attachment comments
	attCommentsCmd := &cobra.Command{
		Use:   "comments [attachment-id]",
		Short: "Get footer comments for attachment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			if v := getIntFlag(cmd, "version"); v > 0 {
				q.Set("version", fmt.Sprintf("%d", v))
			}
			data, err := confGet(cmd, "/attachments/"+args[0]+"/footer-comments", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(attCommentsCmd)
	addSortFlag(attCommentsCmd)
	addBodyFormatFlag(attCommentsCmd)
	attCommentsCmd.Flags().Int("version", 0, "Filter by version")
	confAttachmentCmd.AddCommand(attCommentsCmd)

	// attachment operations
	attOpsCmd := &cobra.Command{
		Use:   "operations [attachment-id]",
		Short: "Get permitted operations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/attachments/"+args[0]+"/operations", nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confAttachmentCmd.AddCommand(attOpsCmd)

	// attachment versions
	attVersionsCmd := &cobra.Command{
		Use:   "versions [attachment-id]",
		Short: "Get attachment versions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := getPaginationQuery(cmd)
			data, err := confGet(cmd, "/attachments/"+args[0]+"/versions", q)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	addPaginationFlags(attVersionsCmd)
	addSortFlag(attVersionsCmd)
	confAttachmentCmd.AddCommand(attVersionsCmd)

	// attachment version-details
	attVersionDetailCmd := &cobra.Command{
		Use:   "version-details [attachment-id] [version-number]",
		Short: "Get version details for attachment version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := confGet(cmd, "/attachments/"+args[0]+"/versions/"+args[1], nil)
			if err != nil {
				return err
			}
			printJSONBytes(data)
			return nil
		},
	}
	confAttachmentCmd.AddCommand(attVersionDetailCmd)
}
