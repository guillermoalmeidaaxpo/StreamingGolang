# Configuration

Configuration is Go-native and container-friendly.

Load order:

1. `configs/default.yaml`
2. `configs/<environment>.yaml`
3. `OUTBOUND_*` environment variables

The environment defaults to `development` and can be changed with:

```powershell
$env:OUTBOUND_ENV = "productive"
```

The config directory defaults to `configs` and can be changed with:

```powershell
$env:OUTBOUND_CONFIG_DIR = "C:\path\to\configs"
```

Canonical file shape:

```yaml
meta:
  build_number: local
  stage: development

http:
  host: ""
  port: 8080
  read_header_timeout: 5s

auth:
  mode: jwt
  issuer: https://login.microsoftonline.com/<tenant-id>/v2.0
  audiences:
    - api://<app-id>
  allowed_roles:
    - DataReader
  require_https_metadata: true

authorization_api:
  base_url: https://authorization.example
  authorize_path: /api/v1/TimeSeries/
  universe_authorize_path: /api/v1/DataUniverse/BulkAuthorize
  timeout: 30s

datastores:
  cmdp_sql:
    dsn: ${OUTBOUND_CMDP_SQL_DSN}
  redis:
    address: localhost:6379
    tls: false
  cassandra:
    data_centers:
      primary:
        - cassandra1.example
    keyspace: ts
    primary_data_center: primary
    connection_timeout: 10s
    read_timeout: 12s

database:
  connect_retry:
    timeout: 60s
    count: 3
    interval: 10s
  command_retry:
    command_timeout: 120s
    count: 3
    interval: 5s
    max_interval: 120s

logging:
  default_level: debug
  microsoft_level: info
  system_level: info
  application_insights_default_level: info

execution:
  batch_size: 10000
  stream_batch_size: 1000
  max_sql_parallel: 5
```

Common environment overrides:

```powershell
$env:OUTBOUND_HTTP_PORT = "8080"
$env:OUTBOUND_AUTH_MODE = "jwt"
$env:OUTBOUND_AUTH_ISSUER = "https://login.microsoftonline.com/<tenant-id>/v2.0"
$env:OUTBOUND_AUTH_AUDIENCES = "api://<app-id>,<other-audience>"
$env:OUTBOUND_AUTH_ALLOWED_ROLES = "DataReader"
$env:OUTBOUND_AUTHORIZATION_API_BASE_URL = "https://authorization.example"
$env:OUTBOUND_CMDP_SQL_DSN = "<connection-string>"
$env:OUTBOUND_REDIS_ADDRESS = "<host>:6380"
$env:OUTBOUND_CASSANDRA_KEYSPACE = "ts"
$env:OUTBOUND_DATABASE_COMMAND_TIMEOUT = "120s"
$env:OUTBOUND_LOG_LEVEL = "debug"
```

Secrets should be provided by environment variables or the deployment platform,
not committed to config files.
