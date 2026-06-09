package authz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/platform/auth"
)

type LicenseValidator interface {
	ValidateReadAccess(ctx context.Context, request LicenseRequest) error
}

type LicenseRequest struct {
	Identifiers           []int64
	Stage                 string
	InternalCorrelationID string
}

type NoopLicenseValidator struct{}

func (NoopLicenseValidator) ValidateReadAccess(context.Context, LicenseRequest) error {
	return nil
}

type HttpLicenseValidator struct {
	client  *http.Client
	baseURL string
	path    string
}

func NewHttpLicenseValidator(baseURL, path string, timeout time.Duration) *HttpLicenseValidator {
	return &HttpLicenseValidator{
		client:  &http.Client{Timeout: timeout},
		baseURL: baseURL,
		path:    path,
	}
}

type authRequest struct {
	Action      string  `json:"action"`
	Type        string  `json:"type"`
	Identifiers []int64 `json:"identifiers"`
	Stage       string  `json:"stage"`
}

func (v *HttpLicenseValidator) ValidateReadAccess(ctx context.Context, req LicenseRequest) error {
	if v.baseURL == "" || v.baseURL == "NOT SET" {
		return nil
	}

	body := authRequest{
		Action:      "Read",
		Type:        "TransactionalDataOutbound",
		Identifiers: req.Identifiers,
		Stage:       req.Stage,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return apperr.Wrap(apperr.Internal, "marshal license request", err)
	}

	url := v.baseURL + v.path
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return apperr.Wrap(apperr.Internal, "create license request", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if req.InternalCorrelationID != "" {
		httpReq.Header.Set("X-Correlation-ID", req.InternalCorrelationID)
	}
	
	// Forward the Bearer token from the incoming request context to the Authorization API
	if _, ok := auth.PrincipalFromContext(ctx); ok {
		if rawToken, ok := ctx.Value("raw_bearer_token").(string); ok {
			httpReq.Header.Set("Authorization", "Bearer "+rawToken)
		}
	}

	resp, err := v.client.Do(httpReq)
	if err != nil {
		return apperr.Wrap(apperr.Unavailable, "call authorization service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return apperr.New(apperr.Forbidden, "Access denied: Missing license for one or more identifiers")
	}

	return apperr.New(apperr.Unavailable, fmt.Sprintf("authorization service returned status %d", resp.StatusCode))
}

type UserChecker interface {
	IsUserAllowed(ctx context.Context, userID string) (bool, error)
}

type AllowedUserLicenseValidator struct {
	inner                   LicenseValidator
	checker                 UserChecker
	ignoreAllowedUsersCheck bool
	allowedUsersInCache     map[string]struct{}
}

func NewAllowedUserLicenseValidator(inner LicenseValidator, checker UserChecker, ignoreCheck bool, staticAllowedUsers []string) *AllowedUserLicenseValidator {
	staticSet := make(map[string]struct{}, len(staticAllowedUsers))
	for _, u := range staticAllowedUsers {
		staticSet[strings.ToLower(strings.TrimSpace(u))] = struct{}{}
	}

	return &AllowedUserLicenseValidator{
		inner:                   inner,
		checker:                 checker,
		ignoreAllowedUsersCheck: ignoreCheck,
		allowedUsersInCache:     staticSet,
	}
}

func (v *AllowedUserLicenseValidator) ValidateReadAccess(ctx context.Context, req LicenseRequest) error {
	if !v.ignoreAllowedUsersCheck {
		principal, ok := auth.PrincipalFromContext(ctx)
		if !ok || principal.Subject == "" {
			return apperr.New(apperr.Unauthorized, "missing user identity")
		}

		userID := strings.ToLower(strings.TrimSpace(principal.Subject))
		
		// 1. Check static list from config
		_, staticallyAllowed := v.allowedUsersInCache[userID]
		
		// 2. Check Redis cache
		dynamicallyAllowed := false
		if !staticallyAllowed && v.checker != nil {
			var err error
			dynamicallyAllowed, err = v.checker.IsUserAllowed(ctx, userID)
			if err != nil {
				return apperr.Wrap(apperr.Unavailable, "check allowed user cache", err)
			}
		}

		if !staticallyAllowed && !dynamicallyAllowed {
			return apperr.New(apperr.Forbidden, "user is not in the allowed users list")
		}
	}

	// Fallback to the real inner validator (HTTP MDO check)
	return v.inner.ValidateReadAccess(ctx, req)
}
