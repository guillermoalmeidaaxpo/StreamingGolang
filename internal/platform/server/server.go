package server

import (
	"net/http"

	"streaming-golang/internal/platform/config"
)

func New(cfg config.HTTP, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Address,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}
}
