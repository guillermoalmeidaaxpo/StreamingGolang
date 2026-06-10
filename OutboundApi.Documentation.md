# MDS DPS Transactional Data Outbound API — Technical Documentation

## 1. Overview

This API provides transactional market data access for:
- **Timeseries**
- **Curves**
- **Surfaces**
- **Metadata / Metadata Range**
- **Data Trace**
- **MESAP transition paths**
- **Synchronous and streaming responses**

Target runtime: **.NET 8**.

---

## 2. Solution Structure (high level)

```text
API
 ├─ Endpoint mapping (Minimal APIs)
 ├─ Middleware, validation, authorization, telemetry
 └─ DTO contracts

Service
 ├─ Handlers (MdoHandler, GenericHandler, MetadataHandler, DataTraceHandler)
 ├─ Orchestration (TransactionalDataService, TransactionalDataCsvService)
 ├─ Strategy selection (MdoDataFetchingStrategyParser)
 └─ Mapping/translation and caching orchestration

Domain
 ├─ Entities, filters, parsers, expressions
 └─ Repository interfaces and shared abstractions

Infrastructure
 ├─ Repository implementations (SQL/Cassandra/MESAP)
 ├─ External API clients
 ├─ Throttling + concurrency gates
 └─ Configuration + resilience policies
```

---

## 3. Runtime Request Pipeline

```mermaid
flowchart LR
    Client[Client / APIM] --> API[Minimal API Endpoints]
    API --> Val[Validators + Request Strategies]
    Val --> H[Service Handlers]
    H --> Orch[Orchestrators]
    Orch --> Strat[MDO Strategy Parser]
    Strat --> SQL[SQL Repositories]
    Strat --> CAS[Cassandra Repository]
    Strat --> MESAP[MESAP Repository]
    H --> Cache[Redis Cache]
    API --> Tel[Application Insights + Logging]
```

---

## 4. API Endpoint Catalog

Base route:
- `api/v{version:apiVersion}`
- Default version: `v1`

Info endpoint:
- `GET /api/v{version}/info`

Health endpoints:
- `/health/startup`
- `/health/liveness`
- `/health/readiness`

### 4.1 Transactional data endpoints (POST)

Resource groups:
- `curves`
- `timeseries`
- `surfaces`

Route variants:
- `{resource}`
- `design/{resource}`
- `validation/{resource}`
- `productive/{resource}`
- `migration/{resource}`

### 4.2 Transactional stream endpoints (POST)

Route patterns:
- `{resource}/streaming`
- `design/{resource}/streaming`
- `validation/{resource}/streaming`
- `productive/{resource}/streaming`
- `migration/{resource}/streaming`

Supported stream content types:
- `application/x-ndjson`
- `application/json`

### 4.3 Generic CSV endpoints

Synchronous CSV download (POST):
- `generic`
- `design/generic`
- `validation/generic`
- `productive/generic`
- `migration/generic`

Streaming CSV (POST):
- `generic/streaming`
- `design/generic/streaming`
- `validation/generic/streaming`
- `productive/generic/streaming`
- `migration/generic/streaming`

### 4.4 Lite endpoint (GET)

- `lite`
- `design/lite`
- `validation/lite`
- `productive/lite`

### 4.5 Metadata endpoints (POST)

Metadata:
- `{resource}/metadata`
- `design/{resource}/metadata`
- `validation/{resource}/metadata`
- `productive/{resource}/metadata`

Metadata range:
- `{resource}/metadata/range`
- `design/{resource}/metadata/range`
- `validation/{resource}/metadata/range`
- `productive/{resource}/metadata/range`

(`resource` = `curves | timeseries | surfaces`)

### 4.6 Data Trace endpoints (POST)

- `datatrace`
- `design/datatrace`
- `validation/datatrace`
- `productive/datatrace`

### 4.7 MESAP transition endpoints (POST)

MESAP generic:
- `mesaptransition/generic`
- `validation/mesaptransition/generic`
- `productive/mesaptransition/generic`

---

## 5. Request Objects (Contract Layer)

