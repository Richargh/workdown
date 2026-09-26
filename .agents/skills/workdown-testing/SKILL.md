---
name: workdown-testing
description: Tests Workdown changes. Use when adding, updating, diagnosing, or reviewing tests for Go services, Cobra CLI behavior, plugin contracts, filesystem/config/credential stores, or Jira API behavior in this repository.
---

# Workdown Testing

STARTER_CHARACTER: 🧪

## Test focus

Prefer focused tests for:

- Behavior

## Rules

- Do not require real external services in default tests.
- Use fake APIs, fake HTTP clients/servers, and in-memory stores.
- Keep provider service tests independent of Cobra unless testing CLI adapters.
- Structure tests with explicit `// given`, `// when`, and `// then` comment blocks.
- Always name the system under test `testee`.

