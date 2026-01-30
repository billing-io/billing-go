package billingio

// Chain represents a supported blockchain network.
type Chain string

const (
	ChainTron     Chain = "tron"
	ChainArbitrum Chain = "arbitrum"
)

// Token represents a supported stablecoin.
type Token string

const (
	TokenUSDT Token = "USDT"
	TokenUSDC Token = "USDC"
)

// CheckoutStatus represents the lifecycle state of a checkout.
type CheckoutStatus string

const (
	CheckoutStatusPending    CheckoutStatus = "pending"
	CheckoutStatusDetected   CheckoutStatus = "detected"
	CheckoutStatusConfirming CheckoutStatus = "confirming"
	CheckoutStatusConfirmed  CheckoutStatus = "confirmed"
	CheckoutStatusExpired    CheckoutStatus = "expired"
	CheckoutStatusFailed     CheckoutStatus = "failed"
)

// EventType represents a webhook event type.
type EventType string

const (
	EventTypeCheckoutCreated         EventType = "checkout.created"
	EventTypeCheckoutPaymentDetected EventType = "checkout.payment_detected"
	EventTypeCheckoutConfirming      EventType = "checkout.confirming"
	EventTypeCheckoutCompleted       EventType = "checkout.completed"
	EventTypeCheckoutExpired         EventType = "checkout.expired"
	EventTypeCheckoutFailed          EventType = "checkout.failed"
)

// WebhookEndpointStatus represents the status of a webhook endpoint.
type WebhookEndpointStatus string

const (
	WebhookEndpointStatusActive   WebhookEndpointStatus = "active"
	WebhookEndpointStatusDisabled WebhookEndpointStatus = "disabled"
)

// Checkout represents a crypto payment checkout.
type Checkout struct {
	CheckoutID            string            `json:"checkout_id"`
	DepositAddress        string            `json:"deposit_address"`
	Chain                 Chain             `json:"chain"`
	Token                 Token             `json:"token"`
	AmountUSD             float64           `json:"amount_usd"`
	AmountAtomic          string            `json:"amount_atomic"`
	Status                CheckoutStatus    `json:"status"`
	TxHash                *string           `json:"tx_hash"`
	Confirmations         int               `json:"confirmations"`
	RequiredConfirmations int               `json:"required_confirmations"`
	ExpiresAt             string            `json:"expires_at"`
	DetectedAt            *string           `json:"detected_at"`
	ConfirmedAt           *string           `json:"confirmed_at"`
	CreatedAt             string            `json:"created_at"`
	Metadata              map[string]string `json:"metadata,omitempty"`
}

// CheckoutStatusResponse is the lightweight status polling response.
type CheckoutStatusResponse struct {
	CheckoutID            string         `json:"checkout_id"`
	Status                CheckoutStatus `json:"status"`
	TxHash                *string        `json:"tx_hash"`
	Confirmations         int            `json:"confirmations"`
	RequiredConfirmations int            `json:"required_confirmations"`
	DetectedAt            *string        `json:"detected_at"`
	ConfirmedAt           *string        `json:"confirmed_at"`
	PollingIntervalMs     int            `json:"polling_interval_ms"`
}

// CheckoutList is a paginated list of checkouts.
type CheckoutList struct {
	Data       []Checkout `json:"data"`
	HasMore    bool       `json:"has_more"`
	NextCursor *string    `json:"next_cursor"`
}

// WebhookEndpoint represents a registered webhook endpoint.
type WebhookEndpoint struct {
	WebhookID   string                `json:"webhook_id"`
	URL         string                `json:"url"`
	Events      []EventType           `json:"events"`
	Secret      string                `json:"secret,omitempty"`
	Description *string               `json:"description"`
	Status      WebhookEndpointStatus `json:"status"`
	CreatedAt   string                `json:"created_at"`
}

// WebhookEndpointList is a paginated list of webhook endpoints.
type WebhookEndpointList struct {
	Data       []WebhookEndpoint `json:"data"`
	HasMore    bool              `json:"has_more"`
	NextCursor *string           `json:"next_cursor"`
}

// Event represents a webhook event.
type Event struct {
	EventID    string    `json:"event_id"`
	Type       EventType `json:"type"`
	CheckoutID string    `json:"checkout_id"`
	Data       Checkout  `json:"data"`
	CreatedAt  string    `json:"created_at"`
}

// EventList is a paginated list of events.
type EventList struct {
	Data       []Event `json:"data"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

// HealthResponse is returned by the health check endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// WebhookEvent is the parsed payload of an incoming webhook delivery.
// It is returned by VerifyWebhookSignature after successful verification.
type WebhookEvent struct {
	EventID    string    `json:"event_id"`
	Type       EventType `json:"type"`
	CheckoutID string    `json:"checkout_id"`
	Data       Checkout  `json:"data"`
	CreatedAt  string    `json:"created_at"`
}

// CreateCheckoutParams are the parameters for creating a checkout.
type CreateCheckoutParams struct {
	AmountUSD        float64           `json:"amount_usd"`
	Chain            Chain             `json:"chain"`
	Token            Token             `json:"token"`
	ExpiresInSeconds *int              `json:"expires_in_seconds,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`

	// IdempotencyKey is sent as the Idempotency-Key header. Optional.
	IdempotencyKey string `json:"-"`
}

// ListCheckoutsParams are the parameters for listing checkouts.
type ListCheckoutsParams struct {
	Cursor *string         `json:"cursor,omitempty"`
	Limit  *int            `json:"limit,omitempty"`
	Status *CheckoutStatus `json:"status,omitempty"`
}

// CreateWebhookParams are the parameters for creating a webhook endpoint.
type CreateWebhookParams struct {
	URL         string      `json:"url"`
	Events      []EventType `json:"events"`
	Description *string     `json:"description,omitempty"`
}

// ListParams are generic pagination parameters for list endpoints.
type ListParams struct {
	Cursor *string `json:"cursor,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

// ListEventsParams are the parameters for listing events.
type ListEventsParams struct {
	Cursor     *string    `json:"cursor,omitempty"`
	Limit      *int       `json:"limit,omitempty"`
	Type       *EventType `json:"type,omitempty"`
	CheckoutID *string    `json:"checkout_id,omitempty"`
}
