# ACLI - Atlassian CLI

A command-line interface for managing Atlassian Cloud products — Jira, Confluence, and Bitbucket — directly from your terminal.

## Features

- **Jira** — Manage issues, projects, boards, and sprints
- **Confluence** — Manage spaces and pages
- **Bitbucket** — Manage repositories, pull requests, and pipelines
- **Multiple profiles** — Switch between different Atlassian instances easily

## Installation

### Quick install (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/chinmaymk/acli/main/install.sh | sh
```

To install a specific version or to a custom directory:

```bash
ACLI_VERSION=v1.0.0 ACLI_INSTALL_DIR=~/.local/bin curl -fsSL https://raw.githubusercontent.com/chinmaymk/acli/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
# Download the latest release for Windows
Invoke-WebRequest -Uri "https://github.com/chinmaymk/acli/releases/latest/download/acli-windows-amd64.exe" -OutFile "$env:LOCALAPPDATA\acli.exe"

# Add to PATH (current user, persists across sessions)
$path = [Environment]::GetEnvironmentVariable("Path", "User")
if ($path -notlike "*$env:LOCALAPPDATA*") {
    [Environment]::SetEnvironmentVariable("Path", "$path;$env:LOCALAPPDATA", "User")
}
```

### From source

```bash
git clone https://github.com/chinmaymk/acli.git
cd acli
make install
```

### Pre-built binaries

Download the latest release from [GitHub Releases](https://github.com/chinmaymk/acli/releases) for your platform (macOS, Linux, Windows).

## Configuration

### Quick Start

The easiest way to get set up is the interactive setup command:

```bash
acli config setup
```

This will walk you through creating a `default` profile. To create a named profile:

```bash
acli config setup work
```

### Getting Your Credentials

**Jira & Confluence** use email + API token (Basic Auth):

1. Go to [Atlassian API Tokens](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Click **Create API token**, give it a label, and copy the token
3. Your Atlassian URL looks like `https://your-instance.atlassian.net`

**Bitbucket** uses app passwords or access tokens:

1. Go to [Bitbucket App Passwords](https://bitbucket.org/account/settings/app-passwords/)
2. Click **Create app password**, select the permissions you need, and copy the password
3. You'll also need your Bitbucket username (shown at the top of your account settings)

Alternatively, Bitbucket can use workspace/repository access tokens (Bearer auth) — in that case, leave the username blank.

### Profile Management

```bash
acli config setup [name]    # Create or update a profile interactively
acli config list             # List all profiles
acli config show [name]      # Show profile details (tokens masked)
acli config delete <name>    # Delete a profile
```

### Using Profiles

Use `--profile` or `-p` to select a profile (defaults to `default`):

```bash
acli -p work jira issue list
acli -p personal bb repo list
```

### Config File

Profiles are stored in `~/.config/acli/config.json` (created automatically by `config setup`):

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

### Auth Modes

ACLI supports two authentication modes, detected automatically. The same credentials are used for Jira, Confluence, and Bitbucket:

| Mode | When | How |
|---|---|---|
| **Basic Auth** | Email is set in profile | `email:api_token` (personal API tokens) |
| **Bearer Auth** | Email is blank | `Authorization: Bearer <token>` (OAuth 2.0 / scoped tokens) |

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
