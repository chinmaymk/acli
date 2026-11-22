package acli

import (
	"testing"

	"github.com/chinmaymk/acli/internal/jira"
)

func TestTextSearchJQL(t *testing.T) {
	for _, tc := range []struct {
		name    string
		text    string
		filters issueSearchFilters
		want    string
	}{
		{"terms and filters", " release notes ", issueSearchFilters{project: "DEV", mine: true, updatedSince: "30d"},
			`text ~ "\"release\"" AND text ~ "\"notes\"" AND project = "DEV" AND assignee = currentUser() AND updated >= "-30d"`},
		{"issue key", "dev-123", issueSearchFilters{}, `key = "DEV-123"`},
		{"literal punctuation", `"café" OR (labels=ops)`, issueSearchFilters{},
			`text ~ "\"café\"" AND text ~ "\"OR\"" AND text ~ "\"labels\"" AND text ~ "\"ops\""`},
		{"filter only", "", issueSearchFilters{assignee: `a"b`}, `assignee = "a\"b"`},
		{"no input", "", issueSearchFilters{}, ""},
		{"no words", "...", issueSearchFilters{}, ""},
		{"conflicting assignee", "", issueSearchFilters{mine: true, assignee: "a"}, ""},
		{"invalid age", "term", issueSearchFilters{updatedSince: "-1d"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := textSearchJQL(tc.text, "all", tc.filters)
			if got != tc.want || (err != nil) != (tc.want == "") {
				t.Errorf("got %q, %v; want %q", got, err, tc.want)
			}
		})
	}
}

func TestIssuePreviews(t *testing.T) {
	issues := []jira.IssueDetailed{
		{Key: "DEV-1", Fields: jira.IssueFields{Summary: "Publish docs", Description: "Review release notes"}},
		{Key: "DEV-2", Fields: jira.IssueFields{Summary: "Release notes"}},
	}
	items := issuePreviews(issues, "release notes", "https://example.test", "all")
	if len(items) != 2 || items[0].Key != "DEV-2" || items[1].Key != "DEV-1" {
		t.Fatalf("unexpected preview order: %#v", items)
	}
	if items[0].Excerpt != "" || items[1].Excerpt != "Review release notes" || items[1].URL != "https://example.test/browse/DEV-1" {
		t.Errorf("unexpected previews: %#v", items)
	}
}

func TestAnyWordSearch(t *testing.T) {
	filters := issueSearchFilters{project: "DEV", mine: true, updatedSince: "2w"}
	got, err := textSearchJQL("draft notes", "any", filters)
	want := `(text ~ "\"draft\"" OR text ~ "\"notes\"") AND project = "DEV" AND assignee = currentUser() AND updated >= "-2w"`
	if err != nil || got != want {
		t.Fatalf("got %q, %v; want %q", got, err, want)
	}
	got, err = textSearchJQL("dev-21", "any", issueSearchFilters{})
	if err != nil || got != `key = "DEV-21"` {
		t.Errorf("unexpected key lookup: %q, %v", got, err)
	}
}

func TestAnyWordPreviewOrder(t *testing.T) {
	issues := []jira.IssueDetailed{
		{Key: "DEV-1", Fields: jira.IssueFields{Summary: "Draft draft draft"}},
		{Key: "DEV-2", Fields: jira.IssueFields{Description: "Draft notes"}},
		{Key: "DEV-3", Fields: jira.IssueFields{Summary: "Draft", Description: "Notes"}},
		{Key: "DEV-4", Fields: jira.IssueFields{Summary: "Notes", Description: "Draft"}},
	}
	previews := issuePreviews(issues, "draft draft notes", "https://example.test", "any")
	for i, key := range []string{"DEV-3", "DEV-4", "DEV-2", "DEV-1"} {
		if previews[i].Key != key {
			t.Errorf("rank %d: got %s, want %s", i+1, previews[i].Key, key)
		}
	}
}
