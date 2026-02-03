// Package billingio provides a Go client for the billing.io API.
//
// Create a client with your API key:
//
//	client := billingio.New("sk_live_...")
//
// All service methods accept a context.Context as their first parameter:
//
//	checkout, err := client.Checkouts.Create(ctx, &billingio.CreateCheckoutParams{
//	    AmountUSD: 49.99,
//	    Chain:     billingio.ChainTron,
//	    Token:     billingio.TokenUSDT,
//	})
package billingio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
)

const (
	defaultBaseURL = "https://api.billing.io/v1"
	sdkVersion     = "0.1.0"
)

// Client is the billing.io API client. Use New to create one.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	userAgent  string

	// Resource services
	Checkouts            *CheckoutService
	Webhooks             *WebhookService
	Events               *EventService
	Health               *HealthService
	Customers            *CustomerService
	PaymentMethods       *PaymentMethodService
	PaymentLinks         *PaymentLinkService
	SubscriptionPlans    *SubscriptionPlanService
	Subscriptions        *SubscriptionService
	SubscriptionRenewals *SubscriptionRenewalService
	Entitlements         *EntitlementService
	Payouts              *PayoutService
	Settlements          *SettlementService
	RevenueEvents        *RevenueEventService
	Adjustments          *AdjustmentService
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithHTTPClient sets a custom *http.Client for all requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// New creates a new billing.io API client.
//
//	client := billingio.New("sk_live_...",
//	    billingio.WithBaseURL("https://api.billing.io/v1"),
//	    billingio.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
//	)
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
		userAgent:  fmt.Sprintf("billing-go/%s Go/%s", sdkVersion, runtime.Version()),
	}
	for _, opt := range opts {
		opt(c)
	}

	c.Checkouts = &CheckoutService{client: c}
	c.Webhooks = &WebhookService{client: c}
	c.Events = &EventService{client: c}
	c.Health = &HealthService{client: c}
	c.Customers = &CustomerService{client: c}
	c.PaymentMethods = &PaymentMethodService{client: c}
	c.PaymentLinks = &PaymentLinkService{client: c}
	c.SubscriptionPlans = &SubscriptionPlanService{client: c}
	c.Subscriptions = &SubscriptionService{client: c}
	c.SubscriptionRenewals = &SubscriptionRenewalService{client: c}
	c.Entitlements = &EntitlementService{client: c}
	c.Payouts = &PayoutService{client: c}
	c.Settlements = &SettlementService{client: c}
	c.RevenueEvents = &RevenueEventService{client: c}
	c.Adjustments = &AdjustmentService{client: c}

	return c
}

// do executes an HTTP request and decodes the response into dest.
// If dest is nil the response body is discarded (used for 204 responses).
func (c *Client) do(ctx context.Context, method, path string, body any, dest any, headers map[string]string) error {
	u := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("billingio: failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return fmt.Errorf("billingio: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("billingio: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("billingio: failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	if dest != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("billingio: failed to decode response: %w", err)
		}
	}

	return nil
}

// get is a convenience wrapper for GET requests.
func (c *Client) get(ctx context.Context, path string, dest any) error {
	return c.do(ctx, http.MethodGet, path, nil, dest, nil)
}

// post is a convenience wrapper for POST requests.
func (c *Client) post(ctx context.Context, path string, body any, dest any, headers map[string]string) error {
	return c.do(ctx, http.MethodPost, path, body, dest, headers)
}

// patch is a convenience wrapper for PATCH requests.
func (c *Client) patch(ctx context.Context, path string, body any, dest any) error {
	return c.do(ctx, http.MethodPatch, path, body, dest, nil)
}

// del is a convenience wrapper for DELETE requests.
func (c *Client) del(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// parseAPIError decodes an error response body into an *Error.
func parseAPIError(statusCode int, body []byte) error {
	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err != nil || errResp.Err == nil {
		return &Error{
			Type:       "internal_error",
			Code:       "unknown",
			StatusCode: statusCode,
			Message:    fmt.Sprintf("unexpected error (HTTP %d): %s", statusCode, string(body)),
		}
	}
	errResp.Err.StatusCode = statusCode
	return errResp.Err
}

// addQueryParams builds a URL path with optional query parameters.
func addQueryParams(basePath string, params map[string]string) string {
	if len(params) == 0 {
		return basePath
	}
	v := url.Values{}
	for key, val := range params {
		if val != "" {
			v.Set(key, val)
		}
	}
	encoded := v.Encode()
	if encoded == "" {
		return basePath
	}
	return basePath + "?" + encoded
}

// intToString converts an *int to its string representation, or returns "".
func intToString(v *int) string {
	if v == nil {
		return ""
	}
	return strconv.Itoa(*v)
}

// strOrEmpty dereferences a *string or returns "".
func strOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
