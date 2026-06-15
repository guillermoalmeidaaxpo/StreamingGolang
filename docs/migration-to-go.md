# Outbound API Migration To Go

This document captures the migration of the C# transactional outbound API to Go.
The goal is not a line-by-line rewrite. The Go service should keep the observable
API behavior and business rules of the C# implementation while using Go-native
architecture, configuration, error handling, concurrency, and dependency wiring.

## Scope

The migration targets the API currently represented by the C# solution under:

```text
C:\Projects\Streaming\src
```

The Go implementation lives in:

```text
C:\Projects\StreamingGolang
```

The Go module targets:

```text
go 1.26.0
toolchain go1.26.4
```

## Main Goals

- Preserve the C# API contract for transactional, generic, CSV, streaming JSON,
  streaming NDJSON, and lite endpoints.
- Keep the C# business behavior for mapping resolution, data-source selection,
  split strategies, quote-index generation, filters, transformations, CSV header
  unification, and authorization.
- Use Go-native practices:
  - `net/http` instead of a third-party router.
  - YAML configuration with environment overrides instead of `appsettings.json`
    cloning.
  - Small packages with explicit interfaces at application boundaries.
  - Context-aware repositories and streaming iterators.
  - Structured logging with `log/slog`.
  - JSON Schema validation before decoding request bodies.

## Architecture

The project is organized around Go's `internal` package boundary.

```text
cmd/outbound-api
  Application entry point.

internal/httpapi
  HTTP routing, handlers, middleware, JSON Schema validation,
  CSV/JSON/NDJSON response writers, lite/generic request conversion.

internal/app/transactional
  Use-case orchestration: request planning, validation, strategy selection,
  split query planning, quote-index planning, execution, streaming pipeline,
  transformations.

internal/app/authz
  License/MDO authorization model and external authorization client contract.

internal/domain
  Shared domain types: commands, mappings, filters, executable queries,
  source kinds, identifiers.

internal/query/parser/antlr
  ANTLR-based filter parser and visitor that converts parse trees into the Go
  filter AST.

internal/infra/mssql
  MSSQL/Azure SQL adapters: datastore setup, mapping resolver,
  CMDP/Hyperscale query builders, executor.

internal/infra/cassandra
  Cassandra datastore, query builder, and repository.

internal/infra/redis
  Redis client setup.

internal/platform/config
  Go-native YAML configuration loader and environment overrides.

internal/platform/auth
  Entra ID/OIDC JWT authentication middleware.

internal/platform/server
  HTTP server wrapper.
```

## Configuration

The C# app uses `appsettings.json` and environment-specific overlays such as
`appsettings.Development.json`.

The Go service uses a Go-native configuration model:

1. `configs/default.yaml`
2. `configs/<environment>.yaml`
3. `OUTBOUND_*` environment variables

The default environment is `development`.

Useful overrides:

```powershell
$env:OUTBOUND_ENV = "development"
$env:OUTBOUND_CONFIG_DIR = "C:\Projects\StreamingGolang\configs"
$env:OUTBOUND_HTTP_PORT = "8080"
$env:OUTBOUND_AUTH_MODE = "jwt"
$env:OUTBOUND_AUTHORIZATION_API_BASE_URL = "https://authorization-hs.lab.mds.axpo.com"
```

Connection strings and secrets should live in environment variables or the
deployment secret store. For local development, the development YAML can carry
non-secret defaults, but credentials should not be committed.

More detail is in:

```text
docs/configuration.md
```

The source-of-truth map for keeping this document aligned with code is in:

```text
docs/documentation-sync.md
```

## Authentication And Authorization

The C# API has two distinct security steps and the Go API follows the same model:

1. JWT bearer authentication with Entra ID.
2. MDO/license authorization against the authorization API.

JWT authentication validates issuer, audience, roles/scopes, expiration, and
signature using OIDC metadata.

License authorization calls the authorization service before handlers execute.
The C# flow calls two authorization endpoints:

```text
POST /api/v1/DataUniverse/BulkAuthorize
POST /api/v1/TimeSeries/Authorize
```

The Go payloads follow the same shape:

