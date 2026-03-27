package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTestConfig(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	return tmpDir, func() {
		os.Setenv("HOME", origHome)
	}
}

func writeTestConfig(t *testing.T, tmpDir string, cfg *Config) {
	t.Helper()
	dir := filepath.Join(tmpDir, ".config", "acli")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestLoadNoConfigFile(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Profiles) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(cfg.Profiles))
	}
}

func TestLoadExistingConfig(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	expected := &Config{
		DefaultProfile: "test",
		Profiles: map[string]Profile{
			"test": {
				Name:         "test",
				AtlassianURL: "https://test.atlassian.net",
				Email:        "user@example.com",
				APIToken:     "tok123",
				Defaults: Defaults{
					Project:   "PROJ",
					Workspace: "myws",
					BBProject: "BBPROJ",
				},
			},
		},
	}
	writeTestConfig(t, tmpDir, expected)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultProfile != "test" {
		t.Errorf("expected default profile 'test', got %q", cfg.DefaultProfile)
	}
	p, ok := cfg.Profiles["test"]
	if !ok {
		t.Fatal("expected profile 'test' to exist")
	}
	if p.AtlassianURL != "https://test.atlassian.net" {
		t.Errorf("unexpected URL: %s", p.AtlassianURL)
	}
	if p.Email != "user@example.com" {
		t.Errorf("unexpected email: %s", p.Email)
	}
	if p.APIToken != "tok123" {
		t.Errorf("unexpected token: %s", p.APIToken)
	}
	if p.Defaults.Project != "PROJ" {
		t.Errorf("unexpected default project: %s", p.Defaults.Project)
	}
	if p.Defaults.Workspace != "myws" {
		t.Errorf("unexpected default workspace: %s", p.Defaults.Workspace)
	}
}

func TestSaveAndLoad(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	cfg := &Config{
		DefaultProfile: "prod",
		Profiles: map[string]Profile{
			"prod": {
				Name:         "prod",
				AtlassianURL: "https://prod.atlassian.net",
				Email:        "admin@example.com",
				APIToken:     "secret",
			},
		},
	}

	if err := cfg.Save(); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.DefaultProfile != "prod" {
		t.Errorf("expected default profile 'prod', got %q", loaded.DefaultProfile)
	}
	p := loaded.Profiles["prod"]
	if p.AtlassianURL != "https://prod.atlassian.net" {
		t.Errorf("unexpected URL: %s", p.AtlassianURL)
	}
}

func TestGetProfileExplicit(t *testing.T) {
	cfg := &Config{
		DefaultProfile: "default",
		Profiles: map[string]Profile{
			"default": {Name: "default", APIToken: "tok1"},
			"other":   {Name: "other", APIToken: "tok2"},
		},
	}

	p, err := cfg.GetProfile("other")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "other" {
		t.Errorf("expected profile 'other', got %q", p.Name)
	}
}

func TestGetProfileDefault(t *testing.T) {
	cfg := &Config{
		DefaultProfile: "default",
		Profiles: map[string]Profile{
			"default": {Name: "default", APIToken: "tok1"},
			"other":   {Name: "other", APIToken: "tok2"},
		},
	}

	p, err := cfg.GetProfile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "default" {
		t.Errorf("expected profile 'default', got %q", p.Name)
	}
}

func TestGetProfileSingleFallback(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"only": {Name: "only", APIToken: "tok"},
		},
	}

	p, err := cfg.GetProfile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "only" {
		t.Errorf("expected profile 'only', got %q", p.Name)
	}
}

func TestGetProfileNotFound(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"test": {Name: "test"},
		},
	}

	_, err := cfg.GetProfile("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
}

func TestGetProfileNoProfiles(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{},
	}

	_, err := cfg.GetProfile("")
	if err == nil {
		t.Fatal("expected error when no profiles exist")
	}
}

func TestGetProfileMultipleNoDefault(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {Name: "a"},
			"b": {Name: "b"},
		},
	}

	_, err := cfg.GetProfile("")
	if err == nil {
		t.Fatal("expected error when multiple profiles and no default")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	dir := filepath.Join(tmpDir, ".config", "acli")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte("{invalid"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
