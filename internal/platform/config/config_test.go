package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadMergesDefaultEnvironmentAndEnvVars(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("OUTBOUND_CONFIG_DIR", dir)
	t.Setenv("OUTBOUND_ENV", "development")
	t.Setenv("OUTBOUND_HTTP_PORT", "9090")

	defaultConfig := `
meta:
  build_number: base
  stage: productive
http:
  port: 8080
  read_header_timeout: 5s
auth:
  mode: jwt
  issuer: https://issuer.example
  audiences:
    - aud-a
    - aud-b
  allowed_roles:
    - DataReader
  require_https_metadata: true
datastores:
  redis:
    address: redis-base
    tls: true
execution:
  batch_size: 10000
  stream_batch_size: 1000
`
	developmentConfig := `
meta:
  stage: development
auth:
  mode: disabled
datastores:
  redis:
    tls: false
`

	if err := os.WriteFile(filepath.Join(dir, "default.yaml"), []byte(defaultConfig), 0o600); err != nil {
		t.Fatalf("write default config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "development.yaml"), []byte(developmentConfig), 0o600); err != nil {
		t.Fatalf("write development config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.HTTP.Address != ":9090" {
		t.Fatalf("expected env var port override, got %q", cfg.HTTP.Address)
	}
	if cfg.HTTP.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("expected read header timeout from default file, got %s", cfg.HTTP.ReadHeaderTimeout)
	}
	if cfg.Build.Stage != "development" {
		t.Fatalf("expected development stage override, got %q", cfg.Build.Stage)
	}
	if cfg.Auth.Mode != "disabled" {
		t.Fatalf("expected development auth mode override, got %q", cfg.Auth.Mode)
	}
	if cfg.Redis.UseSSL {
		t.Fatal("expected redis tls override to be false")
	}
	if cfg.Stream.BatchStreamSize != 1000 || cfg.Split.BatchSize != 10000 {
		t.Fatalf("expected default execution settings to remain, got stream=%d split=%d", cfg.Stream.BatchStreamSize, cfg.Split.BatchSize)
	}
}
