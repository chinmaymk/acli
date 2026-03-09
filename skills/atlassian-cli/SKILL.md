# Atlassian CLI (acli)

A Go CLI for managing Atlassian Cloud products — Jira, Confluence, and Bitbucket — from the terminal.

## Setup

### 1. Build and install

```bash
make build      # Build binary → bin/acli
make install    # Install to $GOPATH/bin
```

### 2. Configure a profile

```bash
acli config setup <profile-name>
```

This creates `~/.config/acli/config.json` with your Atlassian instance URL, email, and API token. You can create multiple profiles and switch between them.

**Authentication modes:**

- **Basic Auth** — When email is provided (used for Atlassian Cloud API tokens)
- **Bearer Auth** — When email is left blank (used for OAuth/PAT tokens)

### 3. Manage profiles

```bash
acli config list                        # List all profiles
acli config show <profile>              # Show profile details (tokens masked)
acli config set-default <profile>       # Set the default profile
acli config delete <profile>            # Delete a profile
```

Use `--profile` / `-p` on any command to override the default profile:

```bash
acli --profile work jira issue list --jql "assignee = currentUser()"
```

## Command Structure

Commands follow a consistent `product → resource → action` pattern:

```
acli <product> <resource> <action> [args] [flags]
```

### Product aliases

| Product      | Aliases       |
|-------------|---------------|
| `jira`      | `j`           |
| `confluence` | `conf`, `c`  |
| `bitbucket`  | `bb`         |

## Jira

Manage issues, projects, boards, sprints, and more.

### Issues (`jira issue` / `j i`)

```bash
acli jira issue list --jql "project = PROJ AND status = Open"
acli jira issue get PROJ-123
acli jira issue create --project PROJ --type Task --summary "Fix the bug"
acli jira issue update PROJ-123 --summary "Updated summary"
acli jira issue delete PROJ-123
acli jira issue assign PROJ-123 --assignee <account-id>
acli jira issue transition PROJ-123 --transition "In Progress"
acli jira issue add-label PROJ-123 --label backend
acli jira issue link PROJ-123 PROJ-456 --type Blocks
```

### Search (`jira search` / `j s`)

```bash
acli jira search --jql "assignee = currentUser() ORDER BY updated DESC" --max-results 25
```

Supports `--jql`, `--max-results`, `--start-at`, `--fields`, and `--json`.

### Projects (`jira project` / `j p`)

```bash
acli jira project list
acli jira project get PROJ
```

### Boards (`jira board` / `j b`)

```bash
acli jira board list
acli jira board get <board-id>
```

### Sprints (`jira sprint` / `j sp`)

```bash
acli jira sprint get <sprint-id>
acli jira sprint create --board <board-id> --name "Sprint 1"
acli jira sprint update <sprint-id> --name "Sprint 1 (revised)"
acli jira sprint start <sprint-id>
acli jira sprint end <sprint-id>
acli jira sprint delete <sprint-id>
```

### Filters (`jira filter` / `j f`)

```bash
acli jira filter list
acli jira filter get <filter-id>
acli jira filter create --name "My Filter" --jql "project = PROJ"
acli jira filter delete <filter-id>
```

### Users (`jira user` / `j u`)

```bash
acli jira user get <account-id>
acli jira user search --query "jane"
acli jira user list
```

### Groups (`jira group` / `j g`)

```bash
acli jira group list
acli jira group get <group-id>
```

### Dashboards (`jira dashboard`)

```bash
acli jira dashboard list
acli jira dashboard get <dashboard-id>
```

### Administrative resources

```bash
# Roles
acli jira role list
acli jira role get <role-id>

# Components
acli jira component get <id>
acli jira component create --project PROJ --name "Backend"
acli jira component update <id> --name "New Name"
acli jira component delete <id>

# Issue link types
acli jira issuelink-type list

# Screens
acli jira screen list

# Workflows and workflow schemes
acli jira workflow list
acli jira workflow-scheme list
acli jira workflow-scheme get <id>
```

