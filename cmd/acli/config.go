package acli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Manage ACLI configuration and profiles",
	Aliases: []string{"cfg"},
	RunE:    helpRunE,
}

var configSetupCmd = &cobra.Command{
	Use:   "setup <profile-name>",
	Short: "Create or update a configuration profile",
	Long: `Interactively set up a configuration profile with your Atlassian credentials.

The first profile created is automatically set as the default.

You'll be prompted for:
  - Atlassian instance URL (for Jira/Confluence)
  - Email address (leave blank for OAuth/scoped token auth)
  - API token or OAuth/scoped access token

The same credentials are used for Jira, Confluence, and Bitbucket.

Auth modes:
  - Email + API token → Basic Auth (personal API tokens)
  - Token only (no email) → Bearer Auth (OAuth 2.0 / scoped tokens)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		// Pre-fill with existing values if profile exists
		existing, exists := cfg.Profiles[profileName]
		if exists {
			fmt.Printf("Updating profile %q (press Enter to keep current value)\n\n", profileName)
		} else {
			fmt.Printf("Creating profile %q\n\n", profileName)
		}

		reader := bufio.NewReader(os.Stdin)

		atlassianURL := promptWithDefault(reader, "Atlassian URL", existing.AtlassianURL, "https://your-instance.atlassian.net")
		email := promptWithDefault(reader, "Email", existing.Email, "")
		apiToken := promptWithDefault(reader, "API Token", maskToken(existing.APIToken), "")
		// If user just pressed enter and there was an existing token, keep it
		if apiToken == maskToken(existing.APIToken) && existing.APIToken != "" {
			apiToken = existing.APIToken
		}

		profile := config.Profile{
			Name:         profileName,
			AtlassianURL: strings.TrimRight(atlassianURL, "/"),
			Email:        email,
			APIToken:     apiToken,
		}

		if cfg.Profiles == nil {
			cfg.Profiles = make(map[string]config.Profile)
		}

		// Set as default if it's the first profile or no default is set yet
		isFirst := len(cfg.Profiles) == 0
		cfg.Profiles[profileName] = profile
		if isFirst || cfg.DefaultProfile == "" {
			cfg.DefaultProfile = profileName
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("\nProfile %q saved to ~/.config/acli/config.json\n", profileName)
		if isFirst || cfg.DefaultProfile == profileName {
			fmt.Printf("Profile %q is the default profile\n", profileName)
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if len(cfg.Profiles) == 0 {
			fmt.Println("No profiles configured. Run 'acli config setup' to create one.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintf(w, "PROFILE\tDEFAULT\tURL\tEMAIL\n")
		for name, p := range cfg.Profiles {
			def := ""
			if name == cfg.DefaultProfile {
				def = "*"
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, def, p.AtlassianURL, p.Email)
		}
		return w.Flush()
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show [profile-name]",
	Short: "Show details for a profile (tokens are masked)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		profileName := ""
		if len(args) > 0 {
			profileName = args[0]
		}
		p, err := cfg.GetProfile(profileName)
		if err != nil {
			return err
		}

		fmt.Printf("Profile: %s\n", p.Name)
		fmt.Printf("  Atlassian URL:    %s\n", p.AtlassianURL)
		fmt.Printf("  Email:            %s\n", p.Email)
		fmt.Printf("  API Token:        %s\n", maskToken(p.APIToken))
		if p.Defaults.Project != "" || p.Defaults.Workspace != "" || p.Defaults.BBProject != "" {
			fmt.Println("  Defaults:")
			if p.Defaults.Project != "" {
				fmt.Printf("    Project:        %s\n", p.Defaults.Project)
			}
			if p.Defaults.Workspace != "" {
				fmt.Printf("    Workspace:      %s\n", p.Defaults.Workspace)
			}
			if p.Defaults.BBProject != "" {
				fmt.Printf("    BB Project:     %s\n", p.Defaults.BBProject)
			}
		}
		return nil
	},
}

var configDeleteCmd = &cobra.Command{
	Use:     "delete <profile-name>",
	Short:   "Delete a configuration profile",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if _, ok := cfg.Profiles[profileName]; !ok {
			return fmt.Errorf("profile %q not found", profileName)
		}

		delete(cfg.Profiles, profileName)

		if cfg.DefaultProfile == profileName {
			cfg.DefaultProfile = ""
			// If one profile remains, make it the new default
			for k := range cfg.Profiles {
				cfg.DefaultProfile = k
				break
			}
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Profile %q deleted\n", profileName)
		return nil
	},
}

var configSetDefaultCmd = &cobra.Command{
	Use:   "set-default <profile-name>",
	Short: "Set the default profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if _, ok := cfg.Profiles[profileName]; !ok {
			return fmt.Errorf("profile %q not found", profileName)
		}

		cfg.DefaultProfile = profileName

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Default profile set to %q\n", profileName)
		return nil
	},
}

var configSetDefaultsCmd = &cobra.Command{
	Use:   "set-defaults [profile-name]",
	Short: "Set default project and workspace for a profile",
	Long: `Set per-profile defaults so you don't have to specify --project or workspace on every command.

These defaults are used as fallbacks when the flag/argument is not provided:
  - project: Default Jira project key (used by issue list, issue create, etc.)
  - workspace: Default Bitbucket workspace (used by repo, pr, pipeline, etc.)
  - bb_project: Default Bitbucket project key (used by repo create, etc.)`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		profileName := ""
		if len(args) > 0 {
			profileName = args[0]
		}
		profile, err := cfg.GetProfile(profileName)
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Setting defaults for profile %q\n\n", profile.Name)

		project := promptWithDefault(reader, "Default Jira project key", profile.Defaults.Project, "")
		workspace := promptWithDefault(reader, "Default Bitbucket workspace", profile.Defaults.Workspace, "")
		bbProject := promptWithDefault(reader, "Default Bitbucket project key", profile.Defaults.BBProject, "")

		profile.Defaults = config.Defaults{
			Project:   project,
			Workspace: workspace,
			BBProject: bbProject,
		}

		cfg.Profiles[profile.Name] = profile
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("\nDefaults saved for profile %q\n", profile.Name)
		if project != "" {
			fmt.Printf("  Default project:   %s\n", project)
		}
		if workspace != "" {
			fmt.Printf("  Default workspace: %s\n", workspace)
		}
		if bbProject != "" {
			fmt.Printf("  Default BB project: %s\n", bbProject)
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetupCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configSetDefaultCmd)
	configCmd.AddCommand(configSetDefaultsCmd)
}

func promptWithDefault(reader *bufio.Reader, label, current, placeholder string) string {
	if current != "" {
		fmt.Printf("  %s [%s]: ", label, current)
	} else if placeholder != "" {
		fmt.Printf("  %s (%s): ", label, placeholder)
	} else {
		fmt.Printf("  %s: ", label)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		if current != "" {
			return current
		}
		return placeholder
	}
	return input
}

func maskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