### TransactionalDataRequest[]
- `Ids: List<long>`
- `VersionAsOf: DateTime?`
- `Filters: Filters?`
- `Transformations: Transformations?`
- `Columns: List<string>?`
- `IncludeDeleted: bool?`

### GenericRequest
- `Ids: long[]?` or `Id: long?`
- `MdoIdArray` computed from Ids/Id
- `VersionAsOf`, `Filters`, `Transformations`, `Columns`, `IncludeDeleted`

### LiteRequest (query)
- `id: long`
- `from: string`
- `to: string?`

### MetadataRequest
- `Id: long`
- `ReferenceTime: DateTimeOffset`

### MetadataRangeRequest
- `Ids: long[]`
- `StartTime: DateTimeOffset`
- `EndTime: DateTimeOffset`

### DataTraceRequest
- `Id: long`
- `ReferenceTime: DateTimeOffset`
- `VersionAsOf: DateTimeOffset?`

---

## 6. Object Connections (DI and Runtime)

```mermaid
classDiagram
    class OutboundApi
    class OutboundStreamApi
    class MetadataApi
    class DataTraceApi
    class MesapApi

    class IMdoHandler
    class IGenericHandler
    class IMetadataHandler
    class IMetadataRangeHandler
    class IDataTraceHandler

    class ITransactionalDataService
    class ITransactionalDataCsvService
    class ITransactionalDataCommandParser
    class IMdoDataFetchingStrategyParser
    class IMdoDataFetchingStrategy

    class IDataAccessRepository
    class ICassandraMdoRepository
    class IMetadataRepository
    class IDataTraceRepository
    class IMdoMappingRepository
    class IMesapRepository

    OutboundApi --> IMdoHandler
    OutboundApi --> IGenericHandler
    OutboundStreamApi --> IMdoHandler
    OutboundStreamApi --> IGenericHandler
    MetadataApi --> IMetadataHandler
    MetadataApi --> IMetadataRangeHandler
    DataTraceApi --> IDataTraceHandler
    MesapApi --> IGenericHandler

    IMdoHandler --> ITransactionalDataService
    IMdoHandler --> ITransactionalDataCommandParser
    IGenericHandler --> ITransactionalDataCsvService
    IGenericHandler --> ITransactionalDataCommandParser

    ITransactionalDataService --> IMdoDataFetchingStrategyParser
    ITransactionalDataCsvService --> IMdoDataFetchingStrategyParser
    IMdoDataFetchingStrategyParser --> IMdoDataFetchingStrategy

    IMdoDataFetchingStrategy --> IDataAccessRepository
    IMdoDataFetchingStrategy --> ICassandraMdoRepository
    IMdoDataFetchingStrategy --> IMesapRepository

    IMetadataHandler --> IMetadataRepository
    IMetadataHandler --> IMdoMappingRepository
    IMetadataRangeHandler --> IMetadataRepository
    IMetadataRangeHandler --> IMdoMappingRepository
    IDataTraceHandler --> IDataTraceRepository
```

---

## 7. Data Source Strategy Selection

`MdoDataFetchingStrategyParser` chooses strategy using mapping and request context:

Priority:
1. **MESAP strategy** (when mesap endpoint + mapped mesap id)
2. **Hyperscale strategy** (when HyperScaleId exists)
3. **CMDP strategy** when any of:
   - shape requested
   - aggregations requested
   - id belongs to `Hpfc_Ids_To_Cmdp`
   - no Cassandra id
   - timezone rule requires CMDP path
4. **Cassandra strategy** otherwise

```mermaid
flowchart TD
    A[Incoming DataFetchingCommand] --> B{Mesap endpoint + id in MesapMappingStorage?}
    B -- yes --> M[MESAP Strategy]
    B -- no --> C{HyperScaleId exists?}
    C -- yes --> H[Hyperscale Strategy]
    C -- no --> D{ShouldUseCmdpStrategy?}
    D -- yes --> P[CMDP Strategy]
    D -- no --> S[Cassandra Strategy]
```

---

## 8. Infrastructure and External Dependencies