## Confluence

Manage spaces, pages, blog posts, comments, and more.

### Spaces (`confluence space` / `c s`)

```bash
acli confluence space list
acli confluence space get <space-id>
acli confluence space create --name "Engineering" --key ENG
acli confluence space update <space-id> --name "New Name"
acli confluence space delete <space-id>
```

### Pages (`confluence page` / `c p`)

```bash
acli confluence page list --space-id <space-id>
acli confluence page get <page-id>
acli confluence page create --space-id <space-id> --title "New Page" --body "<p>Content</p>"
acli confluence page update <page-id> --title "Updated Title"
acli confluence page update-title <page-id> --title "Renamed Page"
acli confluence page delete <page-id>
```

### Blog posts (`confluence blogpost` / `c blog`)

```bash
acli confluence blogpost list
acli confluence blogpost get <blogpost-id>
acli confluence blogpost create --space-id <space-id> --title "Announcement"
acli confluence blogpost update <blogpost-id>
acli confluence blogpost delete <blogpost-id>
```

### Comments (`confluence comment` / `c cm`)

```bash
# Footer comments
acli confluence comment footer list --page-id <page-id>
acli confluence comment footer get <comment-id>
acli confluence comment footer create --page-id <page-id> --body "Nice work!"

# Inline comments
acli confluence comment inline list --page-id <page-id>
acli confluence comment inline get <comment-id>
acli confluence comment inline create --page-id <page-id> --body "Suggestion here"
```

### Labels (`confluence label` / `c l`)

```bash
acli confluence label list
acli confluence label create --name "important"
acli confluence label pages --label-id <label-id>
```

### Attachments (`confluence attachment` / `c a`)

```bash
acli confluence attachment list --page-id <page-id>
acli confluence attachment get <attachment-id>
acli confluence attachment delete <attachment-id>
```

### Tasks (`confluence task` / `c t`)

```bash
acli confluence task list
acli confluence task get <task-id>
acli confluence task create
acli confluence task update <task-id>
acli confluence task delete <task-id>
```

### Other resources

```bash
# Whiteboards
acli confluence whiteboard list
acli confluence whiteboard get <id>

# Databases
acli confluence database list
acli confluence database get <id>

# Folders
acli confluence folder list
acli confluence folder get <id>

# Smart links
acli confluence smart-link list
acli confluence smart-link get <id>

# Properties
acli confluence property list --page-id <page-id>
acli confluence property get --page-id <page-id> --key <key>
acli confluence property set --page-id <page-id> --key <key> --value <value>

# Space permissions
acli confluence space-permission list --space-id <space-id>
acli confluence space-permission grant --space-id <space-id>
acli confluence space-permission revoke --space-id <space-id>
```

## Bitbucket

Manage repositories, pull requests, pipelines, branches, and more.

### Repositories (`bitbucket repo` / `bb r`)

```bash
acli bitbucket repo list <workspace>
acli bitbucket repo get <workspace> <repo-slug>
acli bitbucket repo create <workspace> --name "my-repo"
acli bitbucket repo delete <workspace> <repo-slug>
```

### Pull requests (`bitbucket pr` / `bb pr`)

```bash
acli bitbucket pr list <workspace> <repo-slug>
acli bitbucket pr get <workspace> <repo-slug> <pr-id>
acli bitbucket pr create <workspace> <repo-slug> --title "Feature" --source feature-branch
acli bitbucket pr update <workspace> <repo-slug> <pr-id> --title "Updated"
acli bitbucket pr approve <workspace> <repo-slug> <pr-id>
acli bitbucket pr decline <workspace> <repo-slug> <pr-id>
acli bitbucket pr request-changes <workspace> <repo-slug> <pr-id>
acli bitbucket pr comment <workspace> <repo-slug> <pr-id> --body "LGTM"
acli bitbucket pr comments <workspace> <repo-slug> <pr-id>
```

