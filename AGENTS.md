# AGENTS.md

## Context

This is ACLI, a Go CLI tool for Atlassian Cloud (Jira, Confluence, Bitbucket). It uses cobra for command structure and supports multiple configuration profiles.

## Guidelines for AI Agents

### Adding a new command

1. Identify which product it belongs to (jira, confluence, or bitbucket)
2. Add the cobra command in `cmd/acli/<product>.go`
3. Follow existing patterns: resource group commands use `RunE: helpRunE`, leaf commands use `Run` with the actual logic
4. Add short aliases consistent with existing ones
5. API logic goes in `internal/<product>/` — keep CLI and API concerns separate

### Implementing API calls

- Reference the OpenAPI specs in `specs/` for endpoint details and request/response schemas
- Config is loaded via `config.Load()` and profiles retrieved with `config.GetProfile(name)`
- The profile provides `AtlassianURL`, `Email`, and `APIToken` for authentication (Atlassian Cloud uses Basic Auth with email + API token)

### Testing

- Run `make test` to verify changes
- Run `make build` to ensure the binary compiles
- Run `make lint` if golangci-lint is available

### Things to watch out for

- Do not commit `.env` or any file containing API tokens
- The `specs/` JSON files are large (up to 2.4MB) — avoid reading them fully unless needed for a specific endpoint
- Version variables in `root.go` (`version`, `commit`, `date`) are set via ldflags at build time — do not hardcode values
