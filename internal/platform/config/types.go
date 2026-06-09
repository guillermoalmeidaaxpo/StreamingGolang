package config

import "time"

type Config struct {
	HTTP              HTTP
	Build             Build
	Auth              Auth
	AuthorizationAPI  AuthorizationAPI
	ConnectionStrings ConnectionStrings
	Redis             Redis
	Cassandra         Cassandra
	Database          Database
	Logging           Logging
	Stream            Stream
	Split             Split
}

type HTTP struct {
	Address           string
	ReadHeaderTimeout time.Duration
}

type Build struct {
	Number string
	Stage  string
}

type Auth struct {
	Mode                 string
	Authority            string
	Audiences            []string
	AllowedRoles         []string
	RequireHTTPSMetadata bool
}

type AuthorizationAPI struct {
	BaseURL          string
	Endpoint         string
	UniverseEndpoint string
	Timeout          time.Duration
}

type ConnectionStrings struct {
	CmdpSQLDatabase      string
	CmdpMappingDatabase  string
	MdsDatabase          string
	MesapMappingDatabase string
}

type Redis struct {
	URL                     string
	UseSSL                  bool
	IgnoreAllowedUsersCheck bool
	AllowedUsersInCache     []string
}

type Cassandra struct {
	DataCenters          map[string][]string
	Keyspace             string
	PrimaryDataCenter    string
	LocalConnections     int
	LocalMaxConnections  int
	RemoteConnections    int
	RemoteMaxConnections int
	Port                 int
	MaxParallelQueries   int
	ConnectionTimeout    time.Duration
	ReadTimeout          time.Duration
	TableMappings        map[string]string
}

type Database struct {
	ConnectRetry ConnectRetry
	CommandRetry CommandRetry
}

type ConnectRetry struct {
	Timeout  time.Duration
	Count    int
	Interval time.Duration
}

type CommandRetry struct {
	CommandTimeout time.Duration
	Count          int
	Interval       time.Duration
	MaxInterval    time.Duration
}

type Logging struct {
	DefaultLevel                    string
	MicrosoftLevel                  string
	SystemLevel                     string
	ApplicationInsightsDefaultLevel string
}

type Stream struct {
	BatchStreamSize                   int
	BatchOptimizedSize                int
	MaxQueriesInParallel              int
	MaxQueriesCassandraInParallel     int
	MaxConcurrentDatabaseConnections  int
	MaxConcurrentCassandraConnections int
	ReferenceTimeSplitDays            int
	ReferenceTimeCassandraSplitDays   int
}

type Split struct {
	BatchSize int
}
