package acli

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
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

func TestExtraFields(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringArray("field", nil, "")
	cmd.Flags().String("fields-json", "", "")
	_ = cmd.Flags().Set("fields-json", `{"a":1,"b":"x"}`)
	_ = cmd.Flags().Set("field", `b=y`)
	_ = cmd.Flags().Set("field", `fixVersions=[{"name":"1.0"}]`)
	_ = cmd.Flags().Set("field", `customfield_1=5`)

	got, err := extraFields(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if got["b"] != "y" {
		t.Errorf("--field should override --fields-json, got %v", got["b"])
	}
	if got["a"] != float64(1) || got["customfield_1"] != float64(5) {
		t.Errorf("unexpected numbers: %v", got)
	}
	if _, ok := got["fixVersions"].([]interface{}); !ok {
		t.Errorf("fixVersions not parsed as JSON array: %v", got["fixVersions"])
	}

	fields := map[string]interface{}{"b": "keep"}
	mergeExtraFields(fields, got)
	if fields["b"] != "keep" || fields["a"] != float64(1) {
		t.Errorf("merge must not override dedicated flags: %v", fields)
	}
}

func TestExtraFieldsInvalid(t *testing.T) {
	for _, tc := range [][2]string{{"field", "novalue"}, {"fields-json", "[1]"}, {"fields-json", "@/nonexistent"}} {
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("field", nil, "")
		cmd.Flags().String("fields-json", "", "")
		_ = cmd.Flags().Set(tc[0], tc[1])
		if _, err := extraFields(cmd); err == nil {
			t.Errorf("%v: expected error", tc)
		}
	}
}
