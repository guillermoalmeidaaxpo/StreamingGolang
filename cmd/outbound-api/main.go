package main

import (
	"context"
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

	authenticator, err := auth.New(context.Background(), cfg.Auth)
	if err != nil {
		logger.Error("authentication configuration failed", slog.Any("error", err))
		os.Exit(1)
	}

	mappingResolver := transactional.MappingResolver(transactional.StaticMappingResolver{})
	queryBuilder := transactional.QueryBuilder(transactional.PlaceholderQueryBuilder{})
	cmdpMappingDB, err := mssql.OpenSQLServer(cfg.ConnectionStrings.CmdpMappingDatabase)
	if err != nil {
		logger.Error("cmdp mapping sql configuration failed", slog.Any("error", err))
		os.Exit(1)
	}
	if cmdpMappingDB != nil {
		defer cmdpMappingDB.Close()
		mdsDB, err := mssql.OpenSQLServer(cfg.ConnectionStrings.MdsDatabase)
		if err != nil {
			logger.Error("mds mapping sql configuration failed", slog.Any("error", err))
			os.Exit(1)
		}
		if mdsDB != nil {
			defer mdsDB.Close()
		}
		mappingResolver = mssql.NewMappingResolver(cmdpMappingDB, mdsDB)
		queryBuilder = mssql.NewCMDPQueryBuilder()
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
			cassandra.NewCassandraQueryBuilder(cfg.Cassandra.TableMappings),
		)
	}

	cmdpSQLDB, err := mssql.OpenSQLServer(cfg.ConnectionStrings.CmdpSQLDatabase)
	if err != nil {
		logger.Error("cmdp sql database configuration failed", slog.Any("error", err))
		os.Exit(1)
	}
	repositories := make(map[domain.SourceKind]transactional.Repository)
	if cmdpSQLDB != nil {
		defer cmdpSQLDB.Close()
		repositories[domain.SourceCMDP] = mssql.NewRepository(cmdpSQLDB)
		repositories[domain.SourceHyperscale] = mssql.NewRepository(cmdpSQLDB)
	}
	if cassandraSession != nil {
		repositories[domain.SourceCassandra] = cassandra.NewRepository(cassandraSession)
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
