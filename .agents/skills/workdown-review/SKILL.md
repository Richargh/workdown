---
name: workdown-review
description: Reviews Workdown code changes. Use when inspecting diffs for correctness, architecture, tests, safety, CLI behavior, provider plugins, kernel contracts, configuration, filesystem, credentials, or Jira behavior.
---

# Workdown Review

STARTER_CHARACTER: 🔎

## Required context

Before a review, load and apply these skills:

- `/workspace/.agents/skills/workdown-implementation/SKILL.md`
- `/workspace/.agents/skills/workdown-testing/SKILL.md`

Use those skills as the source of truth for implementation and test expectations. Do not duplicate their rules here.

## Review focus

- Find defects that can cause wrong behavior, data loss, security problems, flaky tests, or broken CLI/API contracts.
- Check that changes keep the documented package boundaries.

## Review output

- Report findings first, ordered by severity.
- Include file path and line/range for each finding.
- Explain the user-visible impact and why the change is risky.
- If no findings exist, state that no blocking findings were found.
- Include validation commands run and their result.