```json
{
  "type": "TransactionalDataOutbound",
  "universeName": null,
  "action": "Read",
  "internalCorrelationId": "<correlation-id>",
  "mdoIds": [312091001],
  "stageId": 3
}
```

```json
{
  "identifiers": [312091001],
  "stageId": 3,
  "internalCorrelationId": "<correlation-id>"
}
```

The incoming bearer token is forwarded to the authorization API.

## Endpoints

Implemented endpoint families:

- `/api/v1/timeseries`
- `/api/v1/timeseries/streaming`
- `/api/v1/curves`
- `/api/v1/curves/streaming`
- `/api/v1/surfaces`
- `/api/v1/surfaces/streaming`
- `/api/v1/generic`
- `/api/v1/generic/streaming`
- `/api/v1/lite`
- `/api/v1/info`
- `/health`

Endpoint behavior:

- Transactional endpoints support JSON responses.
- Transactional streaming endpoints support JSON array streaming and NDJSON
  content negotiation.
- Generic endpoints match the C# behavior and return CSV, even if the client
  sends JSON or NDJSON `Accept` headers.
- Lite endpoint builds CSV requests from query parameters.

## Request Validation

The Go API validates raw request bodies using real JSON Schema files before
decoding into Go structs:

```text
internal/httpapi/schemas/transactional.schema.json
internal/httpapi/schemas/generic.schema.json
internal/httpapi/schemas/lite.schema.json
```

Important C#-compatible schema behavior:

- `filters` may be `null`.
- `transformations` may be `null`.
- `transformations.nested` may be `null`.
- Unknown fields are rejected.
- Invalid schema failures return RFC7807-style problem responses.

After schema validation, application validation still checks business rules that
cannot be expressed cleanly in JSON Schema.

## ANTLR Filter Parser

The C# app uses ANTLR, so the Go migration also uses ANTLR.

Grammar files:

```text
internal/query/parser/antlr/grammar/OutboundAPILexer.g4
internal/query/parser/antlr/grammar/OutboundAPIParser.g4
```

Generated Go parser files:

```text
internal/query/parser/antlr/generated
```

Regeneration:

```powershell
go generate ./internal/query/parser/antlr
```

The Go visitor converts ANTLR parse trees into the domain filter AST. It handles:

- scalar comparisons
- numeric comparisons
- text comparisons
- point-in-time values
- point-in-time arithmetic such as `now()+P1D`
- explicit intervals such as `ti(start,end)`
- interval functions such as day/week/month/quarter/year
- gas interval functions
- `begin(...)` and `end(...)`
- `latest(...)`
- `latestGlobal()`
- rank-over filter AST nodes

Important parser parity rules:

- `latest()` and `latestGlobal()` are only valid for `ReferenceTime`.
- `latest()` and `latestGlobal()` must use equality.
- `latest()` accepts one reference-time argument.
- `ReferenceTime in ti(...)` is expanded into lower/upper bound filters.

## Mapping Resolution

The Go mapping resolver reads both CMDP mappings and MDS/Hyperscale mappings.

CMDP mapping source:

```sql
CMDP_TO_MDS_MAPPING
```

MDS/Hyperscale mapping source:

```sql
[Api].[VI_MdsMappingDetails]
```

Mapped fields include:

- MDO/timeseries id
- data category
- CMDP view
- MDS structure
- resolution
- CMDP/MDS column names
- datatype
- key/projectable flags
- ordering metadata
- Cassandra id
- Hyperscale id
- split-query flag
- indexed column information
- switchover metadata
- Hyperscale latest/get-by-created-on views
- mapping timezone

Generic endpoint behavior:

- Generic requests are not forced into one data category.
- Mappings are grouped by data category and planned as separate commands.
- This matches the C# generic behavior where one generic request can retrieve
  different data categories.

## Data-Source Selection

The Go selector follows the C# `MdoDataFetchingStrategyParser` priority:

1. Mesap endpoint and Mesap id -> Mesap strategy.
2. Mapping has `HyperScaleId` -> Hyperscale strategy.
3. Use CMDP if any of these are true:
   - shape filters are present
   - aggregations are present
   - id is in the HPFC force-to-CMDP list
   - Cassandra id is missing
   - Cassandra timezone is not Europe/Zurich/CET
4. Otherwise use Cassandra.

