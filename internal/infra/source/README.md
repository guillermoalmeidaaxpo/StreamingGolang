# Source Adapters

Data fetching should be source-oriented rather than a direct copy of the C#
strategy hierarchy.

Expected adapters:

- `mssql` for CMDP and mapping SQL Server access.
- `cassandra` for Cassandra-hosted transactional data.
- `hyperscale` if it remains a distinct storage/backend contract.
- `mesap` for Mesap transition endpoints.

Application code should depend on small source interfaces and execution plans,
not on database-specific implementation details.
