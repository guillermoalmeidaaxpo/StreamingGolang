# C# to Rust Migration Guide: Transactional Outbound API

## 1. Overview and Goals

This document outlines the architectural strategy and technical roadmap for migrating the Transactional Outbound API from its original C# implementation (and its Go port) to **Rust**.

The primary goals are:
- **Performance:** Leverage Rust's zero-cost abstractions and safe concurrency.
- **Safety:** Utilize Rust's ownership model to eliminate data races and memory leaks.
- **Maintainability:** Strictly follow **SOLID** principles and **Hexagonal Architecture**.
- **Correctness:** Ensure feature parity with the C# implementation, including complex strategy selection and hybrid splits.

---

## 2. Framework and Library Mapping

To ensure architectural consistency, we map each layer of the C# and Go projects to the most idiomatic and high-performance Rust equivalents.

| Category | C# (Source of Truth) | Go (Port) | Rust (Migration) |
| :--- | :--- | :--- | :--- |
| **Async Runtime** | Task Parallel Library | Goroutines / Channels | **Tokio** |
| **Web Framework** | ASP.NET Core (Minimal APIs) | Custom Router / Net/Http | **Axum** |
| **Serialization** | System.Text.Json / Newtonsoft | Encoding/Json | **Serde** |
| **Validation** | FluentValidation | GoJsonSchema / Manual | **Validator / JsonSchema** |
| **Parsing** | ANTLR4 | ANTLR4 (Go target) | **antlr-rust** |
| **Database (SQL)** | Entity Framework / Dapper | Go-Mssqldb | **Tiberius** |
| **Database (Cassandra)** | Cassandra C# Driver | Gocql | **Scylla-Rust-Driver** |
| **Caching** | StackExchange.Redis | Go-Redis | **Fred** |
| **Logging** | Serilog / App Insights | Slog | **Tracing / Tracing-AppInsights** |
| **DI Container** | Microsoft.Extensions.DependencyInjection | Manual (Constructors) | **Manual (Arc/Traits)** |

---

## 3. Architectural Alignment: Hexagonal Architecture

We will adopt a **Hexagonal Architecture (Ports & Adapters)** to decouple business logic from infrastructure.

### 3.1 Domain Layer (`src/domain/`)
- Contains pure business logic, entities, and value objects.
- **No dependencies** on frameworks or external libraries (except maybe `serde` for DTOs).
- **Core Types:** `Identifier`, `DataItem`, `Request`, `Command`, `FilterSet`, `SourceKind`.

### 3.2 Application Layer (`src/application/`)
- Orchestrates use cases.
- Defines **Ports** (Traits) that infrastructure must implement.
- **Key Components:**
    - `Pipeline`: Validates, parses, plans, and executes requests.
    - `Planner`: Logic for mapping resolution and strategy selection.
    - `Executor`: Parallel execution of plan steps.
    - `TransformationProcessor`: Handles time zone conversions and aggregations.

### 3.4 Logic Alignment with C# Source of Truth

The original C# implementation uses a sophisticated strategy selection and filter processing pipeline. The Rust implementation must mirror these specific components:

| C# Class | Rust Responsibility | Logic Detail |
| :--- | :--- | :--- |
| `MdoDataFetchingStrategyParser` | `StrategySelector` | Determines if the request goes to MESAP, Hyperscale, CMDP, or Cassandra. |
| `CassandraFilterProcessingStrategy` | `HybridExecutor` | Orchestrates the split between historical (Cassandra) and recent (CMDP) data. |
| `CassandraFilterProcessor` | `RangeSplitter` | Logic for dividing `ReferenceTime` filters using the `latest_global` watermark. |
| `CassandraDataFetcherService` | `CassandraRepository` | Async fetching from Cassandra with batching, retries, and RDP calculation. |

---

## 4. Business Logic Consistency (Reference: OutboundApi.Documentation.md & C# Source)

