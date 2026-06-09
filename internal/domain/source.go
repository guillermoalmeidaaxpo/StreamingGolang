package domain

type DataCategory string

const (
	Curves     DataCategory = "curves"
	Surfaces   DataCategory = "surfaces"
	TimeSeries DataCategory = "timeseries"
)

type SourceKind string

const (
	SourceCMDP       SourceKind = "cmdp"
	SourceCassandra  SourceKind = "cassandra"
	SourceHyperscale SourceKind = "hyperscale"
	SourceMesap      SourceKind = "mesap"
)
