# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Fixed release workflow dirty git state by writing generated release notes outside the repository.

## [0.0.5] - 2026-09-28

## Changed

- Improved release preparation tasks and automatic release tagging

## [0.0.3] - 2026-09-28

### Added

- Added `remotes` command for listing remote providers and hosting provider commands.
- Added Jira connectivity check command: `workdown remotes jira check`.
- Added Jira PAT validation via `--pat` and one-off Jira URL support via `--url`.
- Added Jira REST validation against `/rest/api/2/myself` and server info reporting.
- Added Jira pull command: `workdown remotes jira pull`.
- Added core Markdown pull command: `workdown pull`.
- Added pull filters for selected issue keys, assigned issues, and result limits.
- Added `mise run docs:cli` to generate CLI reference docs.

## [0.0.0] - 2026-09-23

### Added

- Initial development version marker.