### 4.1 Strategy Selection Priority
Following `MdoDataFetchingStrategyParser.cs`:
1.  **MESAP:** If `isMesapEndpoint` and mapping ID exists in `MesapMappingStorage`.
2.  **Hyperscale:** If `mappings.HyperScaleId.HasValue`.
3.  **CMDP (Forced):** If `hasShape`, `hasAggregations`, ID is in `Hpfc_Ids_To_Cmdp`, or TimeZone != `Europe/Zurich`.
4.  **Cassandra:** Default strategy for transactional data when the above don't apply.

### 4.2 The "Watermark" Split (CassandraFilterProcessor)
The "Hybrid" logic in `CassandraFilterProcessingStrategy` splits the request based on a `latestGlobal` timestamp (the watermark):

*   **Watermark Discovery:** Call `GetMinMaxReferenceTimeDeliveryStart` to find the `MaxReferenceTime` for the mapping.
*   **Filter Splitting:**
    *   **Cassandra Range:** `[LowerLimit, Watermark)`
    *   **CMDP Range:** `[Watermark, UpperLimit]`
*   **Execution:**
    *   If aggregations are present, **Skip Cassandra** and route entirely to CMDP.
    *   Otherwise, execute both ranges in parallel and merge results.

### 4.3 Cassandra Execution Details
Mirroring `CassandraDataFetcherService.cs`:
*   **Batching:** Use `tokio::sync::Semaphore` to limit concurrent queries (mapping to `MaxQueriesCassandraInParallel`).
*   **RDP Calculation:** Post-process results to calculate `RelativeDeliveryPeriod` if the column was requested.
*   **Retry Policy:** Implement robust async retries with exponential backoff for `NoHostAvailable`, `ReadTimeout`, etc.

### Single Responsibility Principle (SRP)
- Each module has one reason to change.
- `src/domain/filters.rs` handles filter models.
- `src/application/planner.rs` handles only the execution plan construction.

### Open/Closed Principle (OCP)
- Use **Traits** to allow extending functionality without modifying existing code.
- New fetching strategies can be added by implementing the `QueryStrategy` trait.

```rust
pub trait QueryStrategy: Send + Sync {
    fn plan(&self, command: Command) -> Vec<Command>;
}

pub struct SplitQueryStrategy {
    pub queries_count: usize,
    pub split_days: i64,
}

impl QueryStrategy for SplitQueryStrategy {
    fn plan(&self, command: Command) -> Vec<Command> {
        // Implementation logic for splitting queries
    }
}
```

### Liskov Substitution Principle (LSP)
- Implementations of traits (like `Repository`) must be interchangeable.
- The `Executor` depends on `Box<dyn Repository>` or generic `R: Repository`, ensuring any repository can be used.

### Interface Segregation Principle (ISP)
- Instead of one giant `Repository` trait, use smaller traits like `MappingResolver`, `WatermarkProvider`, and `StreamProvider`.
- Components only depend on the specific behavior they need.

### Dependency Inversion Principle (DIP)
- High-level modules (`Pipeline`) depend on abstractions (`traits`), not low-level modules (`MssqlRepository`).
- Use **Dependency Injection** via constructors (passing `Arc<dyn Trait>` or using generics).

---

## 5. Key Design Patterns & Implementation Details

### 5.1 Request Pipeline (Async Orchestration)
The pipeline will mirror the Go implementation but with Rust's stronger error handling:

```rust
pub struct Pipeline {
    validator: Arc<dyn Validator>,
    parser: Arc<dyn FilterParser>,
    planner: Arc<dyn Planner>,
    executor: Arc<dyn Executor>,
}

impl Pipeline {
    pub async fn execute(&self, ctx: RequestContext, reqs: Vec<Request>) -> Result<Response, AppError> {
        self.validator.validate(&reqs)?;
        let plan = self.planner.build_plan(ctx, reqs).await?;
        let response = self.executor.execute(plan).await?;
        Ok(response)
    }
}
```

### 5.2 Strategy Selection & Hybrid Split
The logic for choosing between Cassandra and CMDP (and splitting requests across the watermark) will be encapsulated in the `Planner`.

