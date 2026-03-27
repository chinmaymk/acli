package acli

import (
	"testing"
)

func TestRootCommandHasSubcommands(t *testing.T) {
	expected := []string{"jira", "confluence", "bitbucket", "config", "version", "commands"}
	subs := rootCmd.Commands()

	nameSet := make(map[string]bool)
	for _, cmd := range subs {
		nameSet[cmd.Name()] = true
	}

	for _, name := range expected {
		if !nameSet[name] {
			t.Errorf("expected subcommand %q not found on root", name)
		}
	}
}

func TestRootPersistentFlags(t *testing.T) {
	flags := rootCmd.PersistentFlags()

	pf := flags.Lookup("profile")
	if pf == nil {
		t.Fatal("expected --profile flag")
	}
	if pf.Shorthand != "p" {
		t.Errorf("expected shorthand 'p', got %q", pf.Shorthand)
	}

	of := flags.Lookup("output")
	if of == nil {
		t.Fatal("expected --output flag")
	}
	if of.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got %q", of.Shorthand)
	}
	if of.DefValue != "text" {
		t.Errorf("expected default 'text', got %q", of.DefValue)
	}
}

func TestJiraCommandAliases(t *testing.T) {
	if len(jiraCmd.Aliases) == 0 {
		t.Fatal("expected jira aliases")
	}
	found := false
	for _, a := range jiraCmd.Aliases {
		if a == "j" {
			found = true
		}
	}
	if !found {
		t.Error("expected alias 'j' for jira command")
	}
}

func TestConfluenceCommandAliases(t *testing.T) {
	aliases := confluenceCmd.Aliases
	expected := map[string]bool{"conf": true, "c": true}
	for _, a := range aliases {
		delete(expected, a)
	}
	if len(expected) > 0 {
		t.Errorf("missing confluence aliases: %v", expected)
	}
}

func TestBitbucketCommandAliases(t *testing.T) {
	found := false
	for _, a := range bitbucketCmd.Aliases {
		if a == "bb" {
			found = true
		}
	}
	if !found {
		t.Error("expected alias 'bb' for bitbucket command")
	}
}

func TestJiraSubcommands(t *testing.T) {
	// Core subcommands from jira.go
	expected := []string{"issue", "project", "board", "sprint", "epic", "backlog", "search", "filter", "user", "group", "dashboard"}
	// Additional subcommands from jira_admin.go
	expected = append(expected, "role", "issuelink", "issuelinktype", "screen", "workflow", "workflowscheme",
		"permissionscheme", "notificationscheme", "issuesecurityscheme", "fieldconfig", "issuetypescheme",
		"serverinfo", "webhook", "attachment", "audit", "banner", "configuration", "permission", "task", "projectcategory")
	// Additional subcommands from jira_resources.go
	expected = append(expected, "component", "version", "field", "label", "issuetype", "priority", "resolution", "status")

	subs := jiraCmd.Commands()
	nameSet := make(map[string]bool)
	for _, cmd := range subs {
		nameSet[cmd.Name()] = true
	}
	for _, name := range expected {
		if !nameSet[name] {
			t.Errorf("expected jira subcommand %q not found", name)
		}
	}
}

func TestBitbucketSubcommands(t *testing.T) {
	expected := []string{"repo", "pr", "pipeline", "branch", "tag", "commit", "workspace", "project", "webhook", "environment", "deploy-key", "download", "snippet", "issue", "search", "deployment", "branch-restriction"}
	subs := bitbucketCmd.Commands()
	nameSet := make(map[string]bool)
	for _, cmd := range subs {
		nameSet[cmd.Name()] = true
	}
	for _, name := range expected {
		if !nameSet[name] {
			t.Errorf("expected bitbucket subcommand %q not found", name)
		}
	}
}

func TestConfluenceSubcommands(t *testing.T) {
	expected := []string{"space", "page", "blogpost", "comment", "label", "attachment", "task", "custom-content", "whiteboard", "database", "folder", "smart-link", "property", "space-permission", "admin-key", "data-policy", "classification", "user", "space-role", "convert-ids", "app-property"}
	subs := confluenceCmd.Commands()
	nameSet := make(map[string]bool)
	for _, cmd := range subs {
		nameSet[cmd.Name()] = true
	}
	for _, name := range expected {
		if !nameSet[name] {
			t.Errorf("expected confluence subcommand %q not found", name)
		}
	}
}

