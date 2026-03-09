# ACLI - Atlassian CLI

A command-line interface for managing Atlassian Cloud products — Jira, Confluence, and Bitbucket — directly from your terminal.

## Features

- **Jira** — Manage issues, projects, boards, and sprints
- **Confluence** — Manage spaces and pages
- **Bitbucket** — Manage repositories, pull requests, and pipelines
- **Multiple profiles** — Switch between different Atlassian instances easily

## Installation

### From source

```bash
git clone https://github.com/chinmaymk/acli.git
cd acli
make install
```

### Pre-built binaries

Download the latest release from [GitHub Releases](https://github.com/chinmaymk/acli/releases) for your platform (macOS, Linux, Windows).

## Configuration

ACLI reads configuration from `~/.config/acli/config.json`. Create it with the following structure:

```json
{
  "profiles": {
    "default": {
      "name": "default",
      "atlassian_url": "https://your-instance.atlassian.net",
      "email": "you@example.com",
      "api_token": "your-api-token"
    }
  }
}
```

You can generate an API token from your [Atlassian account settings](https://id.atlassian.com/manage-profile/security/api-tokens).

Use `--profile` or `-p` to switch between profiles:

```bash
acli -p work jira issue list
```

## Usage

```bash
# Jira
acli jira issue list
acli jira issue get PROJ-123
acli jira issue create
acli jira project list
acli jira board list
acli jira sprint list

# Confluence
acli confluence space list
acli confluence page list
acli confluence page get <page-id>

# Bitbucket
acli bitbucket repo list
acli bitbucket pr list
acli bitbucket pr get <pr-id>
acli bitbucket pipeline list

# Version
acli version
```

Short aliases are available: `j` for jira, `c`/`conf` for confluence, `bb` for bitbucket.

## Development

```bash
make build      # Build for current platform → bin/acli
make test       # Run tests
make lint       # Run linter
make clean      # Remove build artifacts
make all        # Cross-compile for all platforms
```

## License

MIT
