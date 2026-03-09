# CLAUDE.md

## Project Overview

ACLI is a Go CLI for managing Atlassian Cloud products (Jira, Confluence, Bitbucket) from the terminal. Built with [cobra](https://github.com/spf13/cobra).

## Build & Test

```bash
make build      # Build binary → bin/acli
make test       # Run all tests
make lint       # Run golangci-lint
make install    # Install to $GOPATH/bin
```

Version info is injected via ldflags at build time (see Makefile).

## Project Structure

- `main.go` — Entry point, calls `acli.Execute()`
- `cmd/acli/` — All CLI commands (cobra). One file per Atlassian product:
  - `root.go` — Root command, version command, global `--profile` flag
  - `jira.go` — Jira subcommands (issue, project, board, sprint)
  - `confluence.go` — Confluence subcommands (space, page)
  - `bitbucket.go` — Bitbucket subcommands (repo, pr, pipeline)
  - `helpers.go` — Shared utilities (e.g. `helpRunE` for group commands)
- `internal/config/` — Config loading/saving from `~/.config/acli/config.json`
- `internal/jira/`, `internal/confluence/`, `internal/bitbucket/` — API client packages (stubs, to be implemented)
- `specs/` — OpenAPI specs for Jira, Confluence, and Bitbucket Cloud APIs (reference material for implementation)

## Conventions

- Group commands (commands that only have subcommands) use `RunE: helpRunE` so they print help instead of appearing as "additional help topics"
- Each product's commands follow a consistent pattern: resource → action (e.g. `jira issue list`, `bitbucket pr get`)
- Short aliases exist for all product commands (`j`, `c`/`conf`, `bb`) and resource commands (`i`, `p`, `s`, etc.)
- Config supports multiple named profiles, selected with `--profile`/`-p` flag
- API implementations go in `internal/<product>/`, CLI wiring in `cmd/acli/`
- Many commands are currently TODO stubs — implementation should use the OpenAPI specs in `specs/` as reference

## Sensitive Files

- `.env` — Contains API tokens, **never commit**
- `~/.config/acli/config.json` — User config with API tokens at runtime
