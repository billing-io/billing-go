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

// ---------------------------------------------------------------------------
// Customers
// ---------------------------------------------------------------------------

// CustomerStatus represents the status of a customer.
type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusArchived CustomerStatus = "archived"
)

// Customer represents a billing.io customer.
type Customer struct {
	CustomerID string            `json:"customer_id"`
	Email      string            `json:"email"`
	Name       *string           `json:"name"`
	Status     CustomerStatus    `json:"status"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  string            `json:"created_at"`
	UpdatedAt  string            `json:"updated_at"`
}

// CustomerList is a paginated list of customers.
type CustomerList struct {
	Data       []Customer `json:"data"`
	HasMore    bool       `json:"has_more"`
	NextCursor *string    `json:"next_cursor"`
}

// CreateCustomerParams are the parameters for creating a customer.
type CreateCustomerParams struct {
	Email    string            `json:"email"`
	Name     *string           `json:"name,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateCustomerParams are the parameters for updating a customer.
type UpdateCustomerParams struct {
	Email    *string           `json:"email,omitempty"`
	Name     *string           `json:"name,omitempty"`
	Status   *CustomerStatus   `json:"status,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ListCustomersParams are the parameters for listing customers.
type ListCustomersParams struct {
	Cursor *string         `json:"cursor,omitempty"`
	Limit  *int            `json:"limit,omitempty"`
	Status *CustomerStatus `json:"status,omitempty"`
}

// ---------------------------------------------------------------------------
// Payment Methods
// ---------------------------------------------------------------------------

// PaymentMethodType represents the type of a payment method.
type PaymentMethodType string

const (
	PaymentMethodTypeWallet PaymentMethodType = "wallet"
)

// PaymentMethodStatus represents the status of a payment method.
type PaymentMethodStatus string

const (
	PaymentMethodStatusActive   PaymentMethodStatus = "active"
	PaymentMethodStatusDisabled PaymentMethodStatus = "disabled"
)

// PaymentMethod represents a stored payment method.
type PaymentMethod struct {
	PaymentMethodID string              `json:"payment_method_id"`
	CustomerID      string              `json:"customer_id"`
	Type            PaymentMethodType   `json:"type"`
	Chain           Chain               `json:"chain"`
	WalletAddress   string              `json:"wallet_address"`
	IsDefault       bool                `json:"is_default"`
	Status          PaymentMethodStatus `json:"status"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
}

// PaymentMethodList is a paginated list of payment methods.
type PaymentMethodList struct {
	Data       []PaymentMethod `json:"data"`
	HasMore    bool            `json:"has_more"`
	NextCursor *string         `json:"next_cursor"`
}

// CreatePaymentMethodParams are the parameters for creating a payment method.
type CreatePaymentMethodParams struct {
	CustomerID    string            `json:"customer_id"`
	Type          PaymentMethodType `json:"type"`
	Chain         Chain             `json:"chain"`
	WalletAddress string            `json:"wallet_address"`
}

// UpdatePaymentMethodParams are the parameters for updating a payment method.
type UpdatePaymentMethodParams struct {
	Status *PaymentMethodStatus `json:"status,omitempty"`
}

// ListPaymentMethodsParams are the parameters for listing payment methods.
type ListPaymentMethodsParams struct {
	Cursor     *string `json:"cursor,omitempty"`
	Limit      *int    `json:"limit,omitempty"`
	CustomerID *string `json:"customer_id,omitempty"`
}

// ---------------------------------------------------------------------------
// Payment Links
// ---------------------------------------------------------------------------

// PaymentLinkStatus represents the status of a payment link.
type PaymentLinkStatus string

const (
	PaymentLinkStatusActive   PaymentLinkStatus = "active"
	PaymentLinkStatusInactive PaymentLinkStatus = "inactive"
)

// PaymentLink represents a reusable payment link.
type PaymentLink struct {
	PaymentLinkID string            `json:"payment_link_id"`
	URL           string            `json:"url"`
	AmountUSD     *float64          `json:"amount_usd"`
	Chain         *Chain            `json:"chain"`
	Token         *Token            `json:"token"`
	Description   *string           `json:"description"`
	Status        PaymentLinkStatus `json:"status"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	CreatedAt     string            `json:"created_at"`
}

// PaymentLinkList is a paginated list of payment links.
type PaymentLinkList struct {
	Data       []PaymentLink `json:"data"`
	HasMore    bool          `json:"has_more"`
	NextCursor *string       `json:"next_cursor"`
}

// CreatePaymentLinkParams are the parameters for creating a payment link.
type CreatePaymentLinkParams struct {
	AmountUSD   *float64          `json:"amount_usd,omitempty"`
	Chain       *Chain            `json:"chain,omitempty"`
	Token       *Token            `json:"token,omitempty"`
	Description *string           `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ListPaymentLinksParams are the parameters for listing payment links.
type ListPaymentLinksParams struct {
	Cursor *string `json:"cursor,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

// ---------------------------------------------------------------------------
// Subscription Plans
// ---------------------------------------------------------------------------

// BillingInterval represents the billing cycle of a subscription plan.
type BillingInterval string

const (
	BillingIntervalWeekly  BillingInterval = "weekly"
	BillingIntervalMonthly BillingInterval = "monthly"
	BillingIntervalYearly  BillingInterval = "yearly"
)

// SubscriptionPlanStatus represents the status of a subscription plan.
type SubscriptionPlanStatus string

const (
	SubscriptionPlanStatusActive   SubscriptionPlanStatus = "active"
	SubscriptionPlanStatusArchived SubscriptionPlanStatus = "archived"
)

// SubscriptionPlan represents a recurring billing plan.
type SubscriptionPlan struct {
	PlanID          string                 `json:"plan_id"`
	Name            string                 `json:"name"`
	Description     *string                `json:"description"`
	AmountUSD       float64                `json:"amount_usd"`
	BillingInterval BillingInterval        `json:"billing_interval"`
	Status          SubscriptionPlanStatus `json:"status"`
	Metadata        map[string]string      `json:"metadata,omitempty"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

// SubscriptionPlanList is a paginated list of subscription plans.
type SubscriptionPlanList struct {
	Data       []SubscriptionPlan `json:"data"`
	HasMore    bool               `json:"has_more"`
	NextCursor *string            `json:"next_cursor"`
}

// CreateSubscriptionPlanParams are the parameters for creating a subscription plan.
type CreateSubscriptionPlanParams struct {
	Name            string            `json:"name"`
	Description     *string           `json:"description,omitempty"`
	AmountUSD       float64           `json:"amount_usd"`
	BillingInterval BillingInterval   `json:"billing_interval"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// UpdateSubscriptionPlanParams are the parameters for updating a subscription plan.
type UpdateSubscriptionPlanParams struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Status      *SubscriptionPlanStatus `json:"status,omitempty"`
	Metadata    map[string]string       `json:"metadata,omitempty"`
}

// ListSubscriptionPlansParams are the parameters for listing subscription plans.
type ListSubscriptionPlansParams struct {
	Cursor *string                 `json:"cursor,omitempty"`
	Limit  *int                    `json:"limit,omitempty"`
	Status *SubscriptionPlanStatus `json:"status,omitempty"`
}

// ---------------------------------------------------------------------------
// Subscriptions
// ---------------------------------------------------------------------------

// SubscriptionStatus represents the lifecycle state of a subscription.
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPaused    SubscriptionStatus = "paused"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
)

// Subscription represents a customer subscription.
type Subscription struct {
	SubscriptionID  string             `json:"subscription_id"`
	CustomerID      string             `json:"customer_id"`
	PlanID          string             `json:"plan_id"`
	Status          SubscriptionStatus `json:"status"`
	CurrentPeriodStart string          `json:"current_period_start"`
	CurrentPeriodEnd   string          `json:"current_period_end"`
	CancelledAt     *string            `json:"cancelled_at"`
	Metadata        map[string]string  `json:"metadata,omitempty"`
	CreatedAt       string             `json:"created_at"`
	UpdatedAt       string             `json:"updated_at"`
}

// SubscriptionList is a paginated list of subscriptions.
type SubscriptionList struct {
	Data       []Subscription `json:"data"`
	HasMore    bool           `json:"has_more"`
	NextCursor *string        `json:"next_cursor"`
}

// CreateSubscriptionParams are the parameters for creating a subscription.
type CreateSubscriptionParams struct {
	CustomerID string            `json:"customer_id"`
	PlanID     string            `json:"plan_id"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// UpdateSubscriptionParams are the parameters for updating a subscription.
type UpdateSubscriptionParams struct {
	Status   *SubscriptionStatus `json:"status,omitempty"`
	Metadata map[string]string   `json:"metadata,omitempty"`
}

// ListSubscriptionsParams are the parameters for listing subscriptions.
type ListSubscriptionsParams struct {
	Cursor     *string             `json:"cursor,omitempty"`
	Limit      *int                `json:"limit,omitempty"`
	CustomerID *string             `json:"customer_id,omitempty"`
	PlanID     *string             `json:"plan_id,omitempty"`
	Status     *SubscriptionStatus `json:"status,omitempty"`
}

// ---------------------------------------------------------------------------
// Subscription Renewals
// ---------------------------------------------------------------------------

// RenewalStatus represents the status of a subscription renewal.
type RenewalStatus string

const (
	RenewalStatusPending   RenewalStatus = "pending"
	RenewalStatusPaid      RenewalStatus = "paid"
	RenewalStatusFailed    RenewalStatus = "failed"
	RenewalStatusRetrying  RenewalStatus = "retrying"
)

// SubscriptionRenewal represents a single billing-cycle renewal.
type SubscriptionRenewal struct {
	RenewalID      string        `json:"renewal_id"`
	SubscriptionID string        `json:"subscription_id"`
	PlanID         string        `json:"plan_id"`
	AmountUSD      float64       `json:"amount_usd"`
	Status         RenewalStatus `json:"status"`
	PeriodStart    string        `json:"period_start"`
	PeriodEnd      string        `json:"period_end"`
	PaidAt         *string       `json:"paid_at"`
	FailedAt       *string       `json:"failed_at"`
	CreatedAt      string        `json:"created_at"`
}

// SubscriptionRenewalList is a paginated list of subscription renewals.
type SubscriptionRenewalList struct {
	Data       []SubscriptionRenewal `json:"data"`
	HasMore    bool                  `json:"has_more"`
	NextCursor *string               `json:"next_cursor"`
}

// ListSubscriptionRenewalsParams are the parameters for listing subscription renewals.
type ListSubscriptionRenewalsParams struct {
	Cursor         *string        `json:"cursor,omitempty"`
	Limit          *int           `json:"limit,omitempty"`
	SubscriptionID *string        `json:"subscription_id,omitempty"`
	Status         *RenewalStatus `json:"status,omitempty"`
}

// ---------------------------------------------------------------------------
// Entitlements
// ---------------------------------------------------------------------------

// Entitlement represents a feature or access grant tied to a subscription.
type Entitlement struct {
	EntitlementID  string            `json:"entitlement_id"`
	SubscriptionID string            `json:"subscription_id"`
	FeatureKey     string            `json:"feature_key"`
	Value          string            `json:"value"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
}

// EntitlementList is a paginated list of entitlements.
type EntitlementList struct {
	Data       []Entitlement `json:"data"`
	HasMore    bool          `json:"has_more"`
	NextCursor *string       `json:"next_cursor"`
}

// EntitlementCheckResponse is the response from the entitlement check endpoint.
type EntitlementCheckResponse struct {
	Entitled bool   `json:"entitled"`
	Value    string `json:"value"`
}

// CreateEntitlementParams are the parameters for creating an entitlement.
type CreateEntitlementParams struct {
	SubscriptionID string            `json:"subscription_id"`
	FeatureKey     string            `json:"feature_key"`
	Value          string            `json:"value"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// UpdateEntitlementParams are the parameters for updating an entitlement.
type UpdateEntitlementParams struct {
	Value    *string           `json:"value,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ListEntitlementsParams are the parameters for listing entitlements.
type ListEntitlementsParams struct {
	Cursor         *string `json:"cursor,omitempty"`
	Limit          *int    `json:"limit,omitempty"`
	SubscriptionID *string `json:"subscription_id,omitempty"`
	FeatureKey     *string `json:"feature_key,omitempty"`
}

// CheckEntitlementParams are the parameters for checking an entitlement.
type CheckEntitlementParams struct {
	CustomerID string `json:"customer_id"`
	FeatureKey string `json:"feature_key"`
}

// ---------------------------------------------------------------------------
// Payout Intents
// ---------------------------------------------------------------------------

// PayoutStatus represents the lifecycle state of a payout.
type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "pending"
	PayoutStatusExecuting PayoutStatus = "executing"
	PayoutStatusCompleted PayoutStatus = "completed"
	PayoutStatusFailed    PayoutStatus = "failed"
)

// Payout represents a payout intent.
type Payout struct {
	PayoutID      string            `json:"payout_id"`
	AmountUSD     float64           `json:"amount_usd"`
	Chain         Chain             `json:"chain"`
	Token         Token             `json:"token"`
	WalletAddress string            `json:"wallet_address"`
	TxHash        *string           `json:"tx_hash"`
	Status        PayoutStatus      `json:"status"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	ExecutedAt    *string           `json:"executed_at"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

// PayoutList is a paginated list of payouts.
type PayoutList struct {
	Data       []Payout `json:"data"`
	HasMore    bool     `json:"has_more"`
	NextCursor *string  `json:"next_cursor"`
}

// CreatePayoutParams are the parameters for creating a payout.
type CreatePayoutParams struct {
	AmountUSD     float64           `json:"amount_usd"`
	Chain         Chain             `json:"chain"`
	Token         Token             `json:"token"`
	WalletAddress string            `json:"wallet_address"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// UpdatePayoutParams are the parameters for updating a payout.
type UpdatePayoutParams struct {
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ListPayoutsParams are the parameters for listing payouts.
type ListPayoutsParams struct {
	Cursor *string       `json:"cursor,omitempty"`
	Limit  *int          `json:"limit,omitempty"`
	Status *PayoutStatus `json:"status,omitempty"`
}

// ---------------------------------------------------------------------------
// Settlements
// ---------------------------------------------------------------------------

// Settlement represents a completed settlement record.
type Settlement struct {
	SettlementID string  `json:"settlement_id"`
	PayoutID     string  `json:"payout_id"`
	AmountUSD    float64 `json:"amount_usd"`
	Chain        Chain   `json:"chain"`
	Token        Token   `json:"token"`
	TxHash       string  `json:"tx_hash"`
	SettledAt    string  `json:"settled_at"`
	CreatedAt    string  `json:"created_at"`
}

// SettlementList is a paginated list of settlements.
type SettlementList struct {
	Data       []Settlement `json:"data"`
	HasMore    bool         `json:"has_more"`
	NextCursor *string      `json:"next_cursor"`
}

// ListSettlementsParams are the parameters for listing settlements.
type ListSettlementsParams struct {
	Cursor   *string `json:"cursor,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
	PayoutID *string `json:"payout_id,omitempty"`
}

// ---------------------------------------------------------------------------
// Revenue Events
// ---------------------------------------------------------------------------

// RevenueEventType represents the type of a revenue event.
type RevenueEventType string

const (
	RevenueEventTypeCharge     RevenueEventType = "charge"
	RevenueEventTypeRefund     RevenueEventType = "refund"
	RevenueEventTypeAdjustment RevenueEventType = "adjustment"
)

// RevenueEvent represents a revenue event for accounting.
type RevenueEvent struct {
	RevenueEventID string           `json:"revenue_event_id"`
	Type           RevenueEventType `json:"type"`
	AmountUSD      float64          `json:"amount_usd"`
	CustomerID     *string          `json:"customer_id"`
	CheckoutID     *string          `json:"checkout_id"`
	SubscriptionID *string          `json:"subscription_id"`
	Description    *string          `json:"description"`
	CreatedAt      string           `json:"created_at"`
}

// RevenueEventList is a paginated list of revenue events.
type RevenueEventList struct {
	Data       []RevenueEvent `json:"data"`
	HasMore    bool           `json:"has_more"`
	NextCursor *string        `json:"next_cursor"`
}

// ListRevenueEventsParams are the parameters for listing revenue events.
type ListRevenueEventsParams struct {
	Cursor *string           `json:"cursor,omitempty"`
	Limit  *int              `json:"limit,omitempty"`
	Type   *RevenueEventType `json:"type,omitempty"`
}

// AccountingSummary is the response from the revenue accounting endpoint.
type AccountingSummary struct {
	TotalRevenueUSD  float64 `json:"total_revenue_usd"`
	TotalRefundsUSD  float64 `json:"total_refunds_usd"`
	NetRevenueUSD    float64 `json:"net_revenue_usd"`
	TotalCharges     int     `json:"total_charges"`
	TotalRefunds     int     `json:"total_refunds"`
	PeriodStart      string  `json:"period_start"`
	PeriodEnd        string  `json:"period_end"`
}

// AccountingSummaryParams are the parameters for querying the accounting summary.
type AccountingSummaryParams struct {
	PeriodStart *string `json:"period_start,omitempty"`
	PeriodEnd   *string `json:"period_end,omitempty"`
}

// ---------------------------------------------------------------------------
// Adjustments
// ---------------------------------------------------------------------------

// AdjustmentType represents the type of a revenue adjustment.
type AdjustmentType string

const (
	AdjustmentTypeCredit AdjustmentType = "credit"
	AdjustmentTypeDebit  AdjustmentType = "debit"
)

// Adjustment represents a revenue adjustment (credit or debit).
type Adjustment struct {
	AdjustmentID string            `json:"adjustment_id"`
	Type         AdjustmentType    `json:"type"`
	AmountUSD    float64           `json:"amount_usd"`
	CustomerID   *string           `json:"customer_id"`
	Description  *string           `json:"description"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CreatedAt    string            `json:"created_at"`
}

// AdjustmentList is a paginated list of adjustments.
type AdjustmentList struct {
	Data       []Adjustment `json:"data"`
	HasMore    bool         `json:"has_more"`
	NextCursor *string      `json:"next_cursor"`
}

// CreateAdjustmentParams are the parameters for creating an adjustment.
type CreateAdjustmentParams struct {
	Type        AdjustmentType    `json:"type"`
	AmountUSD   float64           `json:"amount_usd"`
	CustomerID  *string           `json:"customer_id,omitempty"`
	Description *string           `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ListAdjustmentsParams are the parameters for listing adjustments.
type ListAdjustmentsParams struct {
	Cursor *string         `json:"cursor,omitempty"`
	Limit  *int            `json:"limit,omitempty"`
	Type   *AdjustmentType `json:"type,omitempty"`
}
