package acli

import (
	"bytes"
	"testing"

	"github.com/chinmaymk/acli/internal/jira"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a long string", 10, "this is..."},
		{"ab", 3, "ab"},
		{"abcd", 3, "abc"},
		{"abcdef", 4, "a..."},
		{"", 5, ""},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

func TestExtractArgsFromUse(t *testing.T) {
	tests := []struct {
		use  string
		want string
	}{
		{"get <issue-key>", "<issue-key>"},
		{"list", ""},
		{"create <workspace> <repo>", "<workspace> <repo>"},
		{"delete <id>", "<id>"},
	}

	for _, tt := range tests {
		got := extractArgsFromUse(tt.use)
		if got != tt.want {
			t.Errorf("extractArgsFromUse(%q) = %q, want %q", tt.use, got, tt.want)
		}
	}
}

func TestHelpRunE(t *testing.T) {
	// helpRunE should not return an error (it prints help)
	// We can test it via the jira command which uses it
	err := jiraCmd.RunE(jiraCmd, nil)
	if err != nil {
		t.Errorf("helpRunE should not error, got: %v", err)
	}
}

func TestPrintIssueRow(t *testing.T) {
	var buf bytes.Buffer
	w := newTabWriter()
	// Replace stdout with our buffer for testing
	// We just verify it doesn't panic with nil fields
	issue := jira.IssueDetailed{
		Key: "TEST-1",
		Fields: jira.IssueFields{
			Summary: "Test issue",
		},
	}
	printIssueRow(w, issue)
	w.Flush()
	// If we get here without panic, the nil handling works
	_ = buf
}

func TestPrintIssueRowWithAllFields(t *testing.T) {
	w := newTabWriter()
	issue := jira.IssueDetailed{
		Key: "TEST-2",
		Fields: jira.IssueFields{
			Summary:   "Full issue",
			IssueType: &jira.IssueType{Name: "Bug"},
			Status:    &jira.StatusDetails{Name: "Open"},
			Priority:  &jira.Priority{Name: "High"},
			Assignee:  &jira.UserDetails{DisplayName: "John"},
		},
	}
	printIssueRow(w, issue)
	w.Flush()
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"short", "****"},
		{"12345678", "****"},
		{"123456789", "1234****6789"},
		{"abcdefghijklmnop", "abcd****mnop"},
	}

	for _, tt := range tests {
		got := maskToken(tt.input)
		if got != tt.want {
			t.Errorf("maskToken(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
