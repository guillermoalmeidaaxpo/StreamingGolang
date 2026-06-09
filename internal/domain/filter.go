package domain

type FilterSet struct {
	Expressions []string
	Nodes       []FilterNode
}

type FilterNode interface {
	filterNode()
}

type ComparisonFilter struct {
	Raw      string
	Field    string
	Operator string
	Value    FilterValue
}

func (ComparisonFilter) filterNode() {}

type RankOverFilter struct {
	Raw         string
	PartitionBy []string
	OrderBy     []SortExpression
	Bounds      []RankOverBound
}

func (RankOverFilter) filterNode() {}

type SortExpression struct {
	Field     string
	Direction string
}

type RankOverBound struct {
	Raw   string
	Start string
	End   string
}

type LatestExpression struct {
	Raw      string
	Field    string
	Operator string
	Value    FilterValue
}

type FilterValue struct {
	Kind       FilterValueKind
	Raw        string
	Function   string
	TimeZone   string
	Start      string
	End        string
	Arithmetic *TimeArithmetic
	Arguments  []LatestExpression
}

type FilterValueKind string

const (
	FilterValueNumber                FilterValueKind = "number"
	FilterValueText                  FilterValueKind = "text"
	FilterValuePointInTime           FilterValueKind = "point_in_time"
	FilterValueTimeInterval          FilterValueKind = "time_interval"
	FilterValueTimeIntervalPointTime FilterValueKind = "time_interval_point_in_time"
	FilterValueLatest                FilterValueKind = "latest"
	FilterValueLatestGlobal          FilterValueKind = "latest_global"
	FilterValueGeneric               FilterValueKind = "generic"
)

type TimeArithmetic struct {
	Operator string
	Period   string
}
