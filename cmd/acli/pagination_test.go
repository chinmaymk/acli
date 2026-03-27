package acli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestAddAllFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addAllFlag(cmd)

	f := cmd.Flags().Lookup("all")
	if f == nil {
		t.Fatal("expected --all flag")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}

func TestAddBBPaginationFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addBBPaginationFlags(cmd)

	for _, name := range []string{"page", "pagelen", "all"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected --%s flag", name)
		}
	}
}

func TestGetBBPaginationOpts(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addBBPaginationFlags(cmd)
	cmd.Flags().Set("page", "3")
	cmd.Flags().Set("pagelen", "25")
	cmd.Flags().Set("all", "true")

	opts := getBBPaginationOpts(cmd)
	if opts.Page != 3 {
		t.Errorf("expected page 3, got %d", opts.Page)
	}
	if opts.PageLen != 25 {
		t.Errorf("expected pagelen 25, got %d", opts.PageLen)
	}
	if !opts.All {
		t.Error("expected all=true")
	}
}

func TestGetBBPaginationOptsDefaults(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addBBPaginationFlags(cmd)

	opts := getBBPaginationOpts(cmd)
	if opts.Page != 0 {
		t.Errorf("expected page 0, got %d", opts.Page)
	}
	if opts.PageLen != 0 {
		t.Errorf("expected pagelen 0, got %d", opts.PageLen)
	}
	if opts.All {
		t.Error("expected all=false")
	}
}

func TestPrintPaginationHintAllShown(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	printPaginationHint(cmd, 10, 10)
	if buf.String() != "\nShowing 10 results\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestPrintPaginationHintMoreAvailable(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	printPaginationHint(cmd, 10, 50)
	expected := "\nShowing 10 of 50 results (use --all to fetch all)\n"
	if buf.String() != expected {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestPrintPaginationHintNoTotal(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	printPaginationHint(cmd, 5, 0)
	if buf.String() != "\nShowing 5 results\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestPrintPaginationHintShownExceedsTotal(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	printPaginationHint(cmd, 15, 10)
	if buf.String() != "\nShowing 15 results\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}
