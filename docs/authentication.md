# Authentication And Authorization

The C# API has two distinct security steps:

1. JWT bearer authentication and role authorization.
2. MDO/license authorization against the external authorization service.

The Go service follows the same security model but keeps both concerns separate.

Code source of truth:

```text
internal/platform/auth/auth.go
internal/app/authz/license.go
internal/httpapi/license_middleware.go
```

When JWT validation, role/scope checks, authorization payloads, or authorization
paths change, update this document in the same change.

## JWT / Entra ID

Protected `/api/v1/*` endpoints use OIDC/JWT validation when authentication is
enabled. Health checks and `/api/v1/info` stay public.

Environment variables:

```powershell
$env:OUTBOUND_AUTH_MODE = "jwt"
$env:OUTBOUND_AUTH_ISSUER = "https://login.microsoftonline.com/<tenant-id>/v2.0"
$env:OUTBOUND_AUTH_AUDIENCES = "api://<application-id>,<other-audience>"
$env:OUTBOUND_AUTH_ALLOWED_ROLES = "RoleA,RoleB"
$env:OUTBOUND_AUTH_REQUIRE_HTTPS_METADATA = "true"
```

For local development the default is:

```powershell
$env:OUTBOUND_AUTH_MODE = "disabled"
```

## License Validation

The Go API uses an HTTP license validator that calls the authorization service
before request handlers execute. The incoming bearer token is forwarded to the
authorization API.

The C# flow calls two authorization endpoints, and the Go service mirrors that:

- `POST /api/v1/DataUniverse/BulkAuthorize`
- `POST /api/v1/TimeSeries/Authorize`

Data universe payload:

```json
{
  "type": "TransactionalDataOutbound",
  "universeName": null,
  "action": "Read",
  "internalCorrelationId": "<correlation-id>",
  "mdoIds": [312091001],
  "stageId": 3
}
```

Time series payload:

```json
{
  "identifiers": [312091001],
  "stageId": 3,
  "internalCorrelationId": "<correlation-id>"
}
```

The base URL is configured with:

```powershell
$env:OUTBOUND_AUTHORIZATION_API_BASE_URL = "https://authorization-hs.lab.mds.axpo.com"
```