- **Hybrid Split:** When a request range crosses the "watermark", the planner creates two `Command` objects: one for Cassandra (historical) and one for CMDP (recent).

### 5.3 Streaming with `async-stream`
To handle large datasets efficiently, we will use `async-stream` to yield items as they are received from the database:

```rust
pub fn stream_items(plan: Plan) -> impl Stream<Item = Result<DataItem, AppError>> {
    stream! {
        for step in plan.steps {
            let mut stream = repository.stream(step.query).await?;
            while let Some(item) = stream.next().await {
                yield Ok(transform(item));
            }
        }
    }
}
```

### 5.4 Concurrency
Use `tokio::spawn` and `FuturesUnordered` for parallel query execution in the `Executor`. This ensures that slow queries don't block the entire request.

---

## 6. Migration Roadmap

1.  **Phase 1: Foundation**
    - Setup Rust project with `axum`, `tokio`, and `serde`.
    - Implement `domain` models and `apperr` (Error handling).
2.  **Phase 2: Core Logic**
    - Port `FilterParser` (ANTLR integration).
    - Implement the `Planner` and `QueryStrategy` logic.
    - Implement the `TransformationProcessor`.
3.  **Phase 3: Infrastructure (Adapters)**
    - Implement `MssqlRepository` and `CassandraRepository`.
    - Implement `MappingResolver` with Redis/SQL backing.
4.  **Phase 4: API & Integration**
    - Build `axum` routes and JSON/CSV streaming handlers.
    - End-to-end testing against existing test suites.

---

## 7. Testing and Validation Strategy

To ensure correctness and behavioral parity with C# and Go:

### 7.1 Unit Testing
- Use **Traits** to mock infrastructure (SQL, Cassandra, Redis).
- Recommended Library: `mockall` for automatic trait mocking.
- Focus on `Planner` logic, `FilterParser` accuracy, and `TransformationProcessor` edge cases (e.g., Leap Years, DST transitions).

### 7.2 Integration Testing
- Use **Testcontainers-rs** to run real MSSQL and Cassandra instances during CI.
- Verify that SQL/CQL queries generated by the repositories are compatible with the respective database engines.

### 7.3 Parity Testing
- Create a test suite that runs identical requests against the C#, Go, and Rust implementations.
- Compare JSON/CSV outputs (ignoring field order) to ensure exact data parity.

---

## 8. Comparison: Go vs Rust Implementation

| Feature | Go | Rust |
| :--- | :--- | :--- |
| **Concurrency** | Goroutines + Channels | Async/Await + Futures |
| **Error Handling** | `if err != nil` | `Result<T, E>` + `?` operator |
| **Memory** | Garbage Collected | Ownership & Borrow Checker |
| **Polymorphism** | Interfaces (Implicit) | Traits (Explicit) |
| **Performance** | High (but GC pauses) | Maximum (Predictable) |

---

## 9. API Surface and Contract Mapping (Reference: OutboundApi.Documentation.md)

The Rust implementation will replicate the exact endpoint structure and JSON contracts to ensure client compatibility.

### 9.1 Endpoint Catalog (Axum Router)
The router will be organized into resource groups (`curves`, `timeseries`, `surfaces`, `generic`, `lite`, `metadata`, `datatrace`, `mesaptransition`) with the following route variants:
- `productive/{resource}`
- `design/{resource}`
- `validation/{resource}`
- `migration/{resource}`

### 9.2 Request DTOs (Serde Models)
We will use `serde` with `rename_all = "PascalCase"` to match the C# naming convention:

```rust
#[derive(Debug, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct TransactionalDataRequest {
    pub ids: Vec<i64>,
    pub version_as_of: Option<DateTime<Utc>>,
    pub filters: Option<Filters>,
    pub transformations: Option<Transformations>,
    pub columns: Option<Vec<String>>,
    pub include_deleted: Option<bool>,
}
```

