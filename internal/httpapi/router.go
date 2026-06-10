package httpapi

import (
	"log/slog"
	"net/http"

	"streaming-golang/internal/app/authz"
	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/platform/auth"
	"streaming-golang/internal/platform/config"
)

type Dependencies struct {
	Config                config.Config
	Logger                *slog.Logger
	TransactionalPipeline *transactional.Pipeline
	Authenticator         *auth.Authenticator
	LicenseValidator      authz.LicenseValidator
}

func NewRouter(deps Dependencies) http.Handler {
	h := handlers{
		config:                deps.Config,
		transactionalPipeline: deps.TransactionalPipeline,
	}

	publicMux := http.NewServeMux()
	publicMux.HandleFunc("GET /health/startup", h.health)
	publicMux.HandleFunc("GET /health/liveness", h.health)
	publicMux.HandleFunc("GET /health/readiness", h.health)
	publicMux.HandleFunc("GET /api/v1/info", h.info)

	protectedMux := http.NewServeMux()
	registerPostRoutes(protectedMux, transactionalRoutes, h.transactional)
	registerPostRoutes(protectedMux, transactionalStreamingRoutes, h.transactionalStream)
	registerPostRoutes(protectedMux, genericRoutes, h.generic)
	registerPostRoutes(protectedMux, genericStreamingRoutes, h.genericStream)
	registerGetRoutes(protectedMux, liteRoutes, h.liteCSV)
	registerPostRoutes(protectedMux, metadataRoutes, h.notImplemented("metadata"))
	registerPostRoutes(protectedMux, metadataRangeRoutes, h.notImplemented("metadata-range"))
	registerPostRoutes(protectedMux, dataTraceRoutes, h.notImplemented("datatrace"))
	registerPostRoutes(protectedMux, mesapGenericRoutes, h.notImplemented("mesap-generic"))

	var protected http.Handler = protectedMux
	protected = licenseValidation(deps.LicenseValidator, deps.Config.Build.Stage)(protected)
	if deps.Authenticator != nil {
		protected = deps.Authenticator.Middleware(protected)
	}

	publicMux.Handle("/api/v1/", protected)

	return requestLogger(deps.Logger)(
		recoverer(deps.Logger)(
			correlationID(publicMux),
		),
	)
}

func registerPostRoutes(mux *http.ServeMux, routes []string, handler http.HandlerFunc) {
	for _, route := range routes {
		mux.HandleFunc("POST /api/v1"+route, handler)
	}
}

func registerGetRoutes(mux *http.ServeMux, routes []string, handler http.HandlerFunc) {
	for _, route := range routes {
		mux.HandleFunc("GET /api/v1"+route, handler)
	}
}
