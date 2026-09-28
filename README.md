# Workdown

Workdown is a CLI for syncing issue tracker work items with local Markdown files.

Included plugins:

* Jira

## Usage

Download a released binary from [GitHub Releases](https://github.com/Richargh/workdown/releases) and
put it on your `PATH`, then:

```bash
# Check the version
workdown -v

# Check Jira access
workdown remotes jira check \
  --url https://jira.example.test \
  --project PROJ \
  --pat "$JIRA_PAT"

# Pull Jira issues
# Stores them locally as items/PROJ-*.md
workdown pull \
  --url https://jira.example.test \
  --project PROJ \
  --out items \
  --pat "$JIRA_PAT"
```

## Local usage and contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:

* development setup
* local usage
* build/test commands
* project layout
* dependency management

## Changelog

See [CHANGELOG.md](CHANGELOG.md).
