package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Defaults struct {
	Project   string `json:"project,omitempty"`
	Workspace string `json:"workspace,omitempty"`
}

type Profile struct {
	Name         string   `json:"name"`
	AtlassianURL string   `json:"atlassian_url"`
	Email        string   `json:"email"`
	APIToken     string   `json:"api_token"`
	Defaults     Defaults `json:"defaults,omitempty"`
}

type Config struct {
	DefaultProfile string             `json:"default_profile,omitempty"`
	Profiles       map[string]Profile `json:"profiles"`
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "acli"), nil
}

func Load() (*Config, error) {
	dir, err := ConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Profiles: make(map[string]Profile)}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0600)
}

func (c *Config) GetProfile(name string) (Profile, error) {
	// Explicit profile name — look it up directly.
	if name != "" {
		p, ok := c.Profiles[name]
		if !ok {
			return Profile{}, fmt.Errorf("profile %q not found", name)
		}
		return p, nil
	}
	// No profile specified — use the configured default.
	if c.DefaultProfile != "" {
		if p, ok := c.Profiles[c.DefaultProfile]; ok {
			return p, nil
		}
	}
	// Fall back to the sole profile if there's exactly one.
	if len(c.Profiles) == 1 {
		for _, p := range c.Profiles {
			return p, nil
		}
	}
	if len(c.Profiles) == 0 {
		return Profile{}, fmt.Errorf("no profiles configured, run 'acli config setup' to create one")
	}
	return Profile{}, fmt.Errorf("multiple profiles configured, specify one with -p or set a default with 'acli config set-default <name>'")
}
