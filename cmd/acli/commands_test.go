package acli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildCommandTree(t *testing.T) {
	tree := buildCommandTree(rootCmd)

	if tree.Name != "acli" {
		t.Errorf("expected root name 'acli', got %q", tree.Name)
	}
	if tree.Description == "" {
		t.Error("expected non-empty description")
	}
	if len(tree.Subcommands) == 0 {
		t.Error("expected subcommands in tree")
	}

	// Check that jira, confluence, bitbucket are in the tree
	found := map[string]bool{}
	for _, sub := range tree.Subcommands {
		found[sub.Name] = true
	}
	for _, name := range []string{"jira", "confluence", "bitbucket", "config", "version", "commands"} {
		if !found[name] {
			t.Errorf("expected %q in command tree", name)
		}
	}
}

func TestBuildCommandTreeRecursive(t *testing.T) {
	tree := buildCommandTree(rootCmd)

	// Find the jira subtree
	var jiraTree *CommandInfo
	for i, sub := range tree.Subcommands {
		if sub.Name == "jira" {
			jiraTree = &tree.Subcommands[i]
			break
		}
	}
	if jiraTree == nil {
		t.Fatal("jira not found in command tree")
	}
	if len(jiraTree.Aliases) == 0 {
		t.Error("expected jira aliases in tree")
	}
	if len(jiraTree.Subcommands) == 0 {
		t.Error("expected jira subcommands in tree")
	}
}

func TestBuildCommandTreeHidesHelpAndCompletion(t *testing.T) {
	tree := buildCommandTree(rootCmd)

	for _, sub := range tree.Subcommands {
		if sub.Name == "help" || sub.Name == "completion" {
			t.Errorf("command tree should not include %q", sub.Name)
		}
	}
}

func TestBuildCommandTreeFlags(t *testing.T) {
	// Create a test command with flags
	cmd := &cobra.Command{
		Use:   "test-cmd <arg>",
		Short: "A test command",
		Args:  cobra.ExactArgs(1),
	}
	cmd.Flags().StringP("name", "n", "", "the name")
	cmd.Flags().Bool("verbose", false, "verbose output")

	tree := buildCommandTree(cmd)

	if tree.Args != "<arg>" {
		t.Errorf("expected args '<arg>', got %q", tree.Args)
	}
	if len(tree.Flags) < 2 {
		t.Errorf("expected at least 2 flags, got %d", len(tree.Flags))
	}

	flagMap := map[string]FlagInfo{}
	for _, f := range tree.Flags {
		flagMap[f.Name] = f
	}

	if nf, ok := flagMap["name"]; ok {
		if nf.Shorthand != "n" {
			t.Errorf("expected shorthand 'n', got %q", nf.Shorthand)
		}
		if nf.Type != "string" {
			t.Errorf("expected type 'string', got %q", nf.Type)
		}
	} else {
		t.Error("expected 'name' flag in tree")
	}
}

func TestExtractArgsFromUseNoArgs(t *testing.T) {
	if got := extractArgsFromUse("list"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
