package coreclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("rootguard core returned %d: %s", e.StatusCode, e.Message)
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http:    &http.Client{Timeout: 20 * time.Second},
	}
}

type Dashboard struct {
	Docker struct {
		CPU        float64 `json:"cpu"`
		Memory     float64 `json:"memory"`
		Containers int     `json:"containers"`
		Status     string  `json:"status"`
	} `json:"docker"`
	DNS struct {
		Status   string `json:"status"`
		Resolver string `json:"resolver"`
		DNSSEC   bool   `json:"dnssec"`
	} `json:"dns"`
}

type Service struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ServiceActionResponse struct {
	Service string `json:"service"`
	Action  string `json:"action"`
	Status  string `json:"status"`
}

type UnboundSettings struct {
	QnameMinimisation bool `json:"qname_minimisation"`
	Prefetch          bool `json:"prefetch"`
	ServeExpired      bool `json:"serve_expired"`
	CacheMinTTL       int  `json:"cache_min_ttl"`
	CacheMaxTTL       int  `json:"cache_max_ttl"`
	Threads           int  `json:"threads"`
}

type UnboundChange struct {
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

type UnboundPreview struct {
	Changed        bool            `json:"changed"`
	Changes        []UnboundChange `json:"changes"`
	RenderedConfig string          `json:"rendered_config"`
}

type UnboundHistoryEntry struct {
	ID        string          `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	Settings  UnboundSettings `json:"settings"`
	Config    string          `json:"config,omitempty"`
}

type UnboundDiagnosticCheck struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
	Detail string `json:"detail"`
}

type UnboundDiagnosticReport struct {
	Healthy   bool                     `json:"healthy"`
	CheckedAt time.Time                `json:"checked_at"`
	Checks    []UnboundDiagnosticCheck `json:"checks"`
}

type AdGuardStatus struct {
	Configured    bool   `json:"configured"`
	Healthy       bool   `json:"healthy"`
	Upstream      string `json:"upstream"`
	UpstreamReady bool   `json:"upstream_ready"`
}

func (c *Client) Dashboard(ctx context.Context) (Dashboard, error) {
	var result Dashboard
	err := c.do(ctx, http.MethodGet, "/api/dashboard", nil, &result)
	return result, err
}

func (c *Client) System(ctx context.Context) (map[string]string, error) {
	var result map[string]string
	err := c.do(ctx, http.MethodGet, "/api/system", nil, &result)
	return result, err
}

func (c *Client) Services(ctx context.Context) ([]Service, error) {
	var result []Service
	err := c.do(ctx, http.MethodGet, "/api/services", nil, &result)
	return result, err
}

func (c *Client) ServiceAction(ctx context.Context, service, action string) (ServiceActionResponse, error) {
	var result ServiceActionResponse
	err := c.do(ctx, http.MethodPost, "/api/services/"+service+"/"+action, nil, &result)
	return result, err
}

func (c *Client) UnboundSettings(ctx context.Context) (UnboundSettings, error) {
	var result UnboundSettings
	err := c.do(ctx, http.MethodGet, "/api/unbound/settings", nil, &result)
	return result, err
}

func (c *Client) UpdateUnboundSettings(ctx context.Context, settings UnboundSettings) (UnboundSettings, error) {
	var result UnboundSettings
	err := c.do(ctx, http.MethodPut, "/api/unbound/settings", settings, &result)
	return result, err
}

func (c *Client) PreviewUnboundSettings(ctx context.Context, settings UnboundSettings) (UnboundPreview, error) {
	var result UnboundPreview
	err := c.do(ctx, http.MethodPost, "/api/unbound/preview", settings, &result)
	return result, err
}

func (c *Client) UnboundHistory(ctx context.Context) ([]UnboundHistoryEntry, error) {
	var result []UnboundHistoryEntry
	err := c.do(ctx, http.MethodGet, "/api/unbound/history", nil, &result)
	return result, err
}

func (c *Client) RestoreUnboundVersion(ctx context.Context, id string) (UnboundSettings, error) {
	var result UnboundSettings
	err := c.do(ctx, http.MethodPost, "/api/unbound/history/"+id+"/restore", nil, &result)
	return result, err
}

func (c *Client) UnboundDiagnostics(ctx context.Context) (UnboundDiagnosticReport, error) {
	var result UnboundDiagnosticReport
	err := c.do(ctx, http.MethodGet, "/api/unbound/diagnostics", nil, &result)
	return result, err
}

func (c *Client) AdGuardStatus(ctx context.Context) (AdGuardStatus, error) {
	var result AdGuardStatus
	err := c.do(ctx, http.MethodGet, "/api/adguard/status", nil, &result)
	return result, err
}

func (c *Client) BootstrapAdGuard(ctx context.Context) (AdGuardStatus, error) {
	var result AdGuardStatus
	err := c.do(ctx, http.MethodPost, "/api/adguard/bootstrap", nil, &result)
	return result, err
}

func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
	var requestBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		requestBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("rootguard core request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(string(message)),
		}
	}
	if result == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode rootguard core response: %w", err)
	}
	return nil
}
