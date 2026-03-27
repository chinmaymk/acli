package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bitbucketCmd = &cobra.Command{
	Use:     "bitbucket",
	Aliases: []string{"bb"},
	Short:   "Interact with Bitbucket Cloud",
	Long:    "Manage Bitbucket repositories, pull requests, pipelines, branches, and more.",
	RunE:    helpRunE,
}

var bbWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show the current authenticated Bitbucket user",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getBitbucketClient(cmd)
		if err != nil {
			return err
		}
		user, err := client.GetCurrentUser()
		if err != nil {
			return err
		}
		if isJSONOutput(cmd) {
			return outputJSON(user)
		}
		fmt.Printf("Display name: %s\n", user.DisplayName)
		fmt.Printf("Username:     %s\n", user.Nickname)
		fmt.Printf("UUID:         %s\n", user.UUID)
		fmt.Printf("Account ID:   %s\n", user.AccountID)
		return nil
	},
}

func init() {
	bitbucketCmd.AddCommand(bbWhoamiCmd)
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
