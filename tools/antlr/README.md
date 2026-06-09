# ANTLR Tooling

Place `antlr-4.13.2-complete.jar` in this directory, then run:

```powershell
go generate ./internal/query/parser/antlr
```

The generator requires Java 11+. The helper script prefers `JAVA_HOME`, then
known JDK 17/21 install locations, then `java` on `PATH`.

The generated Go parser is intentionally isolated under
`internal/query/parser/antlr/generated` and consumed through the application-level
`transactional.FilterParser` interface.
