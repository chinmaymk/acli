package acli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
)

func TestIsJSONOutputDefault(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringP("output", "o", "text", "output format")

	if isJSONOutput(cmd) {
		t.Error("expected text output by default")
	}
}

func TestIsJSONOutputFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringP("output", "o", "text", "output format")
	cmd.Flags().Set("output", "json")

	if !isJSONOutput(cmd) {
		t.Error("expected JSON output when --output=json")
	}
}

func TestIsJSONOutputBoolFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringP("output", "o", "text", "output format")
	cmd.Flags().Bool("json", false, "json output")
	cmd.Flags().Set("json", "true")

	if !isJSONOutput(cmd) {
		t.Error("expected JSON output when --json flag is set")
	}
}

func TestOutputJSON(t *testing.T) {
	data := map[string]string{"key": "TEST-1", "status": "open"}
	err := outputJSON(data)
	if err != nil {
		t.Fatalf("outputJSON failed: %v", err)
	}
}

func TestOutputResultText(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringP("output", "o", "text", "output format")
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := outputResult(cmd, "create", "TEST-1", "Issue TEST-1 created", nil)
	if err != nil {
		t.Fatalf("outputResult failed: %v", err)
	}
	if buf.String() != "Issue TEST-1 created\n" {
		t.Errorf("unexpected text output: %q", buf.String())
	}
}

func TestOutputResultJSON(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringP("output", "o", "text", "output format")
	cmd.Flags().Set("output", "json")

	// outputResult with JSON writes to stdout, so we just verify no error
	err := outputResult(cmd, "delete", "TEST-1", "Issue TEST-1 deleted", map[string]string{"id": "123"})
	if err != nil {
		t.Fatalf("outputResult JSON failed: %v", err)
	}
}

func TestOutputResultStruct(t *testing.T) {
	r := OutputResult{
		Status:  "ok",
		Action:  "create",
		Key:     "TEST-1",
		Message: "Created",
		Data:    map[string]string{"url": "https://example.com"},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded OutputResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.Status != "ok" {
		t.Errorf("unexpected status: %s", decoded.Status)
	}
	if decoded.Action != "create" {
		t.Errorf("unexpected action: %s", decoded.Action)
	}
	if decoded.Key != "TEST-1" {
		t.Errorf("unexpected key: %s", decoded.Key)
	}
}
