# Authentication And Authorization

The C# API has two distinct security steps:

1. JWT bearer authentication and role authorization.
2. MDO/license authorization against the external authorization service.

The Go service follows the same security model but keeps both concerns separate.

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

The next implementation step is to replace `NoopLicenseValidator` with an HTTP
adapter that calls the authorization server with:

- action: `Read`
- type: `TransactionalDataOutbound`
- requested MDO identifiers
- stage
- internal correlation id

That should be wired after JWT authentication and before transactional handlers.