### 8.1 SQL connections
- CMDP SQL (`CmdpSqlDatabase`)
- Mapping SQL (`CmdpMappingDatabase`)
- Hyperscale SQL (`MdsDatabase`)
- MESAP mapping SQL (`MesapMappingDatabase`)

### 8.2 Cassandra
- `CassandraSessionFactory`
- `CassandraMdoRepository`
- Prepared statement cache
- Concurrency gate (`CassandraConnectionGate`)

### 8.3 Redis
- `RedisDatabaseProvider` with `TokenCredential`
- Shared for:
  - response cache (`ICacheHandler`)
  - query rate limiting (`ICmdpQueryRateLimiter`)
  - global connection gate (`ICmdpGlobalConnectionGate`)

### 8.4 External HTTP APIs
- License validation API (`ILicenseValidatorApiClient`)
- MESAP API integration (`IMesapRepository` via `MesapDataRepository`)
- Configured resiliency:
  - exponential retry
  - circuit breaker

### 8.5 Observability
- Application Insights telemetry + processors
- SQL dependency enrichment/filtering
- Cassandra telemetry tracker and module
- Correlation and caller enrichers

---

## 9. Main Execution Sequences

### 9.1 Synchronous transactional request

```mermaid
sequenceDiagram
    participant C as Client
    participant A as OutboundApi
    participant H as MdoHandler
    participant S as TransactionalDataService
    participant P as StrategyParser
    participant R as Repository

    C->>A: POST /api/v1/{resource}
    A->>H: ExecuteTransactionalDataProcessingSync(commands)
    H->>S: GetTransactionalData(command)
    S->>P: GetDataFetchingStrategy(...)
    P-->>S: Selected strategy
    S->>R: FetchData(command)
    R-->>S: Raw records
    S-->>H: TransformedData
    H-->>A: TsResponse
    A-->>C: 200 OK JSON
```

### 9.2 Streaming transactional request

```mermaid
sequenceDiagram
    participant C as Client
    participant A as OutboundStreamApi
    participant H as MdoHandler
    participant E as StreamingExecutionService

    C->>A: POST /api/v1/{resource}/streaming
    A->>H: GetTransactionalDataStream(commands)
    H-->>A: IAsyncEnumerable<Result<TransformedData>>
    A->>E: ExecuteTransactionalStreamAsync(...)
    E-->>C: streamed chunks (ndjson/json)
```

### 9.3 Metadata flow

```mermaid
sequenceDiagram
    participant C as Client
    participant A as MetadataApi
    participant H as MetadataHandler/MetadataRangeHandler
    participant M as MappingStorage
    participant R as MetadataRepository

    C->>A: POST /api/v1/{resource}/metadata
    A->>H: GetMetadata(request)
    H->>M: Load mappings
    H->>R: GetMetadata / GetMetadataRange
    R-->>H: metadata rows
    H-->>A: response DTO
    A-->>C: 200 OK
```

---

## 10. Middleware and Cross-cutting Components

Pipeline includes:
- Global exception handling
- Request tracing middleware
- HTTP logging
- Compression
- Authentication / Authorization
- Custom outbound middlewares (validation/raw/correlation/license)
- Health checks mapping

Validation stack:
1. Request/schema-level checks
2. FluentValidation validators
3. Query validators for lite routes
4. Domain/mapping-aware validation

---

## 11. Configuration Map (major sections)

- `SplitOptions`
- `StreamOptions`
- `ParallelTasksSettings`
- `LimitsSettings`
- `CassandraConfig`
- `CassandraRouting`
- `RedisOptions`
- `AuthorizationApiOptions`
- `MesapApiOptions`
- `CmdpRateLimiterOptions`
- `CmdpGlobalConnectionGateOptions`
- `ApplicationInsights`

---

## 12. Quick Dependency Diagram (projects)

```mermaid
flowchart LR
    API[MDS...Outbound.API] --> Service[MDS...Outbound.Service]
    API --> Infra[MDS...Outbound.Infrastructure]
    API --> Contract[MDS...Outbound.API.Contract]
    Service --> Domain[MDS...Outbound.Domain]
    Service --> Contract
    Infra --> Domain
```

---

## 13. Notes

