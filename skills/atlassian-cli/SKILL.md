---
name: Atlassian CLI (acli)
description: Drive Jira, Confluence, and Bitbucket Cloud from the terminal. Use this when a user asks to read or modify Atlassian resources (issues, pages, PRs, sprints, etc.) and `acli` is installed.
---

# Atlassian CLI (acli)

`acli` is a Go CLI for Atlassian Cloud (Jira, Confluence, Bitbucket). All three products share one binary, one config, and one auth.

## When to use this skill

Use `acli` instead of `curl` against the Atlassian REST APIs whenever possible — auth, pagination, and JSON shape are already handled. Reach for the OpenAPI specs in `specs/` only when no command covers what you need.

## Discovery: ask the binary, not the docs

The command surface is large (Jira ~80 subcommands, Confluence ~150, Bitbucket ~60) and changes faster than this file. Get the authoritative list at runtime:

```bash
acli commands               # full command tree as JSON (subcommands, flags, args, aliases)
acli <any-command> --help   # human-readable help for one command
```

`acli commands` is designed for agent consumption. Pipe it through `jq` to find what you need:

```bash
acli commands | jq '.. | objects | select(.full_command? | startswith("acli jira issue"))'
```

## Invocation modes

### Normal mode

```bash
acli jira issue get PROJ-123
```

### Stdin mode (`acli -`) — prefer this for agent-driven calls

When the single argument is `-`, `acli` reads a **JSON array of strings** from stdin and uses it as argv:

```bash
echo '["jira","issue","create","--project","PROJ","--type","Task","--summary","Body with $weird `chars` and \"quotes\""]' \
  | acli -
```

Use this whenever an argument contains:

- Shell metacharacters (`$`, backticks, `;`, `&`, `|`, `>`, `<`, `*`, `?`)
- Embedded quotes or newlines (JQL strings, page bodies, PR descriptions, comments)
- Untrusted user input

It eliminates an entire class of quoting and shell-injection bugs because nothing goes through the shell parser. Empty input or anything that isn't a JSON array of strings is rejected with an error.

## Setup

```bash
which acli || (make build && make install)   # bin/acli (build) or $GOPATH/bin/acli (install)
acli config list >/dev/null 2>&1 || acli config setup <profile-name>
```

`config setup` is interactive — URL, email, API token. If a profile already exists, `acli` is ready to use.

- **Email + token** → Basic auth (Atlassian Cloud API tokens)
- **Token only, no email** → Bearer auth (OAuth / scoped tokens)

Same credentials are used for Jira, Confluence, and Bitbucket. Config lives at `~/.config/acli/config.json`.

### Profile management

```bash
acli config list                       # ls profiles, * marks default
acli config show [profile]             # show details, tokens masked
acli config set-default <profile>
acli config delete <profile>
acli config set-defaults [profile]     # set defaults for --project, workspace, bb-project
```

`set-defaults` lets you skip `--project` on every Jira command and the `<workspace>` arg on every Bitbucket command — they fall back to profile defaults when omitted.

## Global flags

| Flag         | Short | Default | Notes                                                     |
|--------------|-------|---------|-----------------------------------------------------------|
| `--profile`  | `-p`  | default | Override default profile for one call                     |
| `--output`   | `-o`  | `text`  | `text` or `json`. **Use `json` for agent/scripted use.**  |

Many commands also accept a local `--json` flag (kept for backwards compat). `-o json` is the canonical switch and works on every command, including mutations — it returns a structured envelope (`{status, action, key, message, data}`) so agents can confirm writes and pull out IDs.

**Errors are always plain-text on stderr, regardless of `-o json`.** Exit code is non-zero on failure. When scripting, check `$?` and capture stderr separately:

```bash
out=$(acli jira issue get PROJ-123 -o json 2>err.log) || { cat err.log; exit 1; }
```

## Pagination

Conventions differ by product:

- **Jira**: `--max-results`, `--start-at`, `--all` (offset-based)
- **Confluence**: `--limit`, `--cursor`, `--all` (cursor-based)
- **Bitbucket**: `--page`, `--pagelen` (max 100), `--all`

`--all` walks every page and returns the combined result. Use it when you actually need everything; otherwise prefer a bounded request.

## Command structure

```
acli <product> <resource> <action> [args] [flags]
```

Aliases: `jira`→`j`, `confluence`→`conf`/`c`, `bitbucket`→`bb`. Most resources also have aliases (`issue`→`i`, `project`→`p`, `pr`, `repo`→`r`, etc.) — check `acli commands` or `--help`.

## High-frequency examples

### Jira

