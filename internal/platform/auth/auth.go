package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"streaming-golang/internal/platform/config"
)

type Authenticator struct {
	enabled      bool
	verifier     *oidc.IDTokenVerifier
	audiences    map[string]struct{}
	allowedRoles map[string]struct{}
}

type Principal struct {
	Subject string
	Roles   []string
	Claims  map[string]any
}

type contextKey string

const principalKey contextKey = "principal"

func New(ctx context.Context, cfg config.Auth) (*Authenticator, error) {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode == "" || mode == "disabled" || mode == "none" {
		return &Authenticator{}, nil
	}
	if mode != "jwt" && mode != "oidc" {
		return nil, fmt.Errorf("unsupported OUTBOUND_AUTH_MODE %q", cfg.Mode)
	}
	if cfg.Authority == "" {
		return nil, errors.New("OUTBOUND_AUTH_ISSUER is required when authentication is enabled")
	}
	if len(cfg.Audiences) == 0 {
		return nil, errors.New("OUTBOUND_AUTH_AUDIENCES is required when authentication is enabled")
	}

	providerCtx := ctx
	if !cfg.RequireHTTPSMetadata {
		providerCtx = oidc.InsecureIssuerURLContext(providerCtx, cfg.Authority)
		providerCtx = context.WithValue(providerCtx, oauth2.HTTPClient, insecureHTTPClient())
	}

	provider, err := oidc.NewProvider(providerCtx, cfg.Authority)
	if err != nil {
		return nil, fmt.Errorf("discover oidc provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		SkipClientIDCheck: true,
	})

	return &Authenticator{
		enabled:      true,
		verifier:     verifier,
		audiences:    set(cfg.Audiences),
		allowedRoles: set(cfg.AllowedRoles),
	}, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	if a == nil || !a.enabled {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		idToken, err := a.verifier.Verify(r.Context(), rawToken)
		if err != nil {
			http.Error(w, "invalid bearer token", http.StatusUnauthorized)
			return
		}

		var claims tokenClaims
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		if !a.acceptsAudience([]string(claims.Audience)) {
			http.Error(w, "invalid token audience", http.StatusUnauthorized)
			return
		}
		if !a.acceptsRole(claims.Roles) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		principal := Principal{
			Subject: idToken.Subject,
			Roles:   claims.Roles,
			Claims:  claims.Raw,
		}
		ctx := context.WithValue(r.Context(), principalKey, principal)
		ctx = context.WithValue(ctx, "raw_bearer_token", rawToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey).(Principal)
	return principal, ok
}

func (a *Authenticator) acceptsAudience(audiences []string) bool {
	for _, audience := range audiences {
		if _, ok := a.audiences[audience]; ok {
			return true
		}
	}
	return false
}

func (a *Authenticator) acceptsRole(roles []string) bool {
	if len(a.allowedRoles) == 0 {
		return true
	}

	for _, role := range roles {
		if _, ok := a.allowedRoles[role]; ok {
			return true
		}
	}
	return false
}

type tokenClaims struct {
	Audience audience       `json:"aud"`
	Roles    []string       `json:"roles"`
	Raw      map[string]any `json:"-"`
}

func (c *tokenClaims) UnmarshalJSON(data []byte) error {
	type alias tokenClaims
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var parsed alias
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	*c = tokenClaims(parsed)
	c.Raw = raw
	return nil
}

type audience []string

func (a *audience) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*a = []string{single}
		return nil
	}

	var many []string
	if err := json.Unmarshal(data, &many); err != nil {
		return err
	}
	*a = many
	return nil
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	return token, token != ""
}

func set(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result[value] = struct{}{}
		}
	}
	return result
}

func insecureHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
