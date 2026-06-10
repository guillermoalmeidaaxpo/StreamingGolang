package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"streaming-golang/internal/app/authz"
	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
	"streaming-golang/internal/httpapi"
	"streaming-golang/internal/infra/cassandra"
	"streaming-golang/internal/infra/mssql"
	"streaming-golang/internal/infra/redis"
	"streaming-golang/internal/platform/auth"
	"streaming-golang/internal/platform/config"
	"streaming-golang/internal/platform/server"
	antlrparser "streaming-golang/internal/query/parser/antlr"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration failed", slog.Any("error", err))
		os.Exit(1)
	}
	logStartupConfiguration(logger, cfg)

	authenticator, err := auth.New(context.Background(), cfg.Auth)
	if err != nil {
		logger.Error("authentication configuration failed", slog.Any("error", err))
		os.Exit(1)
	}

	cmdpSQLDB, err := mssql.OpenSQLServer(cfg.ConnectionStrings.CmdpSQLDatabase)
	if err != nil {
		logger.Error("cmdp sql database configuration failed", slog.Any("error", err))
		os.Exit(1)
	}

	mappingResolver := transactional.MappingResolver(transactional.StaticMappingResolver{})
	queryBuilder := transactional.QueryBuilder(transactional.PlaceholderQueryBuilder{})
	var mdsDB *sql.DB
	cmdpMappingDB, err := mssql.OpenSQLServer(cfg.ConnectionStrings.CmdpMappingDatabase)
	if err != nil {
		logger.Error("cmdp mapping sql configuration failed", slog.Any("error", err))
		os.Exit(1)
	}
	if cmdpMappingDB != nil {
		defer cmdpMappingDB.Close()
		mdsDB, err = mssql.OpenSQLServer(cfg.ConnectionStrings.MdsDatabase)
		if err != nil {
			logger.Error("mds mapping sql configuration failed", slog.Any("error", err))
			os.Exit(1)
		}
		if mdsDB != nil {
			defer mdsDB.Close()
		}
		mappingResolver = mssql.NewMappingResolver(cmdpMappingDB, mdsDB, cmdpSQLDB, logger)
		queryBuilder = transactional.NewCompositeQueryBuilder(
			mssql.NewCMDPQueryBuilder(),
			mssql.NewHyperscaleQueryBuilder(),
		)
	}

	redisClient, err := redis.Open(cfg.Redis)
	if err != nil {
		logger.Error("redis configuration failed", slog.Any("error", err))
		os.Exit(1)
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	cassandraSession, err := cassandra.Open(cfg.Cassandra)
	if err != nil {
		logger.Error("cassandra session configuration failed", slog.Any("error", err))
		os.Exit(1)
	}
	if cassandraSession != nil {
		defer cassandraSession.Close()
		queryBuilder = transactional.NewCompositeQueryBuilder(
			queryBuilder,
			cassandra.NewCassandraQueryBuilder(cfg.Cassandra.TableMappings, cfg.Cassandra.Keyspace),
		)
	}

	repositories := make(map[domain.SourceKind]transactional.Repository)
	if cmdpSQLDB != nil {
		defer cmdpSQLDB.Close()
		repositories[domain.SourceCMDP] = mssql.NewRepository(cmdpSQLDB, logger)
	}
	if mdsDB != nil {
		repositories[domain.SourceHyperscale] = mssql.NewRepository(mdsDB, logger)
	}
	if cassandraSession != nil {
		repositories[domain.SourceCassandra] = cassandra.NewRepository(cassandraSession, logger)
	}

	pipeline := transactional.NewPipeline(
		transactional.NewValidator(),
		antlrparser.New(),
		transactional.NewPlanner(transactional.WithQueryStrategy(transactional.SplitQueryStrategy{
			QueriesCount:           cfg.Stream.MaxQueriesInParallel,
			ReferenceTimeSplitDays: cfg.Stream.ReferenceTimeSplitDays,
		}), transactional.WithMappingResolver(mappingResolver), transactional.WithQueryBuilder(queryBuilder)),
		transactional.NewExecutor(repositories, cfg.Stream.MaxQueriesInParallel),
	)

	httpLicenseValidator := authz.NewHttpLicenseValidator(
		cfg.AuthorizationAPI.BaseURL,
		cfg.AuthorizationAPI.Endpoint,
		cfg.AuthorizationAPI.UniverseEndpoint,
		cfg.AuthorizationAPI.Timeout,
		logger,
	)

	licenseValidator := authz.NewAllowedUserLicenseValidator(
		httpLicenseValidator,
		redisClient,
		cfg.Redis.IgnoreAllowedUsersCheck,
		cfg.Redis.AllowedUsersInCache,
	)

	router := httpapi.NewRouter(httpapi.Dependencies{
		Config:                cfg,
		Logger:                logger,
		TransactionalPipeline: pipeline,
		Authenticator:         authenticator,
		LicenseValidator:      licenseValidator,
	})

	srv := server.New(cfg.HTTP, router)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("outbound api started", slog.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", slog.Any("error", err))
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("outbound api stopped")
}

func logStartupConfiguration(logger *slog.Logger, cfg config.Config) {
	workingDir, err := os.Getwd()
	if err != nil {
		workingDir = "unknown: " + err.Error()
	}

	logger.Info("configuration loaded",
		slog.String("working_directory", workingDir),
		slog.String("stage", cfg.Build.Stage),
		slog.String("outbound_env", os.Getenv("OUTBOUND_ENV")),
		slog.String("outbound_config_dir", os.Getenv("OUTBOUND_CONFIG_DIR")),
	)

	logSQLDatastore(logger, "cmdp_sql", cfg.ConnectionStrings.CmdpSQLDatabase)
	logSQLDatastore(logger, "mapping_sql", cfg.ConnectionStrings.CmdpMappingDatabase)
	logSQLDatastore(logger, "mds_sql", cfg.ConnectionStrings.MdsDatabase)
	logSQLDatastore(logger, "mesap_mapping_sql", cfg.ConnectionStrings.MesapMappingDatabase)
}

func logSQLDatastore(logger *slog.Logger, name, dsn string) {
	logger.Info("sql datastore configured",
		slog.String("datastore", name),
		slog.Bool("configured", mssql.IsConfiguredDSN(dsn)),
		slog.String("driver", mssql.DriverNameForDSN(dsn)),
		slog.String("auth_mode", mssql.AuthModeForDSN(dsn)),
	)
}
