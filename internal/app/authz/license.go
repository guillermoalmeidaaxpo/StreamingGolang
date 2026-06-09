package authz

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	client        *http.Client
	baseURL       string
	authorizePath string
	universePath  string
}

func NewHttpLicenseValidator(baseURL, authorizePath, universePath string, timeout time.Duration) *HttpLicenseValidator {
	// For dev environments pointing to internal Axpo servers, skip TLS verification if needed
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &HttpLicenseValidator{
		client:        &http.Client{Timeout: timeout, Transport: tr},
		baseURL:       strings.TrimSpace(baseURL),
		authorizePath: strings.TrimSpace(authorizePath),
		universePath:  strings.TrimSpace(universePath),
	}
}

type bulkDataUniverseRequest struct {
	Type                  string  `json:"type"`
	UniverseName          *string `json:"universeName"`
	Action                string  `json:"action"`
	InternalCorrelationID string  `json:"internalCorrelationId,omitempty"`
	MDOIDs                []int64 `json:"mdoIds"`
	StageID               uint8   `json:"stageId"`
}

type timeSeriesRequest struct {
	Identifiers           []int64 `json:"identifiers"`
	StageID               uint8   `json:"stageId"`
	InternalCorrelationID string  `json:"internalCorrelationId,omitempty"`
}

func (v *HttpLicenseValidator) ValidateReadAccess(ctx context.Context, req LicenseRequest) error {
	if v.baseURL == "" || v.baseURL == "NOT SET" {
		return nil
	}

	stageID := stageID(req.Stage)
	universeReq := bulkDataUniverseRequest{
		Type:                  "TransactionalDataOutbound",
		UniverseName:          nil,
		Action:                "Read",
		InternalCorrelationID: req.InternalCorrelationID,
		MDOIDs:                req.Identifiers,
		StageID:               stageID,
	}

	universeURL, err := joinURL(v.baseURL, v.universePath)
	if err != nil {
		return apperr.Wrap(apperr.Internal, "build universe authorization URL", err)
	}
	universeResp, err := v.postJSON(ctx, universeURL, universeReq, req.InternalCorrelationID)
	if err != nil {
		return err
	}
	defer universeResp.Body.Close()

	if universeResp.StatusCode == http.StatusInternalServerError {
		return apperr.New(apperr.Unavailable, "An error occurred while validating the data universe")
	}
	if universeResp.StatusCode == http.StatusUnauthorized {
		return apperr.New(apperr.Forbidden, "Unauthorized")
	}
	if universeResp.StatusCode < 200 || universeResp.StatusCode > 299 {
		return apperr.New(apperr.Unavailable, fmt.Sprintf("authorization universe service returned status %d", universeResp.StatusCode))
	}

	var permissions map[string]bool
	if err := json.NewDecoder(universeResp.Body).Decode(&permissions); err != nil {
		return apperr.Wrap(apperr.Unavailable, "decode universe authorization response", err)
	}

	for _, id := range req.Identifiers {
		if !permissions[strconv.FormatInt(id, 10)] {
			return apperr.New(apperr.Forbidden, fmt.Sprintf("You don't have access to MDO id: %d.", id))
		}
	}

	licenseReq := timeSeriesRequest{
		Identifiers:           req.Identifiers,
		StageID:               stageID,
		InternalCorrelationID: req.InternalCorrelationID,
	}

	authorizeURL, err := joinURL(v.baseURL, v.authorizePath)
	if err != nil {
		return apperr.Wrap(apperr.Internal, "build license authorization URL", err)
	}
	licenseResp, err := v.postJSON(ctx, authorizeURL, licenseReq, req.InternalCorrelationID)
	if err != nil {
		return err
	}
	defer licenseResp.Body.Close()

	if licenseResp.StatusCode >= 200 && licenseResp.StatusCode <= 299 {
		return nil
	}

	if licenseResp.StatusCode == http.StatusForbidden || licenseResp.StatusCode == http.StatusUnauthorized {
		return apperr.New(apperr.Forbidden, "Access denied: Missing license for one or more identifiers")
	}

	return apperr.New(apperr.Unavailable, fmt.Sprintf("authorization license service returned status %d", licenseResp.StatusCode))
}

func (v *HttpLicenseValidator) postJSON(ctx context.Context, endpoint string, payload any, correlationID string) (*http.Response, error) {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, apperr.Wrap(apperr.Internal, "marshal authorization request", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, apperr.Wrap(apperr.Internal, "create authorization request", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if correlationID != "" {
		httpReq.Header.Set("X-Correlation-ID", correlationID)
	}

	if rawToken, ok := ctx.Value("raw_bearer_token").(string); ok && strings.TrimSpace(rawToken) != "" {
		httpReq.Header.Set("Authorization", "Bearer "+rawToken)
	}

	resp, err := v.client.Do(httpReq)
	if err != nil {
		return nil, apperr.Wrap(apperr.Unavailable, fmt.Sprintf("call authorization service %s", endpoint), err)
	}
	return resp, nil
}

func joinURL(baseURL, path string) (string, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	path = strings.Trim(strings.TrimSpace(path), "/")
	if baseURL == "" || path == "" {
		return "", fmt.Errorf("base URL and path are required")
	}

	joined := baseURL + "/" + path
	parsed, err := url.ParseRequestURI(joined)
	if err != nil {
		return "", err
	}
	return parsed.String(), nil
}

func stageID(stage string) uint8 {
	switch strings.ToLower(strings.TrimSpace(stage)) {
	case "design":
		return 1
	case "validation":
		return 2
	case "productive", "production", "prod", "migration", "development", "dev", "":
		return 3
	default:
		return 3
	}
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
