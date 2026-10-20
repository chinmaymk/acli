package acli

import (
	"reflect"
	"testing"
)

func TestNamedRefs(t *testing.T) {
	got := namedRefs([]string{"1.0", "2.0"})
	want := []map[string]interface{}{{"name": "1.0"}, {"name": "2.0"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestIssueVersionFlags(t *testing.T) {
	for _, name := range []string{"fix-versions", "affects-versions"} {
		if jiraIssueCreateCmd.Flags().Lookup(name) == nil {
			t.Errorf("create missing --%s", name)
		}
		if jiraIssueEditCmd.Flags().Lookup(name) == nil {
			t.Errorf("edit missing --%s", name)
		}
	}
}