### 9.3 Content-Type Handling
For streaming endpoints (Section 4.2 of doc), the Rust implementation will support:
- `application/x-ndjson`: Leveraged by `axum` and `tokio-util` for line-delimited JSON.
- `application/json`: Standard buffered or chunked JSON responses.
- `text/csv`: Generated via the `csv` crate and streamed to the client.

### 9.4 Caching Strategy (Redis)
Following Section 8.3 of the documentation, the `ICacheHandler` equivalent in Rust will use `fred` with `TokenCredential` support to manage the response cache, ensuring identical cache keys and TTL behavior as the C# implementation.

---

## 10. Deep-Dive Technical Design

### 10.1 Filtering and ANTLR Parsing
To maintain parity with `FilterExpressionParser.cs` and the Go port, the Rust implementation will use the `antlr-rust` library.

*   **Grammar:** Reuse the existing `.g4` grammars (Lexer and Parser).
*   **Visitor Pattern:** Implement a `FilterASTVisitor` that traverses the parse tree and yields a `Vec<FilterNode>`.
*   **Polymorphic Filters:** Leverage Rust `enums` for `FilterNode` (e.g., `Comparison(ComparisonFilter)`, `RankOver(RankFilter)`).
*   **Error Handling:** Implement a custom `ErrorListener` to capture and format syntax errors exactly as the C# `ExpressionTokenizer` does.

### 10.2 Validation Layer
Following `GenericRequestValidator.cs` and `FiltersParsingValidator.cs`, validation will be a two-stage process:

1.  **Schema Validation:** Use `jsonschema` for fast, declarative check of the request structure.
2.  **Domain Validation:** Implement a `Validator` trait with specific implementations for:
    *   `RequestValidator`: Ensures IDs are valid and unique.
    *   `FilterParsingValidator`: Integrates with the ANTLR parser to validate expressions during the validation phase (DIP).
    *   `AggregationValidator`: Checks if keys/values match supported delivery periods and patterns.

### 10.3 Refined Strategy Selection Logic
The `StrategySelector` (port of `MdoDataFetchingStrategyParser.cs`) will use the following logic:

```rust
pub fn get_strategy(mappings: &MappingSet, context: RequestContext) -> StrategyKind {
    // 1. MESAP Check
    if context.is_mesap_endpoint && mesap_storage.contains(mappings.timeseries_id) {
        return StrategyKind::Mesap;
    }

    // 2. Hyperscale Check
    if mappings.hyperscale_id.is_some() {
        return StrategyKind::Hyperscale;
    }

    // 3. CMDP Forced Conditions
    if context.has_shape || context.has_aggregations || 
       FORCED_CMDP_IDS.contains(&mappings.timeseries_id) ||
       mappings.cassandra_id.is_none() ||
       get_timezone(mappings.timeseries_id) != "Europe/Zurich" {
        return StrategyKind::Cmdp;
    }

    // 4. Default
    StrategyKind::Cassandra
}
```

### 10.4 Mapping Composition and Resolution
The `Planner` will orchestrate the resolution of identifiers using the `MappingResolver` trait:

*   **Multi-Source Resolution:** Fetch data from `CmdpMappingDatabase` and `MesapMappingDatabase`.
*   **Composition:** Combine the results into a `MappingSet` which contains all cross-referenced IDs (CassandraId, HyperScaleId, MesapId).
*   **Watermark Injection:** The `MappingResolver` is also responsible for providing the `latest_global` watermark, which is then used by the `HybridExecutor` to split the workload.


w Checker |
| **Polymorphism** | Interfaces (Implicit) | Traits (Explicit) |
| **Performance** | High (but GC pauses) | Maximum (Predictable) |

---

## 9. API Surface and Contract Mapping (Reference: OutboundApi.Documentation.md)

The Rust implementation will replicate the exact endpoint structure and JSON contracts to ensure client compatibility.

### 9.1 Endpoint Catalog (Axum Router)
The router will be organized into resource groups (`curves`, `timeseries`, `surfaces`, `generic`, `lite`, `metadata`, `datatrace`, `mesaptransition`) with the following route variants:
- `productive/{resource}`
- `design/{resource}`
- `validation/{resource}`
- `migration/{resource}`

