# Documentation Index

This folder documents the Go migration and the runtime behavior of the outbound
API. Keep these files synchronized with code changes.

## Documents

- [Migration to Go](migration-to-go.md)
  - Architecture, C# parity rules, implemented behavior, known gaps, build and
    test commands.
- [Configuration](configuration.md)
  - YAML configuration layout, load order, and environment overrides.
- [Authentication and Authorization](authentication.md)
  - Entra ID JWT validation and authorization API calls.
- [Documentation Sync](documentation-sync.md)
  - Source-of-truth map and checklist for keeping docs aligned with code.

## Rule

When changing behavior in `internal`, update the matching document in the same
change. If behavior is intentionally not documented, add that decision to the PR
or commit message.