Current HPFC force-to-CMDP ids in Go:

```text
536000751
536214287
536346251
```

## Split Strategy

The C# app has `SingleQueryStrategy` and `SplitQueryStrategy`.

The Go planner keeps the same idea:

- single-query planning by default
- split planning where the command/mapping supports split query
- special hybrid CMDP/Cassandra split based on reference-time filters and the
  CMDP watermark

Hybrid behavior:

- If a Cassandra mapping is split-enabled and has reference-time filters, Go
  reads the CMDP watermark.
- Equality filters route to Cassandra when the reference time is before the
  watermark, otherwise CMDP.
- Range filters route entirely to Cassandra, entirely to CMDP, or split into two
  commands around the watermark.

Shape and aggregation requests do not hybrid-split because C# forces them to CMDP.

## Quote Indices

CMDP and Cassandra quote-index behavior are separate, as in C#.

CMDP quote-index logic:

- Derived from reference-time filters for indexed CMDP queries.
- Supports split ranges for query planning.

Cassandra quote-index logic:

- Follows the Cassandra-specific rules from C# services such as
  `QuoteIndexGenerator`, `CassandraFilterProcessor`, and
  `CassandraFilterProcessingStrategy`.
- A midnight `ReferenceTime` in the filter timezone generates a quote index.
- Example:

```json
{
  "filters": {
    "expressions": ["ReferenceTime = 2024-04-26T00:00:00"],
    "filterTimeZone": "Europe/Zurich"
  },
  "ids": [536013751]
}
```

Generates:

```text
20240426
```

But this should not generate a quote index because it is not midnight in the
filter timezone:

```json
{
  "filters": {
    "expressions": ["ReferenceTime = 2024-04-26T22:00:00"],
    "filterTimeZone": "Europe/Zurich"
  },
  "ids": [536013751]
}
```

Delivery/RDP filters are part of the Cassandra filter processing flow and follow
the C# `FiltersExtensions.GenerateDeliveryFilters` behavior:

- `DeliveryStart` and `DeliveryEnd` filters are converted to the Cassandra-local
  delivery hour tuple `(del_y, del_m, del_d, del_h)`.
- `DeliveryEnd` is shifted back by one hour for Cassandra CQL filtering, matching
  the current C# Cassandra filter implementation.
- `RelativeDeliveryPeriod` filters are integer hours from the quote-index local
  midnight, using the first quote index in the Cassandra batch.
- Delivery and RDP windows are intersected before generating CQL.
- An empty delivery/RDP intersection skips the Cassandra query for that batch.

The separate RDP calculator remains used for response enrichment semantics; it
is not the same path as Cassandra filter narrowing.

## Hyperscale Query Rules

The Go Hyperscale builder follows the C# `DynamicQueryBuilder` and `ViewProvider`
rules.

Important C# parity rules:

- If there are real filters and no `LatestGlobal()`, use latest-version views.
- If there are no real filters, use latest-reference-time views.
- If `LatestGlobal()` is present, use latest-reference-time views.
- If `VersionAsOf` is present and latest-reference-time is required, use
  `TVF_Get<DataCategory>ByCreatedOnLatestReferenceTime`.
- If `VersionAsOf` is present and real filters are present, use
  `TVF_Get<DataCategory>ByCreatedOn`.
- `latest(...)` on Hyperscale follows the C# `LatestReference` CTE pattern:
  - normal requests query the latest-version view in the CTE and the main query
  - `VersionAsOf` requests query `Core.<DataCategory>Version` in the CTE and
    `TVF_Get<DataCategory>ByCreatedOn` in the main query
  - the inner latest argument becomes the CTE reference-time predicate
  - the main query adds `ReferenceTime = (SELECT MaxReferenceTimeBefore FROM
    LatestReference)`
- For timeseries, `LatestVersionWithCreatedOnView` and
  `LatestReferenceTimeWithCreatedOnView` fall back to the normal timeseries
  view names. This mirrors the C# `MdoMapping` property behavior.
- `CreatedOn` is appended separately. It is not treated as a mapped value column.
- Value columns are emitted as JSON_VALUE expressions against:
  - `CurveValue`
  - `SurfaceValue`
  - `TimeSeriesValue`
