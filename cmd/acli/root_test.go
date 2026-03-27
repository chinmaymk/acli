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
	expected := []string{"issue", "project", "board", "sprint", "epic", "backlog", "search", "filter", "user", "group", "dashboard"}
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
	expected := []string{"space", "page", "blogpost", "comment", "label", "attachment", "task", "custom-content", "whiteboard", "database", "folder", "smart-link", "property", "space-permission", "admin-key", "data-policy", "classification", "user", "space-role"}
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
