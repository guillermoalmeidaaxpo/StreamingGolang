# Documentation Sync

Documentation is part of the migration deliverable. Every behavior change should
update the document that describes that behavior.

## Source-Of-Truth Map

| Behavior | Code source of truth | Documentation to update |
| --- | --- | --- |
| HTTP routes and endpoint modes | `internal/httpapi/router.go`, `internal/httpapi/routes.go`, `internal/httpapi/handlers.go` | `docs/migration-to-go.md` |
| Request schemas | `internal/httpapi/schemas/*.json`, `internal/httpapi/schema_validation.go` | `docs/migration-to-go.md` |
| Generic/lite request conversion | `internal/httpapi/generic_request.go`, `internal/httpapi/handlers.go` | `docs/migration-to-go.md` |
| CSV output and unified headers | `internal/httpapi/csv.go` | `docs/migration-to-go.md` |
| JSON and NDJSON streaming | `internal/httpapi/json_stream.go` | `docs/migration-to-go.md` |
| JWT authentication | `internal/platform/auth/auth.go` | `docs/authentication.md`, `docs/migration-to-go.md` |
| License authorization | `internal/app/authz/license.go`, `internal/httpapi/license_middleware.go` | `docs/authentication.md`, `docs/migration-to-go.md` |
| Configuration load order and env vars | `internal/platform/config/*.go`, `configs/*.yaml` | `docs/configuration.md`, `docs/migration-to-go.md` |
| ANTLR grammar and visitor behavior | `internal/query/parser/antlr/grammar/*.g4`, `internal/query/parser/antlr/ast_visitor.go`, `internal/query/parser/antlr/parser.go` | `docs/migration-to-go.md` |
| Domain commands, mappings, filters | `internal/domain/*.go` | `docs/migration-to-go.md` |
| Planning and strategy selection | `internal/app/transactional/plan.go`, `internal/app/transactional/query_strategy.go` | `docs/migration-to-go.md` |
| CMDP quote indices | `internal/app/transactional/quote_index.go` | `docs/migration-to-go.md` |
| Cassandra quote indices and RDP | `internal/app/transactional/cassandra_quote_index.go`, `internal/app/transactional/rdp_calculator.go` | `docs/migration-to-go.md` |
| Validation layer | `internal/app/transactional/validator.go`, `internal/app/transactional/schema_validator.go` | `docs/migration-to-go.md` |
| Transformations | `internal/app/transactional/transform.go` | `docs/migration-to-go.md` |
| Mapping resolution | `internal/infra/mssql/mapping_resolver.go`, `internal/infra/mssql/mds_mapping.go`, `internal/infra/mssql/mapping_row.go` | `docs/migration-to-go.md` |
| CMDP/Hyperscale SQL | `internal/infra/mssql/query_builder.go` | `docs/migration-to-go.md` |
| MSSQL execution | `internal/infra/mssql/repository.go`, `internal/infra/mssql/database.go` | `docs/migration-to-go.md`, `docs/configuration.md` |
| Cassandra execution | `internal/infra/cassandra/*.go` | `docs/migration-to-go.md`, `docs/configuration.md` |
| Build and local run commands | `go.mod`, `cmd/outbound-api/main.go` | `docs/migration-to-go.md` |

## Change Checklist

Before finishing a code change, check:

- Did an endpoint path, response format, or content negotiation rule change?
- Did a request field, schema rule, or validation rule change?
- Did parser grammar or visitor behavior change?
- Did strategy selection, split behavior, quote-index generation, or switchover
  logic change?
- Did CMDP, Hyperscale, Cassandra, or Mesap query generation change?
- Did authentication, authorization, config, or environment variables change?
- Did a known gap become implemented?

If yes, update the matching documentation in the same change.

## Verification Commands

Use these commands after documentation-related code changes:

```powershell
go build ./...
go test -c -o internal\app\transactional\transactional.test.exe ./internal/app/transactional
go test -c -o internal\infra\mssql\mssql.test.exe ./internal/infra/mssql
go test -c -o internal\httpapi\httpapi.test.exe ./internal/httpapi
go test -c -o internal\query\parser\antlr\antlrparser.test.exe ./internal/query/parser/antlr
```

On managed Windows machines, test execution may be blocked by group policy, so
`go test -c` is the default compile check.

## Documentation Review Pattern

When reviewing changes, compare:

1. The changed code files.
2. The source-of-truth map above.
3. The corresponding documentation section.

If a behavior changed but the docs did not, either update the docs or explicitly
state why no documentation change is required.