- Latest-reference-time timeseries values are projected as `Property0`,
  `Property1`, etc.
- Datatype casts follow mapping metadata.
- `ORDER BY` follows C# `MappingSet.GetOrderByMdsColumns`: mapping
  `OrderPriority` first, otherwise `ReferenceTime`.

Example C#-compatible request:

```json
[
  {
    "ids": [504078501],
    "filters": {
      "filterTimeZone": "CET",
      "expressions": ["ReferenceTime = LatestGlobal()"]
    },
    "transformations": {
      "targetTimeZone": "CET",
      "nested": null
    },
    "columns": ["CreatedOn"]
  }
]
```

Expected Hyperscale shape:

```sql
SELECT ReferenceTime,
       JSON_VALUE(TimeSeriesValue, '$."Value"') AS Property0,
       CreatedOn
FROM [Api].[VI_TimeseriesLatestVersionLatestReferenceTime]
WHERE Deleted = '0'
  AND MdoId = '504078501'
ORDER BY ReferenceTime
```

The Go query is parameterized and uses `[d]` aliases, but should preserve the
same behavior, selected view, selected columns, filters, and ordering.

## CMDP Query Rules

CMDP query generation uses mappings plus filters plus optional split/index
ranges.

Current behavior:

- identifier predicate uses `TimeSeries_FID`
- mapped key and projectable columns are selected
- requested columns restrict value columns
- reference-time interval filters become SQL bounds
- quote-index split range adds indexed column predicates
- rank-over filters are executed with a CMDP-only derived query using
  `RANK() OVER (PARTITION BY ... ORDER BY ...)`
- rank-over filtering follows C# bound semantics:
  - no third argument -> `rank = 1`
  - single value -> `rank = value`
  - `[n,last]` -> `rank >= n`
  - `[n,m]` -> `rank >= n AND rank <= m`
- order columns come from mapping priority when present
- shape filters are normalized and emitted as CMDP-only predicates against
  `DeliveryStart`, matching C# `ShapeNormalizer` and
  `SqlShapePredicateBuilder`
- shape month filters use `DATEPART(MONTH, DeliveryStart) IN (...)`
- shape day filters use the C# ISO weekday expression based on
  `DATEDIFF(DAY, '19000101', DeliveryStart)`
- shape time filters use `CAST(DeliveryStart AS time)` half-open ranges
- shape delivery-start expressions use `filterTimeZone`; non-UTC zones are
  translated to SQL Server `AT TIME ZONE`

Rank-over validation follows the C# API:

- not allowed with aggregations
- not allowed for timeseries
- not allowed for Cassandra, Hyperscale, or Mesap-hosted ids
- partition columns must be mapped key columns
- `Identifier/MdoId` and `RelativeDeliveryPeriod` are not allowed in rank-over
  partition/order columns

CMDP still needs careful parity testing for:

- `latest(...)` CTE behavior
- delivery/RDP filter combinations

## Cassandra Query Rules

Cassandra query generation uses Cassandra mappings, table mappings, quote
indices, and processed filter constraints.

Current behavior:

- builds Cassandra table queries from configured table mappings
- requires quote indices when Cassandra mapping/filter logic needs them
- logs Cassandra query execution failures with statement and context
- streams rows through the repository pipeline
- applies C#-compatible delivery/RDP tuple filters and skips empty
  delivery/RDP intersections

High-risk parity areas:

- Cassandra quote-index generation
- split/hybrid behavior around CMDP watermark
- Cassandra vs CMDP query shape for old reference-time data

## CSV Behavior

The C# `TransactionalDataCsvService` first computes combined headers:

```csharp
commandParser.GetCombinedHeaders(commandList)
```

The Go CSV writer preserves the same concept:

- headers are unified across all commands in the response
- headers are based on mappings and requested projection columns
- CSV generic endpoints ignore JSON/NDJSON accept negotiation and return CSV
- CSV streaming writes headers once and streams rows/batches

Important C# projection parity:

- no requested columns -> all mapped columns
- requested `CreatedOn` only -> all mapped columns plus `CreatedOn`
- custom requested columns -> non-projectable key columns plus matching requested
  projectable columns
