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

All Atlassian products (Jira, Confluence, and Bitbucket) use email + API token (Basic Auth). ACLI requires a **scoped API token** with access to the specific Atlassian products you plan to use:

1. Go to [Atlassian API Tokens](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Click **Create scoped API token**
3. Give it a label and grant access to the products you need (Jira, Confluence, and/or Bitbucket)
4. Copy the token
5. Your Atlassian URL looks like `https://your-instance.atlassian.net`

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

### Auth Mode

ACLI authenticates using Basic Auth (`email:api_token`) with personal API tokens. Use different profiles to configure separate credentials for each product (e.g., one profile for Jira/Confluence and another for Bitbucket).

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

## API Coverage

### Jira

| Resource | Actions |
|---|---|
| **issue** | list, get, create, edit, delete, assign, transition, transitions, attach, vote, unvote, watch, unwatch, watchers, changelog, notify, createmeta, editmeta |
| issue **comment** | list, get, add, delete |
| issue **worklog** | list, add, delete |
| issue **link** (remote) | list, create, delete |
| **project** | list, get, create, update, delete, components, versions, statuses, roles, archive, restore, features |
| **board** | list, get, config, issues, backlog, sprints, epics |
| **sprint** | get, create, update, delete, issues, move |
| **epic** | get, issues, move |
| **backlog** | move |
| **search** | *(JQL search)* |
| **filter** | list, get, create, update, delete |
| **user** | list, get, search, assignable, me, create, delete |
| **group** | list, get, create, delete, members, add-user, remove-user, search |
| **dashboard** | list, get, create, update, delete, copy |
| dashboard **gadget** | list, add, update, remove |
| **component** | get, create, update, delete |
| **version** | get, create, update, delete |
| **field** | list, create, delete, trash, restore |
| **label** | *(list all labels)* |
| **issuetype** | list, get, create, update, delete |
| **priority** | list, get, create, update, delete |
| **resolution** | list, get, create, update, delete |
| **status** | list, get, categories |
| **role** | list, get, create, delete |
| **issuelink** | create, get, delete |
| **issuelinktype** | list, get, create, update, delete |
| **screen** | list, create, delete, tabs, fields |
| **workflow** | list |
| **workflowscheme** | list, get, create, update, delete |
| **permissionscheme** | list, get, create, delete |
| **notificationscheme** | list, get, create, delete |
| **issuesecurityscheme** | list, get, create, delete |
| **fieldconfig** | list, create, delete |
| **issuetypescheme** | list, create, delete |
| **projectcategory** | list, get, create, update, delete |
| **serverinfo** | *(get server info)* |
| **webhook** | list |
| **attachment** | get, delete, meta |
| **audit** | *(get audit records)* |
| **banner** | get, set |
| **configuration** | *(get global config)* |
| **permission** | mine, all |
| **task** | get, cancel |

### Confluence

| Resource | Actions |
|---|---|
| **space** | list, get, create, pages, blogposts, labels, content-labels, custom-content, operations, permissions, role-assignments, set-role-assignments |
| **page** | list, get, create, update, update-title, delete, children, direct-children, ancestors, descendants, versions, version-details, labels, attachments, footer-comments, inline-comments, custom-content, operations, likes-count, likes-users, redact |
| **blogpost** | list, get, create, update, delete, attachments, labels, footer-comments, inline-comments, custom-content, operations, versions, version-details, likes-count, likes-users, redact |
| **comment** › footer | list, get, create, update, delete, children, operations, versions, likes-count, likes-users, version-details |
| **comment** › inline | list, get, create, update, delete, children, operations, versions, likes-count, likes-users, version-details |
| **label** | list, pages, blogposts, attachments |
| **attachment** | list, get, delete, labels, comments, operations, versions, version-details |
| **task** | list, get, update |
| **custom-content** | list, get, create, update, delete, attachments, children, labels, comments, operations, versions, version-details |
| **whiteboard** | create, get, delete, ancestors, descendants, direct-children, operations, properties |
| **database** | create, get, delete, ancestors, descendants, direct-children, operations, properties |
| **folder** | create, get, delete, ancestors, descendants, direct-children, operations, properties |
| **smart-link** | create, get, delete, ancestors, descendants, direct-children, operations, properties |
| **property** | *(per content type)* list, get, create, update, delete — for page, blogpost, comment, attachment, custom-content, whiteboard, database, folder, smart-link, space |
| **space-permission** | available |
| **admin-key** | get, enable, disable |
| **data-policy** | metadata, spaces |
| **classification** | list; per type (page, blogpost, database, whiteboard): get, set, reset; space: get, set, delete |
| **user** | bulk-lookup, check-access, invite |
| **space-role** | list, get, create, update, delete, mode |
| **convert-ids** | *(convert content IDs)* |
| **app-property** | list, get, set, delete |

### Bitbucket

| Resource | Actions |
|---|---|
| **repo** | list, get, create, delete, fork, forks |
| **pr** | list, get, create, approve, unapprove, decline, merge, request-changes, comments, comment, diff |
| **pipeline** | list, get, run, stop, steps, log, variables, add-variable, delete-variable |
| **branch** | list, get, create, delete |
| **tag** | list, get, create, delete |
| **commit** | list, get, statuses, diff |
| **workspace** | list, get, members, permissions |
| **project** | list, get, create, delete |
| **webhook** | list, get, create, delete, list-workspace, create-workspace, delete-workspace |
| **environment** | list, get, create, delete |
| **deploy-key** | list, get, create, delete |
| **download** | list, delete |
| **snippet** | list, get, create, delete |
| **issue** | list, get, create, update, delete, comments |
| **search** | code |
| **deployment** | list, get |
| **branch-restriction** | list, get, create, delete |

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