```bash
# Search / read
acli jira issue list --jql "assignee = currentUser() AND statusCategory != Done" -o json
acli jira issue get PROJ-123 -o json
acli jira search --jql "project = PROJ ORDER BY updated DESC" --max-results 25

# Write (use stdin mode when bodies/JQL get hairy)
echo '["jira","issue","create","--project","PROJ","--type","Task","--summary","Fix login bug","--description","Repro: ...\nExpected: ..."]' | acli -
acli jira issue transition PROJ-123 --status "In Progress"
acli jira issue comment add PROJ-123 --body "Deployed to staging."
acli jira issue assign PROJ-123 <account-id>
acli jira issue link create PROJ-123 --title "Design doc" --url https://...

# Boards & sprints
acli jira board list --project PROJ
acli jira sprint issues <sprint-id> --all -o json
acli jira sprint move <sprint-id> PROJ-1 PROJ-2 PROJ-3
```

### Confluence

```bash
# Read
acli confluence space list -o json
acli confluence page get <page-id> --body-format storage -o json
acli confluence page list --space-id <space-id> --limit 50 -o json

# Write — page bodies almost always need stdin mode
echo '["confluence","page","create","--space-id","123","--title","Runbook","--body","<p>Step 1: <code>kubectl ...</code></p>","--body-format","storage"]' | acli -
acli confluence comment footer create --page-id <page-id> --body "LGTM"
acli confluence label pages <label-id> --all -o json

# Updates require both --title and --version-number (current version + 1)
v=$(acli confluence page get <page-id> --include-version -o json | jq -r '.version.number')
acli confluence page update <page-id> --title "New title" --body "<p>...</p>" --version-number $((v+1))
```

Confluence v2 page/blogpost/comment updates **require** `--title` and `--version-number` (current + 1). Read the current version with `--include-version`, increment, and pass it back. Body format defaults to `storage` (Confluence storage XML); pass `--body-format atlas_doc_format` or `wiki` if you have ADF/wiki markup instead.

### Bitbucket

```bash
# Read
acli bitbucket repo list --query 'name ~ "api"' -o json
acli bitbucket pr list <repo-slug> --state OPEN -o json
acli bitbucket pr diff <repo-slug> <pr-id>

# Write
acli bitbucket pr create <repo-slug> --title "Add retry logic" --source feature/retry --destination main
acli bitbucket pr approve <repo-slug> <pr-id>
acli bitbucket pr comment <repo-slug> <pr-id> --body "Nit on line 42" --file src/x.go --line 42
acli bitbucket pipeline run <repo-slug> --branch main
acli bitbucket pipeline log <workspace> <repo-slug> <pipeline-uuid> <step-uuid>
```

`<workspace>` is optional on most Bitbucket commands when a profile default is set (`acli config set-defaults`).

## Common pitfalls

- **Jira: edit, not update.** It's `acli jira issue edit <key>`, not `update`.
- **Confluence updates** need both `--title` and `--version-number` (see Confluence section). Forgetting either is the most common error.
- **Bitbucket `<workspace>` is positional**, not a flag. With a profile default it can be omitted, but otherwise: `acli bb pr list <workspace> <repo-slug>`.
- **`--all` walks every page.** Cheap on small result sets, expensive (and rate-limit-prone) on big ones. Prefer a bounded request when you can.
- **Search vs. list.** `acli jira search` is the JQL-first command with `--fields`. `acli jira issue list` is a convenience wrapper with `--project`/`--assignee`/`--status` shortcuts. Either accepts `--jql`.
- **Stdin mode rejects empty input and non-string-array JSON.** `[]` falls through to the root help; pipe in a real argv array.
- **Mutations in dry-run mode**: there is no dry-run. Read first (`get`/`list -o json`) before destructive writes.

## Tips for agent use

1. **Always pass `-o json`** — outputs are stable, parseable envelopes; text output is for humans.
2. **Use `acli -` (stdin mode)** for any argument that isn't a plain ASCII identifier. No `printf %q`, no escaping bugs.
3. **Discover, don't memorize.** If unsure about a flag or subcommand, run `acli commands | jq ...` or `acli <path> --help` rather than guessing.
4. **Mutations return identifiers.** With `-o json`, `create`/`update`/`delete` emit `{status, action, key, id, ...}` — capture `key` or `id` for follow-up calls.
5. **Profile defaults** (`config set-defaults`) eliminate repeated `--project`/workspace flags and reduce error surface.

## Sensitive files

- `.env` — API tokens. Never commit.
- `~/.config/acli/config.json` — runtime tokens. Never paste contents into a PR/issue.

## API reference (for gaps)

- `specs/jira.openapi.json`
- `specs/confluence.openapi.json` (v2)
- `specs/bitbucket.openapi.json`

Use these only when no `acli` command covers the operation — then fall back to authenticated `curl` using `acli config show` for the URL/email.