- JSON endpoints drop `Identifier/MdoId`
- CSV endpoints keep `Identifier/MdoId`
- Aggregation CSV headers are generated from aggregation output columns rather
  than raw mapping columns. This preserves aliases such as `DeliveryBucket` and
  `AveragePrice`.
- Non-stream CSV responses now emit the unified header even when the query
  returns zero rows, so empty aggregation results are still inspectable.

## Transformations

Current Go transformation support includes:

- target timezone field handling
- offset inclusion rules
- nullable `transformations`
- nullable `nested`
- default include-offset behavior based on endpoint mode
- IANA timezone names such as `Europe/Zurich` are supported in Windows release
  builds by embedding Go timezone data; C# aliases such as `CET` are normalized
  through the shared timezone loader

C# behavior to preserve:

- Generic CSV defaults `Offset=false`.
- Generic CSV defaults target timezone to `UTC` when there is no aggregation,
  no explicit target timezone, and offset is false.
- Aggregations default target timezone to UTC when no target timezone is given.
- Aggregation keys normalize `Delivery` to `DeliveryStart`, matching the C#
  `AggregationExpressionHelper`.
- JSON and streaming JSON include offsets by default.

Remaining high-risk transformation areas:

- nested output shape
- additional aggregation period variants beyond the currently covered SQL
  bucket cases
- shape-aware transformations

## Shape Filters

Shape filters are now carried as normalized domain data:

```text
domain.NormalizedShape
  Months          []int
  Days            []int
  TimeSpans       []ShapeTimeSpan
  HolidayCalendar *int
```

The Go normalizer follows the C# `ShapeNormalizer` behavior:

- valid months are `Jan` through `Dec`, normalized to `1..12`
- valid days are `Mon` through `Sun`, normalized to `1..7`
- values are sorted after normalization
- duplicate month/day entries are rejected
- duplicate or overlapping time ranges are rejected
- `T00:00:00` used as a time-range end means end-of-day
- a full-day range emits no SQL time predicate
- `holidayCalendar` is preserved but has no SQL behavior yet, matching the C#
  `NormalizedShape` comment

Shape validation and strategy behavior:

- shape is only allowed for curves
- shape is rejected for Hyperscale ids
- shape forces CMDP strategy and disables hybrid CMDP/Cassandra split

CMDP SQL generation adds active predicates to the normal `WHERE` clause:

```sql
DATEPART(MONTH, <DeliveryStart>) IN (...)
((DATEDIFF(DAY, '19000101', <DeliveryStart>) % 7) + 1) IN (...)
CAST(<DeliveryStart> AS time) >= ...
```

When `filters.filterTimeZone` is non-UTC, `<DeliveryStart>` is wrapped with SQL
Server `AT TIME ZONE`, as in the C# implementation.

## Aggregations

Aggregation requests are now represented explicitly in the Go domain command:

```text
domain.Aggregations
  GroupBy     []AggregationColumn
  Expressions []AggregationColumn
```

The planner builds this model from `transformations.keys` and
`transformations.values`, following C# behavior:

- both `keys` and `values` must be supplied together
- only one aggregation key is allowed
- `Aggregate(Delivery, <period>)` is normalized to
  `Aggregate(DeliveryStart, <period>)`
- when no target timezone is supplied, aggregation bucketing defaults to `UTC`
- aggregation projection columns are generated like C#
  `AggregationSqlBuilder.GetAggregatedColumnNames`

Generated aggregation output columns start with:

```text
Identifier
ReferenceTime
DeliveryStart
DeliveryEnd
RelativeDeliveryPeriod
LegacyDeliveryBucketNumber   # CMDP only
```

Then group aliases and value aliases are appended, respecting requested
projection columns.

CMDP aggregation SQL now uses the C# shape:

- literal MDO id as `Identifier`
- `MIN(ReferenceTime)` unless reference time itself is bucketed
- `MIN(DeliveryStart)`
- `MAX(DeliveryEnd)`
- `NULL AS RelativeDeliveryPeriod`
- `NULL AS LegacyDeliveryBucketNumber`
- SQL Server `DATEADD/DATEDIFF` bucket expressions for aggregate keys
- aggregate values such as `AVG(...)`, `SUM(...)`, `MIN(...)`, `MAX(...)`,
  and `COUNT(...)`
