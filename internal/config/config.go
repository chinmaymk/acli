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
	BBProject string `json:"bb_project,omitempty"`
}

type Profile struct {
	Name         string   `json:"name"`
	AtlassianURL string   `json:"atlassian_url"`
	Email        string   `json:"email"`
	APIToken     string   `json:"api_token"`
	// Optional Bitbucket-specific credentials. When empty, falls back to
	// Email / APIToken. Bitbucket Cloud uses app passwords which are
	// separate from Atlassian API tokens used for Jira/Confluence.
	BitbucketEmail string `json:"bitbucket_email,omitempty"`
	BitbucketToken string `json:"bitbucket_token,omitempty"`
	Defaults       Defaults `json:"defaults,omitempty"`
}

// BitbucketAuth returns the effective email and token for Bitbucket,
// preferring the Bitbucket-specific fields and falling back to the shared ones.
func (p Profile) BitbucketAuth() (email, token string) {
	email = p.BitbucketEmail
	if email == "" {
		email = p.Email
	}
	token = p.BitbucketToken
	if token == "" {
		token = p.APIToken
	}
	return
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
	var p Profile
	var found bool

	// Explicit profile name — look it up directly.
	if name != "" {
		p, found = c.Profiles[name]
		if !found {
			return Profile{}, fmt.Errorf("profile %q not found", name)
		}
	}

	if !found {
		// No profile specified — use the configured default.
		if c.DefaultProfile != "" {
			p, found = c.Profiles[c.DefaultProfile]
		}
	}

	if !found {
		// Fall back to the sole profile if there's exactly one.
		if len(c.Profiles) == 1 {
			for _, p = range c.Profiles {
				found = true
			}
		}
	}

	// Apply environment variable overrides. These take precedence over
	// config-file values, allowing token injection via CI, containers, etc.
	// Supported variables: ACLI_ATLASSIAN_URL, ACLI_EMAIL, ACLI_API_TOKEN.
	if v := os.Getenv("ACLI_ATLASSIAN_URL"); v != "" {
		p.AtlassianURL = v
		found = true
	}
	if v := os.Getenv("ACLI_EMAIL"); v != "" {
		p.Email = v
		found = true
	}
	if v := os.Getenv("ACLI_API_TOKEN"); v != "" {
		p.APIToken = v
		found = true
	}
	if v := os.Getenv("ACLI_BITBUCKET_EMAIL"); v != "" {
		p.BitbucketEmail = v
		found = true
	}
	if v := os.Getenv("ACLI_BITBUCKET_TOKEN"); v != "" {
		p.BitbucketToken = v
		found = true
	}

	if found {
		return p, nil
	}

	if len(c.Profiles) == 0 {
		return Profile{}, fmt.Errorf("no profiles configured, run 'acli config setup' to create one")
	}
	return Profile{}, fmt.Errorf("multiple profiles configured, specify one with -p or set a default with 'acli config set-default <name>'")
}
