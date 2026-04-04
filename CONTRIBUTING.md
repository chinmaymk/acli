# Contributing to ACLI

Thanks for your interest in contributing to ACLI! This guide will help you get started.

## Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [golangci-lint](https://golangci-lint.run/welcome/install-local/) (for linting)
- Make

## Getting Started

1. Fork and clone the repository:

   ```bash
   git clone https://github.com/<your-username>/acli.git
   cd acli
   ```

2. Build and run tests to make sure everything works:

   ```bash
   make build
   make test
   make lint
   ```

3. Create a branch for your changes:

   ```bash
   git checkout -b my-feature
   ```

## Project Structure

- `main.go` -- Entry point
- `cmd/acli/` -- CLI command definitions (one file per Atlassian product)
- `internal/config/` -- Configuration loading/saving
- `internal/jira/`, `internal/confluence/`, `internal/bitbucket/` -- API client packages
- `specs/` -- OpenAPI specs for reference

See [CLAUDE.md](CLAUDE.md) for a more detailed breakdown.

## Development Workflow

### Building

```bash
make build      # Build binary to bin/acli
make install    # Install to $GOPATH/bin
```

### Testing

```bash
make test       # Run all tests
```

Please add tests for any new functionality. Tests live alongside the code they test using Go's standard `_test.go` convention.

### Linting

```bash
make lint       # Run golangci-lint
```

All code must pass linting before it can be merged.

## Code Conventions

- **Command structure**: Commands follow a `product resource action` pattern (e.g., `jira issue list`, `bitbucket pr get`).
- **Aliases**: All product commands have short aliases (`j` for jira, `c`/`conf` for confluence, `bb` for bitbucket). Resource commands also have short aliases.
- **Group commands**: Commands that only have subcommands should use `RunE: helpRunE` so they print help instead of appearing as "additional help topics".
- **API clients**: API implementation goes in `internal/<product>/`, CLI wiring goes in `cmd/acli/`.
- **OpenAPI specs**: When implementing API calls, use the specs in `specs/` as reference.

## Implementing a New Command

1. Identify which product file in `cmd/acli/` the command belongs to (e.g., `jira.go`, `confluence.go`, `bitbucket.go`).
2. Add the cobra command definition following the existing patterns in that file.
3. Implement the API client logic in the corresponding `internal/<product>/` package.
4. Add tests for both the command wiring and the API client.

Many commands are currently stubs marked with TODO -- picking one of these up is a great first contribution.

## Submitting Changes

1. Make sure your code builds, tests pass, and linting is clean:

   ```bash
   make build && make test && make lint
   ```

2. Commit your changes with a clear, descriptive message.
3. Push your branch and open a pull request against `main`.
4. Describe what your PR does and why in the PR description.

## Security

- **Never commit API tokens or secrets.** The `.env` file and `~/.config/acli/config.json` contain sensitive credentials and must not be checked in.
- If you discover a security vulnerability, please report it privately rather than opening a public issue.

## License

By contributing to ACLI, you agree that your contributions will be licensed under the [MIT License](LICENSE).