- `GROUP BY` and `ORDER BY` include reference time and group expressions,
  following `AggregationSqlBuilder`

Hyperscale aggregation SQL uses the same output shape but:

- excludes `LegacyDeliveryBucketNumber`
- reads values from JSON payload columns such as `CurveValue`
- casts JSON values based on mapping datatype before aggregation
- keeps the regular Hyperscale latest-version / TVF source selection rules
- can carry the same Hyperscale `latest(...)` CTE predicate as standard
  Hyperscale queries, although live C# parity still needs proof for aggregation
  plus latest combinations

Validation now rejects:

- aggregations outside curves
- aggregation value columns that are unmapped
- aggregation value columns that are key columns
- duplicate aggregation aliases
- rank-over filters combined with aggregations

Regression coverage currently compiles tests for:

- planner aggregation metadata and generated columns
- CMDP aggregation query shape
- Hyperscale aggregation query shape
- CSV aggregation header selection

## Error Handling And Logging

The Go API uses RFC7807-style problem responses:

```json
{
  "type": "about:blank",
  "title": "invalid-request-body",
  "status": 400,
  "detail": "...",
  "instance": "/api/v1/...",
  "correlationId": "..."
}
```

Structured logs include:

- configuration loading
- datastore configuration and auth mode
- request completion
- mapping resolution
- mapping SQL queries
- authorization API calls and status
- selected strategy/source
- watermark reads
- quote-index generation
- MSSQL/Cassandra query execution and parameters
- query execution failures with duration and context

In development/local mode, wrapped application error causes can be exposed in
problem details to speed up debugging. Productive mode should keep external
details safer.

## Build And Run

Build all packages:

```powershell
go build ./...
```

Build the executable:

```powershell
go build -o bin\outbound-api.exe .\cmd\outbound-api
```

Run:

```powershell
.\bin\outbound-api.exe
```

Generate ANTLR parser:

```powershell
go generate ./internal/query/parser/antlr
```

Compile tests without executing them:

```powershell
go test -c -o internal\app\transactional\transactional.test.exe ./internal/app/transactional
go test -c -o internal\infra\mssql\mssql.test.exe ./internal/infra/mssql
go test -c -o internal\httpapi\httpapi.test.exe ./internal/httpapi
go test -c -o internal\query\parser\antlr\antlrparser.test.exe ./internal/query/parser/antlr
```

On some managed Windows machines, executing Go test binaries from `%TEMP%` may
be blocked by group policy. In that case, `go test -c` is useful as a compile
verification step.

## Useful PowerShell Test Request Pattern

For PowerShell, avoid inline JSON quoting problems by using a file:

```powershell
@'
[
  {
    "ids": [504078501],
    "filters": {
      "filterTimeZone": "CET",
      "expressions": ["ReferenceTime = LatestGlobal()"]
    },
    "transformations": {
      "targetTimeZone": "CET",
      "nested": null
    },
    "columns": ["CreatedOn"]
  }
]
'@ | Set-Content -Encoding utf8 body.json

curl.exe -v http://localhost:8080/api/v1/timeseries/streaming `
  -H "Authorization: Bearer <token>" `
  -H "Content-Type: application/json" `
  --data-binary "@body.json"
