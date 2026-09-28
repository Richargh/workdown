# Workdown

[![CI](https://github.com/Richargh/workdown/actions/workflows/ci.yml/badge.svg)](https://github.com/Richargh/workdown/actions/workflows/ci.yml)
[![GitHub Release](https://img.shields.io/github/v/release/Richargh/workdown)](https://github.com/Richargh/workdown/releases)
[![License](https://img.shields.io/github/license/Richargh/workdown)](LICENSE)

`workdown` syncs issue tracker work items to your local file system as Markdown files.

Included issue tracker plugins:

* Jira

[Usage](#usage)
• [Installation](#installation)
• [Contributing](#contributing)
• [Changelog](#changelog)

## Usage

Check the version:

```bash
workdown -v
```

Check Jira access:

```bash
workdown remotes jira check \
  --url https://jira.example.test \
  --project PROJ \
  --pat "$JIRA_PAT"
```

Pull Jira issues and store them locally as `items/PROJ-*.md`:

```bash
workdown pull \
  --url https://jira.example.test \
  --project PROJ \
  --out items \
  --pat "$JIRA_PAT"
```

## Installation

### From GitHub Releases

Download a released binary from [GitHub Releases](https://github.com/Richargh/workdown/releases) and put it on your `PATH`.

### From source

See [CONTRIBUTING.md](CONTRIBUTING.md) for local setup and source build commands.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:

* development setup
* local usage
* build/test commands
* release flow
* commit message style
* project layout
* dependency management

## Changelog

See [CHANGELOG.md](CHANGELOG.md).
