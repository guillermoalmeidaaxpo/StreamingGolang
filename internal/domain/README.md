# Domain

This package owns the business language of the outbound API.

Keep it small and free of HTTP, JSON, SQL, Redis, Cassandra, and framework
details. Application packages orchestrate use cases with these types; transport
and infrastructure packages adapt them to external protocols.