### 9.2 Request DTOs (Serde Models)
We will use `serde` with `rename_all = "PascalCase"` to match the C# naming convention:

```rust
#[derive(Debug, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct TransactionalDataRequest {
    pub ids: Vec<i64>,
    pub version_as_of: Option<DateTime<Utc>>,
    pub filters: Option<Filters>,
    pub transformations: Option<Transformations>,
    pub columns: Option<Vec<String>>,
    pub include_deleted: Option<bool>,
}
```

### 9.3 Content-Type Handling
For streaming endpoints (Section 4.2 of doc), the Rust implementation will support:
- `application/x-ndjson`: Leveraged by `axum` and `tokio-util` for line-delimited JSON.
- `application/json`: Standard buffered or chunked JSON responses.
- `text/csv`: Generated via the `csv` crate and streamed to the client.

### 9.4 Caching Strategy (Redis)
Following Section 8.3 of the documentation, the `ICacheHandler` equivalent in Rust will use `fred` with `TokenCredential` support to manage the response cache, ensuring identical cache keys and TTL behavior as the C# implementation.

---

## 10. Deep-Dive Technical Design

### 10.1 Filtering and ANTLR Parsing
To maintain parity with `FilterExpressionParser.cs` and the Go port, the Rust implementation will use the `antlr-rust` library.

*   **Grammar:** Reuse the existing `.g4` grammars (Lexer and Parser).
*   **Visitor Pattern:** Implement a `FilterASTVisitor` that traverses the parse tree and yields a `Vec<FilterNode>`.
*   **Polymorphic Filters:** Leverage Rust `enums` for `FilterNode` (e.g., `Comparison(ComparisonFilter)`, `RankOver(RankFilter)`).
*   **Error Handling:** Implement a custom `ErrorListener` to capture and format syntax errors exactly as the C# `ExpressionTokenizer` does.

### 10.2 Validation Layer
Following `GenericRequestValidator.cs` and `FiltersParsingValidator.cs`, validation will be a two-stage process:

1.  **Schema Validation:** Use `jsonschema` for fast, declarative check of the request structure.
2.  **Domain Validation:** Implement a `Validator` trait with specific implementations for:
    *   `RequestValidator`: Ensures IDs are valid and unique.
    *   `FilterParsingValidator`: Integrates with the ANTLR parser to validate expressions during the validation phase (DIP).
    *   `AggregationValidator`: Checks if keys/values match supported delivery periods and patterns.

### 10.3 Refined Strategy Selection Logic
The `StrategySelector` (port of `MdoDataFetchingStrategyParser.cs`) will use the following logic:

```rust
pub fn get_strategy(mappings: &MappingSet, context: RequestContext) -> StrategyKind {
    // 1. MESAP Check
    if context.is_mesap_endpoint && mesap_storage.contains(mappings.timeseries_id) {
        return StrategyKind::Mesap;
    }

    // 2. Hyperscale Check
    if mappings.hyperscale_id.is_some() {
        return StrategyKind::Hyperscale;
    }

    // 3. CMDP Forced Conditions
    if context.has_shape || context.has_aggregations || 
       FORCED_CMDP_IDS.contains(&mappings.timeseries_id) ||
       mappings.cassandra_id.is_none() ||
       get_timezone(mappings.timeseries_id) != "Europe/Zurich" {
        return StrategyKind::Cmdp;
    }

    // 4. Default
    StrategyKind::Cassandra
}
```

### 10.4 Mapping Composition and Resolution
The `Planner` will orchestrate the resolution of identifiers using the `MappingResolver` trait:

*   **Multi-Source Resolution:** Fetch data from `CmdpMappingDatabase` and `MesapMappingDatabase`.
*   **Composition:** Combine the results into a `MappingSet` which contains all cross-referenced IDs (CassandraId, HyperScaleId, MesapId).
*   **Watermark Injection:** The `MappingResolver` is also responsible for providing the `latest_global` watermark, which is then used by the `HybridExecutor` to split the workload.


