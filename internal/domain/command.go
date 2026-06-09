package domain

type Command struct {
	IDs             []Identifier
	DataCategory    DataCategory
	Columns         []string
	IncludeOffset   bool
	TargetTimeZone  string
	HasAggregations bool
	HasShape        bool
	Filters         FilterSet
	Mappings        []Mapping
	Source          SourceKind
	QuoteIndices    []int
	IndexRange      *IndexRange
}

type IndexRange struct {
	Start int
	End   int
}

type Mapping struct {
	ID           Identifier
	DataCategory DataCategory
	Source       SourceKind
	ViewName     string
	IndexField   string
	Resolution   string
	CassandraID  string
	HyperscaleID *Identifier
	SwitchOver   string
	SplitQuery   bool
	Timezone     string
	Columns      []ColumnMapping
}

type ColumnMapping struct {
	MDSName             string
	SourceName          string
	DataType            string
	IsKey               bool
	IsProjectable       bool
	OrderPriority       *int
	KeyColumnOrdering   *int
	ValueColumnOrdering *int
}