### Pipelines (`bitbucket pipeline` / `bb pipe`)

```bash
acli bitbucket pipeline list <workspace> <repo-slug>
acli bitbucket pipeline get <workspace> <repo-slug> <pipeline-uuid>
acli bitbucket pipeline trigger <workspace> <repo-slug> --branch main
acli bitbucket pipeline stop <workspace> <repo-slug> <pipeline-uuid>
acli bitbucket pipeline add-variable <workspace> <repo-slug> --key VAR --value val
```

### Branches (`bitbucket branch`)

```bash
acli bitbucket branch list <workspace> <repo-slug>
acli bitbucket branch get <workspace> <repo-slug> <branch-name>
acli bitbucket branch create <workspace> <repo-slug> --name feature-x --target main
acli bitbucket branch delete <workspace> <repo-slug> <branch-name>
```

### Commits (`bitbucket commit`)

```bash
acli bitbucket commit list <workspace> <repo-slug>
acli bitbucket commit get <workspace> <repo-slug> <commit-hash>
acli bitbucket commit comment <workspace> <repo-slug> <commit-hash> --body "Note"
```

### Tags (`bitbucket tag`)

```bash
acli bitbucket tag list <workspace> <repo-slug>
acli bitbucket tag get <workspace> <repo-slug> <tag-name>
acli bitbucket tag create <workspace> <repo-slug> --name v1.0.0 --target main
acli bitbucket tag delete <workspace> <repo-slug> <tag-name>
```

### Workspaces (`bitbucket workspace`)

```bash
acli bitbucket workspace list
acli bitbucket workspace get <workspace>
```

### Projects (`bitbucket project`)

```bash
acli bitbucket project list <workspace>
acli bitbucket project get <workspace> <project-key>
acli bitbucket project create <workspace> --name "My Project" --key PROJ
acli bitbucket project delete <workspace> <project-key>
```

### Other resources

```bash
# Webhooks
acli bitbucket webhook list <workspace> <repo-slug>
acli bitbucket webhook create <workspace> <repo-slug> --url https://example.com/hook

# Environments
acli bitbucket environment list <workspace> <repo-slug>
acli bitbucket environment create <workspace> <repo-slug> --name production

# Deploy keys
acli bitbucket deploy-key list <workspace> <repo-slug>
acli bitbucket deploy-key create <workspace> <repo-slug> --key "ssh-rsa ..." --label "CI"

# Deployments
acli bitbucket deployment list <workspace> <repo-slug>

# Snippets
acli bitbucket snippet list <workspace>
acli bitbucket snippet get <workspace> <snippet-id>
acli bitbucket snippet create <workspace> --title "Snippet" --filename code.py

# Downloads
acli bitbucket download get <workspace> <repo-slug> <filename>

# Issues (requires issue tracker enabled on repo)
acli bitbucket issue list <workspace> <repo-slug>
acli bitbucket issue create <workspace> <repo-slug> --title "Bug"

# Search
acli bitbucket search <workspace> <repo-slug> --query "TODO"

# Branch restrictions
acli bitbucket branch-restriction list <workspace> <repo-slug>
```

## Common flags

| Flag              | Short | Description                         |
|-------------------|-------|-------------------------------------|
| `--profile`       | `-p`  | Select configuration profile        |
| `--json`          |       | Output raw JSON response            |
| `--max-results`   |       | Limit number of results returned    |
| `--start-at`      |       | Pagination offset                   |

## API Reference Specs

OpenAPI specifications for each Atlassian product are available in the `specs/` directory for detailed endpoint and field documentation:

- `specs/jira.openapi.json` — Jira Cloud REST API
- `specs/confluence.openapi.json` — Confluence Cloud REST API (v2)
- `specs/bitbucket.openapi.json` — Bitbucket Cloud REST API
