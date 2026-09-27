---
name: workdown-implementation
description: Implements Workdown features and refactors. Use when adding or changing Go code, CLI commands, provider plugins, kernel contracts, configuration, filesystem, credentials, or Jira behavior in this repository.
---

# Workdown Implementation

STARTER_CHARACTER: 🛠️

## Design Preferences

Use strong opaque domain types by default. Prefer "parse, don't validate": parse raw boundary input once into valid values, then pass those values through service code. Do not use raw strings for values with domain meaning.

Opaque domain types must guarantee correct construction: keep fields private, construct values only through parsers or constructors, and expose data only through accessors.

Example:

```text
URL(private value)
Token(private value)

parseURL(raw: string) -> URL | error
parseToken(raw: string) -> Token | error
url.string() -> string
token.secret() -> string

Before: newAPI(rawURL: string, token: string) -> API | error
After:  newAPI(baseURL: URL, token: Token) -> API | error
```
