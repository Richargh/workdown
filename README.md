# Workdown

Workdown is a CLI for syncing issue tracker work items with local Markdown files. 

Has plugins to support:

* Jira

## Usage

The full CLI reference is in [docs/cli/workdown.md](docs/cli/workdown.md).

### Check Jira access

```bash
go run ./cmd/workdown remotes jira check \
  --url https://jira.example.test \
  --project PROJ \
  --pat "$JIRA_PAT"
```

### Pull Jira issues to Markdown files

```bash
go run ./cmd/workdown pull \
  --url https://jira.example.test \
  --project PROJ \
  --out items \
  --pat "$JIRA_PAT"
```

This writes files such as:

```text
items/PROJ-1.md
items/PROJ-2.md
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, build/test commands, project layout, and dependency management.
