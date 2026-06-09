package config

import (
	"fmt"
)

func Load() (Config, error) {
	fileCfg, err := loadFiles()
	if err != nil {
		return Config{}, err
	}

	host := stringEnv("OUTBOUND_HTTP_HOST", fileCfg.HTTP.Host)
	port := intEnv("OUTBOUND_HTTP_PORT", fileCfg.HTTP.Port)
	readHeaderTimeout := durationEnv("OUTBOUND_HTTP_READ_HEADER_TIMEOUT", fileCfg.HTTP.ReadHeaderTimeout.Duration())

	return Config{
		HTTP: HTTP{
			Address:           fmt.Sprintf("%s:%d", host, port),
			ReadHeaderTimeout: readHeaderTimeout,
		},
		Build: Build{
			Number: stringEnv("OUTBOUND_BUILD_NUMBER", fileCfg.Meta.BuildNumber),
			Stage:  stringEnv("OUTBOUND_STAGE", fileCfg.Meta.Stage),
		},
		Auth: Auth{
			Mode:                 stringEnv("OUTBOUND_AUTH_MODE", fileCfg.Auth.Mode),
			Authority:            stringEnv("OUTBOUND_AUTH_ISSUER", fileCfg.Auth.Issuer),
			Audiences:            csvEnv("OUTBOUND_AUTH_AUDIENCES", fileCfg.Auth.Audiences),
			AllowedRoles:         csvEnv("OUTBOUND_AUTH_ALLOWED_ROLES", fileCfg.Auth.AllowedRoles),
			RequireHTTPSMetadata: boolEnv("OUTBOUND_AUTH_REQUIRE_HTTPS_METADATA", fileCfg.Auth.RequireHTTPSMetadata),
		},
		AuthorizationAPI: AuthorizationAPI{
			BaseURL:          stringEnv("OUTBOUND_AUTHORIZATION_API_BASE_URL", fileCfg.AuthorizationAPI.BaseURL),
			Endpoint:         stringEnv("OUTBOUND_AUTHORIZATION_API_AUTHORIZE_PATH", fileCfg.AuthorizationAPI.AuthorizePath),
			UniverseEndpoint: stringEnv("OUTBOUND_AUTHORIZATION_API_UNIVERSE_AUTHORIZE_PATH", fileCfg.AuthorizationAPI.UniverseAuthorizePath),
			Timeout:          durationEnv("OUTBOUND_AUTHORIZATION_API_TIMEOUT", fileCfg.AuthorizationAPI.Timeout.Duration()),
		},
		ConnectionStrings: ConnectionStrings{
			CmdpSQLDatabase:      stringEnv("OUTBOUND_CMDP_SQL_DSN", fileCfg.Datastores.CmdpSQL.DSN),
			CmdpMappingDatabase:  stringEnv("OUTBOUND_MAPPING_SQL_DSN", fileCfg.Datastores.MappingSQL.DSN),
			MdsDatabase:          stringEnv("OUTBOUND_MDS_SQL_DSN", fileCfg.Datastores.MdsSQL.DSN),
			MesapMappingDatabase: stringEnv("OUTBOUND_MESAP_MAPPING_SQL_DSN", fileCfg.Datastores.MesapMappingSQL.DSN),
		},
		Redis: Redis{
			URL:                     stringEnv("OUTBOUND_REDIS_ADDRESS", fileCfg.Datastores.Redis.Address),
			UseSSL:                  boolEnv("OUTBOUND_REDIS_TLS", fileCfg.Datastores.Redis.TLS),
			IgnoreAllowedUsersCheck: boolEnv("OUTBOUND_REDIS_IGNORE_ALLOWED_USERS_CHECK", fileCfg.Datastores.Redis.IgnoreAllowedUsersCheck),
			AllowedUsersInCache:     csvEnv("OUTBOUND_REDIS_ALLOWED_USERS", fileCfg.Datastores.Redis.AllowedUsers),
		},
		Cassandra: Cassandra{
			DataCenters:          fileCfg.Datastores.Cassandra.DataCenters,
			Keyspace:             stringEnv("OUTBOUND_CASSANDRA_KEYSPACE", fileCfg.Datastores.Cassandra.Keyspace),
			PrimaryDataCenter:    stringEnv("OUTBOUND_CASSANDRA_PRIMARY_DATA_CENTER", fileCfg.Datastores.Cassandra.PrimaryDataCenter),
			LocalConnections:     intEnv("OUTBOUND_CASSANDRA_LOCAL_CONNECTIONS", fileCfg.Datastores.Cassandra.LocalConnections),
			LocalMaxConnections:  intEnv("OUTBOUND_CASSANDRA_LOCAL_MAX_CONNECTIONS", fileCfg.Datastores.Cassandra.LocalMaxConnections),
			RemoteConnections:    intEnv("OUTBOUND_CASSANDRA_REMOTE_CONNECTIONS", fileCfg.Datastores.Cassandra.RemoteConnections),
			RemoteMaxConnections: intEnv("OUTBOUND_CASSANDRA_REMOTE_MAX_CONNECTIONS", fileCfg.Datastores.Cassandra.RemoteMaxConnections),
			Port:                 intEnv("OUTBOUND_CASSANDRA_PORT", fileCfg.Datastores.Cassandra.Port),
			MaxParallelQueries:   intEnv("OUTBOUND_CASSANDRA_MAX_PARALLEL_QUERIES", fileCfg.Datastores.Cassandra.MaxParallelQueries),
			ConnectionTimeout:    durationEnv("OUTBOUND_CASSANDRA_CONNECTION_TIMEOUT", fileCfg.Datastores.Cassandra.ConnectionTimeout.Duration()),
			ReadTimeout:          durationEnv("OUTBOUND_CASSANDRA_READ_TIMEOUT", fileCfg.Datastores.Cassandra.ReadTimeout.Duration()),
			TableMappings:        fileCfg.Datastores.Cassandra.TableMappings,
		},
		Database: Database{
			ConnectRetry: ConnectRetry{
				Timeout:  durationEnv("OUTBOUND_DATABASE_CONNECT_TIMEOUT", fileCfg.Database.ConnectRetry.Timeout.Duration()),
				Count:    intEnv("OUTBOUND_DATABASE_CONNECT_RETRY_COUNT", fileCfg.Database.ConnectRetry.Count),
				Interval: durationEnv("OUTBOUND_DATABASE_CONNECT_RETRY_INTERVAL", fileCfg.Database.ConnectRetry.Interval.Duration()),
			},
			CommandRetry: CommandRetry{
				CommandTimeout: durationEnv("OUTBOUND_DATABASE_COMMAND_TIMEOUT", fileCfg.Database.CommandRetry.CommandTimeout.Duration()),
				Count:          intEnv("OUTBOUND_DATABASE_COMMAND_RETRY_COUNT", fileCfg.Database.CommandRetry.Count),
				Interval:       durationEnv("OUTBOUND_DATABASE_COMMAND_RETRY_INTERVAL", fileCfg.Database.CommandRetry.Interval.Duration()),
				MaxInterval:    durationEnv("OUTBOUND_DATABASE_COMMAND_RETRY_MAX_INTERVAL", fileCfg.Database.CommandRetry.MaxInterval.Duration()),
			},
		},
		Logging: Logging{
			DefaultLevel:                    stringEnv("OUTBOUND_LOG_LEVEL", fileCfg.Logging.DefaultLevel),
			MicrosoftLevel:                  stringEnv("OUTBOUND_LOG_MICROSOFT_LEVEL", fileCfg.Logging.MicrosoftLevel),
			SystemLevel:                     stringEnv("OUTBOUND_LOG_SYSTEM_LEVEL", fileCfg.Logging.SystemLevel),
			ApplicationInsightsDefaultLevel: stringEnv("OUTBOUND_LOG_APPLICATION_INSIGHTS_LEVEL", fileCfg.Logging.ApplicationInsightsDefaultLevel),
		},
		Stream: Stream{
			BatchStreamSize:                   intEnv("OUTBOUND_EXECUTION_STREAM_BATCH_SIZE", fileCfg.Execution.StreamBatchSize),
			BatchOptimizedSize:                intEnv("OUTBOUND_EXECUTION_OPTIMIZED_BATCH_SIZE", fileCfg.Execution.OptimizedBatchSize),
			MaxQueriesInParallel:              intEnv("OUTBOUND_EXECUTION_MAX_SQL_PARALLEL", fileCfg.Execution.MaxSQLParallel),
			MaxQueriesCassandraInParallel:     intEnv("OUTBOUND_EXECUTION_MAX_CASSANDRA_PARALLEL", fileCfg.Execution.MaxCassandraParallel),
			MaxConcurrentDatabaseConnections:  intEnv("OUTBOUND_EXECUTION_MAX_SQL_CONNECTIONS", fileCfg.Execution.MaxSQLConnections),
			MaxConcurrentCassandraConnections: intEnv("OUTBOUND_EXECUTION_MAX_CASSANDRA_CONNECTIONS", fileCfg.Execution.MaxCassandraConnections),
			ReferenceTimeSplitDays:            intEnv("OUTBOUND_EXECUTION_REFERENCE_TIME_SPLIT_DAYS", fileCfg.Execution.ReferenceTimeSplitDays),
			ReferenceTimeCassandraSplitDays:   intEnv("OUTBOUND_EXECUTION_CASSANDRA_REFERENCE_TIME_SPLIT_DAYS", fileCfg.Execution.CassandraReferenceTimeSplitDays),
		},
		Split: Split{
			BatchSize: intEnv("OUTBOUND_EXECUTION_BATCH_SIZE", fileCfg.Execution.BatchSize),
		},
	}, nil
}