```

## Current Implementation Status

Implemented:

- Go 1.26 module and service entrypoint.
- `net/http` router.
- YAML configuration and environment overrides.
- Entra ID JWT middleware.
- Authorization/license middleware with DataUniverse and TimeSeries calls.
- JSON Schema validation.
- ANTLR parser and visitor-backed filter AST.
- Mapping resolver for CMDP and MDS/Hyperscale mappings.
- Strategy selection matching C# priority.
- Split and hybrid planning.
- CMDP quote-index planning.
- Cassandra quote-index planning.
- CMDP SQL builder.
- Hyperscale SQL builder with latest-reference-time, `latest(...)` CTE, and
  projection parity rules.
- Aggregation command model, validation, CMDP SQL, Hyperscale SQL, and CSV
  header generation.
- Shape filter normalization, validation, CMDP strategy forcing, and CMDP SQL
  predicate generation.
- Cassandra query builder and repository.
- CSV/generic/lite endpoints.
- JSON, streaming JSON, and NDJSON transactional endpoints.
- MSSQL/Azure SQL datastore setup.
- Cassandra and Redis infrastructure setup.
- Structured logging.
- Development-safe error detail exposure.

Known gaps / still needs parity work:

- Aggregation end-to-end parity still needs live CMDP/Hyperscale proof with
  representative ids and expected C# result sets.
- Shape filter end-to-end parity still needs live CMDP proof with representative
  ids and expected C# result sets.
- CMDP `latest(...)` CTE parity.
- Hyperscale `latest(...)` end-to-end parity still needs live proof with
  representative ids and expected C# result sets.
- Mesap endpoint/data-source parity.
- Data trace/migration endpoints if required by production consumers.
- Application Insights exporter integration if required.
- End-to-end contract tests against a realistic dev environment.

When any item moves between "implemented" and "known gaps", update this section
in the same change as the code.

## Migration Principles

- Do not copy C# implementation details blindly.
- Preserve external behavior and business rules.
- Prefer Go-native structure and dependencies.
- Keep domain logic in `internal/app` and `internal/domain`.
- Keep SQL/Cassandra/Redis/Auth adapters in `internal/infra` or `internal/platform`.
- Add tests around every discovered C# behavior quirk.
- Use logs to make data-source selection and query generation observable.

## C# Files Used As Behavioral References

Key C# files consulted during migration:

```text
MDS.DPS.TransactionalData.Outbound.Antlr\OutboundAPILexer.g4
MDS.DPS.TransactionalData.Outbound.Antlr\OutboundAPIParser.g4
MDS.DPS.TransactionalData.Outbound.Domain\Expression\Parsers\FilterExpressionParser.cs
MDS.DPS.TransactionalData.Outbound.Domain\Expression\Parsers\FilterExpressionVisitor.cs
MDS.DPS.TransactionalData.Outbound.Domain\Filtering\FilterProvider.cs
MDS.DPS.TransactionalData.Outbound.Domain\Extensions\FilterSetExtensions.cs
MDS.DPS.TransactionalData.Outbound.Domain\Extensions\CmdpQuoteIndexFilterExtensions.cs
MDS.DPS.TransactionalData.Outbound.Domain\Extensions\RelativeDeliveryPeriodFilterBuilder.cs
MDS.DPS.TransactionalData.Outbound.Domain\Extensions\RankOverClauseBuilder.cs
MDS.DPS.TransactionalData.Outbound.Domain\Entities\Filters\PartitionOverFilter.cs
MDS.DPS.TransactionalData.Outbound.Domain\Entities\MappingSet.cs
MDS.DPS.TransactionalData.Outbound.Domain\Entities\MdoMapping.cs
MDS.DPS.TransactionalData.Outbound.API\Validators\RankFilterValidator.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\Orchestrator\MdoDataFetchingStrategyParser.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\Orchestrator\TransactionalDataCommandParser.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\Orchestrator\DataFetchingCommand.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\DataFetchingStrategies\DynamicQueryBuilder.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\DataFetchingStrategies\ViewProvider.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\DataFetchingStrategies\LatestReferenceBeforeSQLPartBuilder.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\DataFetchingStrategies\CmdpDataFetchingStrategy.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\DataFetchingStrategies\HyperScaleDataFetchingStrategy.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\DataFetchingStrategies\CassandraFilterProcessor.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\DataFetchingStrategies\CassandraFilterProcessingStrategy.cs
MDS.DPS.TransactionalData.Outbound.Service\Services\DataFetchingStrategies\FiltersExtensions.cs
MDS.DPS.TransactionalData.Outbound.API\Builders\CsvTransactionalDataCommandBuilder.cs
```

## Recommended Next Steps

1. Build a C# vs Go contract test matrix for representative ids:
   - CMDP-only
   - Cassandra-only
   - Hyperscale timeseries
   - Hyperscale curves
   - Hyperscale surfaces
   - hybrid switchover
   - HPFC force-to-CMDP
2. Capture generated SQL/CQL and response shape from both apps.
3. Add Go regression tests for every C# query-shape mismatch found.
4. Finish aggregation and shape parity.
5. Finish Mesap parity if the endpoint is required.
6. Add deployment documentation and production config examples.
