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

func TestIssueColumnFor(t *testing.T) {
	issue := jira.IssueDetailed{
		Key: "TEST-1",
		Fields: jira.IssueFields{
			Summary:   "A summary",
			IssueType: &jira.IssueType{Name: "Bug"},
			Status:    &jira.StatusDetails{Name: "Open"},
			Priority:  &jira.Priority{Name: "High"},
			Assignee:  &jira.UserDetails{DisplayName: "Alice"},
			Labels:    []string{"a", "b"},
		},
	}

	tests := []struct {
		field      string
		wantHeader string
		wantValue  string
	}{
		{"summary", "SUMMARY", "A summary"},
		{"type", "TYPE", "Bug"},
		{"issuetype", "TYPE", "Bug"},
		{"status", "STATUS", "Open"},
		{"priority", "PRIORITY", "High"},
		{"assignee", "ASSIGNEE", "Alice"},
		{"labels", "LABELS", "a,b"},
		{"unknowncustom", "UNKNOWNCUSTOM", ""},
	}

	for _, tt := range tests {
		col := issueColumnFor(tt.field)
		if col.header != tt.wantHeader {
			t.Errorf("issueColumnFor(%q).header = %q, want %q", tt.field, col.header, tt.wantHeader)
		}
		if got := col.extract(issue); got != tt.wantValue {
			t.Errorf("issueColumnFor(%q).extract = %q, want %q", tt.field, got, tt.wantValue)
		}
	}

	// Verify nil-safe extractors don't panic on an issue with no parsed fields.
	bare := jira.IssueDetailed{Key: "X-1"}
	for _, f := range []string{"type", "status", "priority", "assignee", "reporter", "creator", "resolution", "project"} {
		if got := issueColumnFor(f).extract(bare); got != "" {
			t.Errorf("issueColumnFor(%q).extract on bare issue = %q, want empty", f, got)
		}
	}
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
