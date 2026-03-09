package acli

import (
	"github.com/spf13/cobra"
)

var bitbucketCmd = &cobra.Command{
	Use:     "bitbucket",
	Aliases: []string{"bb"},
	Short:   "Interact with Bitbucket Cloud",
	Long:    "Manage Bitbucket repositories, pull requests, pipelines, branches, and more.",
	RunE:    helpRunE,
}

func init() {
	bitbucketCmd.AddCommand(bbRepoCmd)
	bitbucketCmd.AddCommand(bbPRCmd)
	bitbucketCmd.AddCommand(bbPipelineCmd)
	bitbucketCmd.AddCommand(bbBranchCmd)
	bitbucketCmd.AddCommand(bbTagCmd)
	bitbucketCmd.AddCommand(bbCommitCmd)
	bitbucketCmd.AddCommand(bbWorkspaceCmd)
	bitbucketCmd.AddCommand(bbProjectCmd)
	bitbucketCmd.AddCommand(bbWebhookCmd)
	bitbucketCmd.AddCommand(bbEnvironmentCmd)
	bitbucketCmd.AddCommand(bbDeployKeyCmd)
	bitbucketCmd.AddCommand(bbDownloadCmd)
	bitbucketCmd.AddCommand(bbSnippetCmd)
	bitbucketCmd.AddCommand(bbIssueCmd)
	bitbucketCmd.AddCommand(bbSearchCmd)
	bitbucketCmd.AddCommand(bbDeploymentCmd)
	bitbucketCmd.AddCommand(bbBranchRestrictionCmd)
}
