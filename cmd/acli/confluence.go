package acli

import (
	"github.com/spf13/cobra"
)

var confluenceCmd = &cobra.Command{
	Use:     "confluence",
	Aliases: []string{"conf", "c"},
	Short:   "Interact with Confluence Cloud",
	Long:    "Manage Confluence spaces, pages, blog posts, comments, labels, attachments, tasks, and more.",
	RunE:    helpRunE,
}

// Resource group commands
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

var confBlogPostCmd = &cobra.Command{
	Use:     "blogpost",
	Short:   "Manage blog posts",
	Aliases: []string{"blog", "bp"},
	RunE:    helpRunE,
}

var confCommentCmd = &cobra.Command{
	Use:     "comment",
	Short:   "Manage comments (footer and inline)",
	Aliases: []string{"cm"},
	RunE:    helpRunE,
}

var confFooterCommentCmd = &cobra.Command{
	Use:     "footer",
	Short:   "Manage footer comments",
	Aliases: []string{"fc"},
	RunE:    helpRunE,
}

var confInlineCommentCmd = &cobra.Command{
	Use:     "inline",
	Short:   "Manage inline comments",
	Aliases: []string{"ic"},
	RunE:    helpRunE,
}

var confLabelCmd = &cobra.Command{
	Use:     "label",
	Short:   "Manage labels",
	Aliases: []string{"l"},
	RunE:    helpRunE,
}

var confAttachmentCmd = &cobra.Command{
	Use:     "attachment",
	Short:   "Manage attachments",
	Aliases: []string{"att", "a"},
	RunE:    helpRunE,
}

var confTaskCmd = &cobra.Command{
	Use:     "task",
	Short:   "Manage tasks",
	Aliases: []string{"t"},
	RunE:    helpRunE,
}

var confCustomContentCmd = &cobra.Command{
	Use:     "custom-content",
	Short:   "Manage custom content",
	Aliases: []string{"cc"},
	RunE:    helpRunE,
}

var confWhiteboardCmd = &cobra.Command{
	Use:     "whiteboard",
	Short:   "Manage whiteboards",
	Aliases: []string{"wb"},
	RunE:    helpRunE,
}

var confDatabaseCmd = &cobra.Command{
	Use:     "database",
	Short:   "Manage databases",
	Aliases: []string{"db"},
	RunE:    helpRunE,
}

var confFolderCmd = &cobra.Command{
	Use:     "folder",
	Short:   "Manage folders",
	Aliases: []string{"f"},
	RunE:    helpRunE,
}

var confSmartLinkCmd = &cobra.Command{
	Use:     "smart-link",
	Short:   "Manage smart links (embeds)",
	Aliases: []string{"sl", "embed"},
	RunE:    helpRunE,
}

var confPropertyCmd = &cobra.Command{
	Use:     "property",
	Short:   "Manage content properties",
	Aliases: []string{"prop"},
	RunE:    helpRunE,
}

var confSpacePermissionCmd = &cobra.Command{
	Use:   "space-permission",
	Short: "Manage space permissions and roles",
	Aliases: []string{"sp"},
	RunE:  helpRunE,
}

var confAdminKeyCmd = &cobra.Command{
	Use:   "admin-key",
	Short: "Manage admin key",
	Aliases: []string{"ak"},
	RunE:  helpRunE,
}

var confDataPolicyCmd = &cobra.Command{
	Use:   "data-policy",
	Short: "Manage data policies",
	Aliases: []string{"dp"},
	RunE:  helpRunE,
}

var confClassificationCmd = &cobra.Command{
	Use:   "classification",
	Short: "Manage classification levels",
	Aliases: []string{"cl"},
	RunE:  helpRunE,
}

var confUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user access",
	Aliases: []string{"u"},
	RunE:  helpRunE,
}

var confSpaceRoleCmd = &cobra.Command{
	Use:   "space-role",
	Short: "Manage space roles",
	Aliases: []string{"sr"},
	RunE:  helpRunE,
}

func init() {
	confluenceCmd.AddCommand(confSpaceCmd)
	confluenceCmd.AddCommand(confPageCmd)
	confluenceCmd.AddCommand(confBlogPostCmd)
	confluenceCmd.AddCommand(confCommentCmd)
	confluenceCmd.AddCommand(confLabelCmd)
	confluenceCmd.AddCommand(confAttachmentCmd)
	confluenceCmd.AddCommand(confTaskCmd)
	confluenceCmd.AddCommand(confCustomContentCmd)
	confluenceCmd.AddCommand(confWhiteboardCmd)
	confluenceCmd.AddCommand(confDatabaseCmd)
	confluenceCmd.AddCommand(confFolderCmd)
	confluenceCmd.AddCommand(confSmartLinkCmd)
	confluenceCmd.AddCommand(confPropertyCmd)
	confluenceCmd.AddCommand(confSpacePermissionCmd)
	confluenceCmd.AddCommand(confAdminKeyCmd)
	confluenceCmd.AddCommand(confDataPolicyCmd)
	confluenceCmd.AddCommand(confClassificationCmd)
	confluenceCmd.AddCommand(confUserCmd)
	confluenceCmd.AddCommand(confSpaceRoleCmd)

	confCommentCmd.AddCommand(confFooterCommentCmd)
	confCommentCmd.AddCommand(confInlineCommentCmd)
}