- Endpoint mapping is centralized in:
  - `OutboundApi`, `OutboundStreamApi`, `MetadataApi`, `DataTraceApi`, `MesapApi`, `InfoApi`
- Service registrations are centralized in:
  - API: `API/Extensions/ServiceCollectionExtensions.cs`
  - Service: `Service/Extensions/ServiceCollectionExtensions.cs`
  - Infrastructure: `Infrastructure/Extensions/ServiceRegistrationExtensions.cs`

---

## 14. Validation Flow (detailed)

Validation is executed in middleware before handlers.

Core components:
- `ValidationMiddleware`
- `RequestValidationStrategyResolver`
- Request strategies:
  - `TransactionalDataRequestValidationStrategy`
  - `GenericRequestValidationStrategy`
  - `LiteRequestValidationStrategy`
  - `MetadataRequestValidationStrategy`
  - `MetadataRangeRequestValidationStrategy`
  - `DataTraceRequestValidationStrategy`

```mermaid
sequenceDiagram
    participant C as Client
    participant M as ValidationMiddleware
    participant R as StrategyResolver
    participant S as IRequestValidationStrategy
    participant V as FluentValidators
    participant N as Next Middleware/Endpoint

    C->>M: HTTP request
    M->>R: Resolve(context)
    R-->>M: concrete strategy or null
    M->>S: ValidateAsync(context)
    S->>V: Validate DTO/query
    alt invalid
        V-->>S: errors
        S-->>M: false + context.Items[ValidationErrors]
        M-->>C: 400 ProblemDetails
    else valid
        V-->>S: success
        S-->>M: true
        M->>N: next(context)
    end
```

### 14.1 Generic request validation internals

`GenericRequestDetailsValidator` performs high-value checks:
- mapping existence (`MappingStorage.Load`)
- duplicates and ids validation
- shape constraints (only curves / cmdp-compatible)
- projection + aggregation validation
- filter parsing (`FilterExpressionParser.Parse`)
- mesap filter validation (`MesapMappingStorage`, `ParsedFiltersMesapValidator`)
- filter column validation (`FilterColumnsRules`, `CustomFilterValidator`)
- estimated rows validation (`DataRowsNumberValidator`)

---

## 15. Filtering and Parsing Objects + Flow

Main filtering/parsing objects:
- `FilterExpressionParser`
- `FilterExpressionVisitor`
- `FilterProvider`
- `FilterMapper`
- `RawFilter`
- `Filter`
- `FilterSet`
- `MdoLimitsRequest`
- `FilterLimits`

```mermaid
flowchart TD
    A[Request.Filters.Expressions] --> B[FilterExpressionParser.Parse]
    B --> C[List<RawFilter>]
    C --> D[FilterMapper.Map]
    D --> E[FilterSet]
    E --> F{Need limits?}
    F -- yes --> G[IDataAccessRepository.GetMinMaxReferenceTimeDeliveryStart]
    G --> H[Set FilterLimits]
    F -- no --> I[Skip limits]
    H --> J{Has latest filter?}
    I --> J
    J -- yes --> K[Resolve latest via repository]
    J -- no --> L[Final FilterSet]
    K --> L
```

### 15.1 Where parsing is used

- Validation time:
  - `GenericRequestDetailsValidator.ParseFilters`
  - `DataRowsNumberValidator.ParseFilters`
- Command conversion/runtime:
  - `TransactionalDataCommandParser` -> `FilterProvider.GetFilters`

---

## 16. Statistics Service and Data Row Estimation Flow

Statistics objects:
- `IStatisticsService` / `StatisticsService`
- `IStatisticsRepository` / `StatisticsRepository`
- `IStatistics` implementations:
  - `TimeseriesStatistics`
  - `CurvesStatistics`
  - `SurfacesStatistics`

Row estimation entry point:
- `DataRowsNumberValidator`

