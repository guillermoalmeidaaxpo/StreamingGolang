package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"

	"streaming-golang/internal/platform/config"
)

func Open(cfg config.Cassandra) (*gocql.Session, error) {
	if len(cfg.DataCenters) == 0 {
		return nil, nil
	}

	hosts := make([]string, 0)
	for _, dcHosts := range cfg.DataCenters {
		hosts = append(hosts, dcHosts...)
	}

	if len(hosts) == 0 {
		return nil, nil
	}

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = cfg.Keyspace
	cluster.Port = cfg.Port
	cluster.Timeout = cfg.ReadTimeout
	cluster.ConnectTimeout = cfg.ConnectionTimeout
	cluster.NumConns = cfg.LocalConnections
	
	if cfg.PrimaryDataCenter != "" {
		cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy(cfg.PrimaryDataCenter)
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("create cassandra session: %w", err)
	}

	return session, nil
}

func IsConfigured(cfg config.Cassandra) bool {
	return len(cfg.DataCenters) > 0 && cfg.Keyspace != ""
}
