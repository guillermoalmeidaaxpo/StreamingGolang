# Filter Parser

This package is reserved for the filter language implementation.

The C# API currently uses ANTLR. The Go service should keep parser details behind
the `transactional.FilterParser` interface so the application does not depend on
ANTLR-generated types. That allows us to start with the existing grammar for
behavioral parity and later replace it with a smaller Go-native parser if the
grammar surface allows it.
