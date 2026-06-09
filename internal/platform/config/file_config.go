package config

type fileConfig struct {
	Meta             metaConfig             `yaml:"meta"`
	HTTP             httpConfig             `yaml:"http"`
	Auth             authConfig             `yaml:"auth"`
	AuthorizationAPI authorizationAPIConfig `yaml:"authorization_api"`
	Datastores       datastoresConfig       `yaml:"datastores"`
	Database         databaseConfig         `yaml:"database"`
	Logging          loggingConfig          `yaml:"logging"`
	Execution        executionConfig        `yaml:"execution"`
}

type metaConfig struct {
	BuildNumber string `yaml:"build_number"`
	Stage       string `yaml:"stage"`
}

type httpConfig struct {
	Host              string   `yaml:"host"`
	Port              int      `yaml:"port"`
	ReadHeaderTimeout duration `yaml:"read_header_timeout"`
}

type authConfig struct {
	Mode                 string   `yaml:"mode"`
	Issuer               string   `yaml:"issuer"`
	Audiences            []string `yaml:"audiences"`
	AllowedRoles         []string `yaml:"allowed_roles"`
	RequireHTTPSMetadata bool     `yaml:"require_https_metadata"`
}

type authorizationAPIConfig struct {
	BaseURL               string   `yaml:"base_url"`
	AuthorizePath         string   `yaml:"authorize_path"`
	UniverseAuthorizePath string   `yaml:"universe_authorize_path"`
	Timeout               duration `yaml:"timeout"`
}

type datastoresConfig struct {
	CmdpSQL         sqlConfig       `yaml:"cmdp_sql"`
	MappingSQL      sqlConfig       `yaml:"mapping_sql"`
	MdsSQL          sqlConfig       `yaml:"mds_sql"`
	MesapMappingSQL sqlConfig       `yaml:"mesap_mapping_sql"`
	Redis           redisConfig     `yaml:"redis"`
	Cassandra       cassandraConfig `yaml:"cassandra"`
}

type sqlConfig struct {
	DSN string `yaml:"dsn"`
}

type redisConfig struct {
	Address                 string   `yaml:"address"`
	TLS                     bool     `yaml:"tls"`
	IgnoreAllowedUsersCheck bool     `yaml:"ignore_allowed_users_check"`
	AllowedUsers            []string `yaml:"allowed_users"`
}

type cassandraConfig struct {
	DataCenters          map[string][]string `yaml:"data_centers"`
	Keyspace             string              `yaml:"keyspace"`
	PrimaryDataCenter    string              `yaml:"primary_data_center"`
	LocalConnections     int                 `yaml:"local_connections"`
	LocalMaxConnections  int                 `yaml:"local_max_connections"`
	RemoteConnections    int                 `yaml:"remote_connections"`
	RemoteMaxConnections int                 `yaml:"remote_max_connections"`
	Port                 int                 `yaml:"port"`
	MaxParallelQueries   int                 `yaml:"max_parallel_queries"`
	ConnectionTimeout    duration            `yaml:"connection_timeout"`
	ReadTimeout          duration            `yaml:"read_timeout"`
	TableMappings        map[string]string   `yaml:"table_mappings"`
}

type databaseConfig struct {
	ConnectRetry connectRetryConfig `yaml:"connect_retry"`
	CommandRetry commandRetryConfig `yaml:"command_retry"`
}

type connectRetryConfig struct {
	Timeout  duration `yaml:"timeout"`
	Count    int      `yaml:"count"`
	Interval duration `yaml:"interval"`
}

type commandRetryConfig struct {
	CommandTimeout duration `yaml:"command_timeout"`
	Count          int      `yaml:"count"`
	Interval       duration `yaml:"interval"`
	MaxInterval    duration `yaml:"max_interval"`
}

type loggingConfig struct {
	DefaultLevel                    string `yaml:"default_level"`
	MicrosoftLevel                  string `yaml:"microsoft_level"`
	SystemLevel                     string `yaml:"system_level"`
	ApplicationInsightsDefaultLevel string `yaml:"application_insights_default_level"`
}

type executionConfig struct {
	BatchSize                       int `yaml:"batch_size"`
	StreamBatchSize                 int `yaml:"stream_batch_size"`
	OptimizedBatchSize              int `yaml:"optimized_batch_size"`
	MaxSQLParallel                  int `yaml:"max_sql_parallel"`
	MaxCassandraParallel            int `yaml:"max_cassandra_parallel"`
	MaxSQLConnections               int `yaml:"max_sql_connections"`
	MaxCassandraConnections         int `yaml:"max_cassandra_connections"`
	ReferenceTimeSplitDays          int `yaml:"reference_time_split_days"`
	CassandraReferenceTimeSplitDays int `yaml:"cassandra_reference_time_split_days"`
}
