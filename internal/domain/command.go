package domain

import "time"

type Command struct {
	IDs                 []Identifier
	DataCategory        DataCategory
	Columns             []string
	VersionAsOf         *time.Time
	IncludeDeleted      bool
	IncludeIdentifier   bool
	IncludeOffset       bool
	FilterTimeZone      string
	TargetTimeZone      string
	HasAggregations     bool
	HasShape            bool
	LatestReferenceTime bool
	Filters             FilterSet
	Mappings            []Mapping
	Source              SourceKind
	QuoteIndices        []int
	IndexRange          *IndexRange
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
	Views        MappingViews
	IndexField   string
	Resolution   string
	CassandraID  string
	HyperscaleID *Identifier
	SwitchOver   string
	SplitQuery   bool
	Timezone     string
	Columns      []ColumnMapping
}

type MappingViews struct {
	LatestVersion                     string
	LatestReferenceTime               string
	LatestVersionWithCreatedOn        string
	LatestReferenceTimeWithCreatedOn  string
	GetByCreatedOn                    string
	GetByCreatedOnLatestReferenceTime string
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