```mermaid
sequenceDiagram
    participant V as DataRowsNumberValidator
    participant S as IStatisticsService
    participant R as IStatisticsRepository
    participant D as IDataAccessRepository(cmdp)
    participant C as IDataPointsCalculator

    V->>S: GetByMdoIdAsync(id, mappings, filters)
    S->>R: GetByMdoIdAsync(id, category)
    R-->>S: stats or null
    S->>D: GetMinMaxReferenceTimeDeliveryStart(...)
    D-->>S: filter limits
    S-->>V: IStatistics
    V->>C: GetEstimatedDataPoints(filters, mappings, stats)
    C-->>V: estimated row count
    V-->>V: compare with LimitsSettings.MaxDataPointsNumber
```

Behavior notes:
- For hyperscale MDOs, existing statistics are required.
- For CMDP/Cassandra paths, limits may be derived from repository queries.
- Surface category is excluded from row-estimation enforcement in validator logic.

---

## 17. Extended Object Connections (all major services)

```mermaid
flowchart LR
    subgraph API
        OA[OutboundApi]
        OSA[OutboundStreamApi]
        MA[MetadataApi]
        DTA[DataTraceApi]
        MSA[MesapApi]
        VM[ValidationMiddleware]
        RSV[RequestValidationStrategyResolver]
    end

    subgraph Service
        MH[MdoHandler]
        GH[GenericHandler]
        MDH[MetadataHandler]
        MDRH[MetadataRangeHandler]
        DTH[DataTraceHandler]
        TDS[TransactionalDataService]
        TDCSV[TransactionalDataCsvService]
        TDCP[TransactionalDataCommandParser]
        MDFSP[MdoDataFetchingStrategyParser]
        STS[StatisticsService]
        DRV[DataRowsNumberValidator]
    end

    subgraph Domain
        FP[FilterProvider]
        FEM[FilterExpressionParser]
        FSM[FilterSet/Filter/RawFilter]
    end

    subgraph Infrastructure
        DAR[IDataAccessRepository]
        CMR[ICassandraMdoRepository]
        MRR[IMetadataRepository]
        DTR[IDataTraceRepository]
        MMR[IMdoMappingRepository]
        MSR[IMesapRepository]
        SR[IStatisticsRepository]
        CH[ICacheHandler]
    end

    OA --> MH
    OA --> GH
    OSA --> MH
    OSA --> GH
    MA --> MDH
    MA --> MDRH
    DTA --> DTH
    MSA --> GH

    VM --> RSV
    RSV --> DRV

    MH --> TDS
    MH --> TDCP
    MH --> CH
    GH --> TDCSV
    GH --> TDCP

    TDS --> MDFSP
    TDCSV --> MDFSP
    TDCP --> FP
    FP --> FEM
    FP --> FSM

    MDFSP --> DAR
    MDFSP --> CMR
    MDFSP --> MSR
    MDH --> MRR
    MDH --> MMR
    MDRH --> MRR
    MDRH --> MMR
    DTH --> DTR

    DRV --> STS
    STS --> SR
    STS --> DAR
```

---

## 18. Parsing + Validation Connection Matrix

| Area | Main Class | Depends On | Purpose |
|---|---|---|---|
| Request strategy selection | `RequestValidationStrategyResolver` | `IEnumerable<IRequestValidationStrategy>` | Pick validator flow by route/method |
| Transactional request body validation | `TransactionalDataRequestValidationStrategy` | `IValidator<TransactionalDataRequest[]>` | Contract + business validation before endpoint |
| Generic request validation | `GenericRequestValidationStrategy` | `IValidator<GenericRequest>` | Validate generic and mesap-generic requests |
| Generic deep validation | `GenericRequestDetailsValidator` | mapping storage, mesap mapping, custom validators, row estimator | Validate ids, mappings, filters, projections, aggregations |
| Filter parsing | `FilterExpressionParser` | tokenizer + visitor | Convert expression strings into `RawFilter` objects |
| Filter mapping/runtime | `FilterProvider` + `FilterMapper` | `IDataAccessRepository`, mappings | Produce runtime `FilterSet` with defaults/latest resolution |
| Statistics lookup | `StatisticsService` | `IStatisticsRepository`, `IDataAccessRepository` | Get statistics and/or derive limits for estimation |
| Rows estimation | `DataRowsNumberValidator` | `IStatisticsService`, `IDataPointsCalculator` | Block oversized non-streaming requests |