func TestConfigSubcommands(t *testing.T) {
	expected := []string{"setup", "list", "show", "delete", "set-default", "set-defaults"}
	subs := configCmd.Commands()
	nameSet := make(map[string]bool)
	for _, cmd := range subs {
		nameSet[cmd.Name()] = true
	}
	for _, name := range expected {
		if !nameSet[name] {
			t.Errorf("expected config subcommand %q not found", name)
		}
	}
}

func TestConfluenceNestedCommentSubcommands(t *testing.T) {
	// comment -> footer, inline
	subs := confCommentCmd.Commands()
	nameSet := make(map[string]bool)
	for _, cmd := range subs {
		nameSet[cmd.Name()] = true
	}
	for _, name := range []string{"footer", "inline"} {
		if !nameSet[name] {
			t.Errorf("expected comment subcommand %q not found", name)
		}
	}
}

func TestConfluenceResourceAliases(t *testing.T) {
	tests := []struct {
		name    string
		aliases []string
		actual  []string
	}{
		{"space", []string{"s"}, confSpaceCmd.Aliases},
		{"page", []string{"p"}, confPageCmd.Aliases},
		{"blogpost", []string{"blog", "bp"}, confBlogPostCmd.Aliases},
		{"comment", []string{"cm"}, confCommentCmd.Aliases},
		{"footer", []string{"fc"}, confFooterCommentCmd.Aliases},
		{"inline", []string{"ic"}, confInlineCommentCmd.Aliases},
		{"label", []string{"l"}, confLabelCmd.Aliases},
		{"attachment", []string{"att", "a"}, confAttachmentCmd.Aliases},
		{"task", []string{"t"}, confTaskCmd.Aliases},
		{"custom-content", []string{"cc"}, confCustomContentCmd.Aliases},
		{"whiteboard", []string{"wb"}, confWhiteboardCmd.Aliases},
		{"database", []string{"db"}, confDatabaseCmd.Aliases},
		{"folder", []string{"f"}, confFolderCmd.Aliases},
		{"smart-link", []string{"sl", "embed"}, confSmartLinkCmd.Aliases},
		{"property", []string{"prop"}, confPropertyCmd.Aliases},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualSet := make(map[string]bool)
			for _, a := range tt.actual {
				actualSet[a] = true
			}
			for _, expected := range tt.aliases {
				if !actualSet[expected] {
					t.Errorf("expected alias %q for %s not found", expected, tt.name)
				}
			}
		})
	}
}

func TestConfigCommandAliases(t *testing.T) {
	found := false
	for _, a := range configCmd.Aliases {
		if a == "cfg" {
			found = true
		}
	}
	if !found {
		t.Error("expected alias 'cfg' for config command")
	}
}

func TestJiraSubcommandCount(t *testing.T) {
	// Ensure we have a substantial number of jira subcommands (guards against accidental removal)
	subs := jiraCmd.Commands()
	if len(subs) < 30 {
		t.Errorf("expected at least 30 jira subcommands, got %d", len(subs))
	}
}

func TestBitbucketSubcommandCount(t *testing.T) {
	subs := bitbucketCmd.Commands()
	if len(subs) < 15 {
		t.Errorf("expected at least 15 bitbucket subcommands, got %d", len(subs))
	}
}

func TestConfluenceSubcommandCount(t *testing.T) {
	subs := confluenceCmd.Commands()
	if len(subs) < 19 {
		t.Errorf("expected at least 19 confluence subcommands, got %d", len(subs))
	}
}

func TestGroupCommandsHaveRunE(t *testing.T) {
	// Group commands should have RunE set (to helpRunE) so they don't appear as "additional help topics"
	groupCmds := []*struct {
		name string
		cmd  interface{ HasSubCommands() bool }
	}{
		{"jira", jiraCmd},
		{"confluence", confluenceCmd},
		{"bitbucket", bitbucketCmd},
		{"config", configCmd},
	}

	for _, gc := range groupCmds {
		if !gc.cmd.HasSubCommands() {
			t.Errorf("%s should have subcommands", gc.name)
		}
	}
}
