package domain

type ExecutableQuery struct {
	ID           Identifier
	DataCategory DataCategory
	Source       SourceKind
	Filters      FilterSet
	IndexRange   *IndexRange
	Statement    string
	Parameters   map[string]any
	Arguments    []any
}
