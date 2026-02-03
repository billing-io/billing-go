# billing-go

Official Go SDK for the [billing.io](https://billing.io) API -- non-custodial crypto checkout infrastructure.

## Installation

```bash
go get github.com/billing-io/billing-go
```

Requires **Go 1.21+**. Zero external dependencies.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	billingio "github.com/billing-io/billing-go"
)

func main() {
	client := billingio.New("sk_live_...")

	ctx := context.Background()

	// Create a checkout
	checkout, err := client.Checkouts.Create(ctx, &billingio.CreateCheckoutParams{
		AmountUSD: 49.99,
		Chain:     billingio.ChainTron,
		Token:     billingio.TokenUSDT,
		Metadata: map[string]string{
			"order_id": "ord_12345",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Checkout %s created\n", checkout.CheckoutID)
	fmt.Printf("Deposit to: %s\n", checkout.DepositAddress)
	fmt.Printf("Amount: %s %s\n", checkout.AmountAtomic, checkout.Token)

	// Poll for status
	status, err := client.Checkouts.GetStatus(ctx, checkout.CheckoutID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %s (%d/%d confirmations)\n",
		status.Status, status.Confirmations, status.RequiredConfirmations)
}
```

## Configuration

```go
import (
	"net/http"
	"time"
)

client := billingio.New("sk_live_...",
	billingio.WithBaseURL("https://api.billing.io/v1"),
	billingio.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)
```

## Webhook verification

Verify incoming webhook signatures in a standard `net/http` handler:

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	billingio "github.com/billing-io/billing-go"
)

func main() {
	secret := os.Getenv("BILLING_WEBHOOK_SECRET") // whsec_...

	http.HandleFunc("/webhooks/billing", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		sig := r.Header.Get(billingio.SignatureHeader)

		event, err := billingio.VerifyWebhookSignature(body, sig, secret)
		if err != nil {
			log.Printf("webhook verification failed: %v", err)
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}

		switch event.Type {
		case billingio.EventTypeCheckoutCompleted:
			fmt.Printf("Payment confirmed for checkout %s\n", event.CheckoutID)
			// Fulfill the order...

		case billingio.EventTypeCheckoutExpired:
			fmt.Printf("Checkout %s expired\n", event.CheckoutID)

		case billingio.EventTypeCheckoutFailed:
			fmt.Printf("Checkout %s failed\n", event.CheckoutID)
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

You can also set a custom timestamp tolerance (in seconds):

```go
event, err := billingio.VerifyWebhookSignatureWithTolerance(body, sig, secret, 600)
```

## Error handling

All API errors are returned as `*billingio.Error` values. Use the helper
functions or type-assert to inspect them:

```go
import "errors"

checkout, err := client.Checkouts.Get(ctx, "co_nonexistent")
if err != nil {
	if billingio.IsNotFound(err) {
		fmt.Println("Checkout not found")
	} else if billingio.IsRateLimited(err) {
		fmt.Println("Slow down -- rate limited")
	} else if billingio.IsAuthError(err) {
		fmt.Println("Check your API key")
	} else {
		// Inspect the full error
		var apiErr *billingio.Error
		if errors.As(err, &apiErr) {
			fmt.Printf("API error: type=%s code=%s status=%d msg=%s\n",
				apiErr.Type, apiErr.Code, apiErr.StatusCode, apiErr.Message)
		}
	}
}
```

## Context usage

Every method accepts a `context.Context`, giving you full control over
timeouts and cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

checkout, err := client.Checkouts.Get(ctx, "co_abc123")
```

## Pagination

### Manual pagination

```go
var cursor *string
for {
	list, err := client.Checkouts.List(ctx, &billingio.ListCheckoutsParams{
		Cursor: cursor,
		Limit:  intPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, co := range list.Data {
		fmt.Println(co.CheckoutID, co.Status)
	}

	if !list.HasMore {
		break
	}
	cursor = list.NextCursor
}

func intPtr(v int) *int { return &v }
```

### Auto-pagination

The SDK provides a generic iterator that fetches pages lazily:

```go
iter := client.Checkouts.ListAutoPaginate(ctx, nil)
for iter.Next() {
	co := iter.Current()
	fmt.Println(co.CheckoutID, co.Status)
}
if err := iter.Err(); err != nil {
	log.Fatal(err)
}
```

Auto-pagination is available on all list endpoints:

```go
client.Checkouts.ListAutoPaginate(ctx, params)
client.Webhooks.ListAutoPaginate(ctx, params)
client.Events.ListAutoPaginate(ctx, params)
client.Customers.ListAutoPaginate(ctx, params)
client.PaymentMethods.ListAutoPaginate(ctx, params)
client.PaymentLinks.ListAutoPaginate(ctx, params)
client.SubscriptionPlans.ListAutoPaginate(ctx, params)
client.Subscriptions.ListAutoPaginate(ctx, params)
client.SubscriptionRenewals.ListAutoPaginate(ctx, params)
client.Entitlements.ListAutoPaginate(ctx, params)
client.Payouts.ListAutoPaginate(ctx, params)
client.Settlements.ListAutoPaginate(ctx, params)
client.RevenueEvents.ListAutoPaginate(ctx, params)
client.Adjustments.ListAutoPaginate(ctx, params)
```

## Idempotency

Pass an idempotency key when creating checkouts to safely retry requests:

```go
checkout, err := client.Checkouts.Create(ctx, &billingio.CreateCheckoutParams{
	AmountUSD:      49.99,
	Chain:          billingio.ChainArbitrum,
	Token:          billingio.TokenUSDC,
	IdempotencyKey: "order-12345-attempt-1",
})
```

## Webhook endpoints

```go
// Create a webhook endpoint
endpoint, err := client.Webhooks.Create(ctx, &billingio.CreateWebhookParams{
	URL: "https://example.com/webhooks/billing",
	Events: []billingio.EventType{
		billingio.EventTypeCheckoutCompleted,
		billingio.EventTypeCheckoutExpired,
	},
})
// Save endpoint.Secret securely -- it is only returned on creation.

// List endpoints
list, err := client.Webhooks.List(ctx, nil)

// Delete an endpoint
err = client.Webhooks.Delete(ctx, "we_abc123")
```

## Events

```go
// List all events for a checkout
events, err := client.Events.List(ctx, &billingio.ListEventsParams{
	CheckoutID: strPtr("co_abc123"),
})

// Get a single event
event, err := client.Events.Get(ctx, "evt_abc123")

func strPtr(s string) *string { return &s }
```

## Customers

```go
// Create a customer
customer, err := client.Customers.Create(ctx, &billingio.CreateCustomerParams{
	Email: "alice@example.com",
	Name:  strPtr("Alice"),
})

// List customers
list, err := client.Customers.List(ctx, nil)

// Get a single customer
customer, err := client.Customers.Get(ctx, "cus_abc123")

// Update a customer
customer, err := client.Customers.Update(ctx, "cus_abc123", &billingio.UpdateCustomerParams{
	Name: strPtr("Alice Smith"),
})
```

## Payment methods

```go
// Create a payment method
pm, err := client.PaymentMethods.Create(ctx, &billingio.CreatePaymentMethodParams{
	CustomerID:    "cus_abc123",
	Type:          billingio.PaymentMethodTypeWallet,
	Chain:         billingio.ChainTron,
	WalletAddress: "T...",
})

// List payment methods for a customer
list, err := client.PaymentMethods.List(ctx, &billingio.ListPaymentMethodsParams{
	CustomerID: strPtr("cus_abc123"),
})

// Set as default
pm, err = client.PaymentMethods.SetDefault(ctx, "pm_abc123")

// Delete a payment method
err = client.PaymentMethods.Delete(ctx, "pm_abc123")
```

## Payment links

```go
// Create a payment link
link, err := client.PaymentLinks.Create(ctx, &billingio.CreatePaymentLinkParams{
	AmountUSD:   floatPtr(25.00),
	Description: strPtr("Pro plan"),
})

// List payment links
list, err := client.PaymentLinks.List(ctx, nil)
```

## Subscription plans

```go
// Create a plan
plan, err := client.SubscriptionPlans.Create(ctx, &billingio.CreateSubscriptionPlanParams{
	Name:            "Pro Monthly",
	AmountUSD:       29.99,
	BillingInterval: billingio.BillingIntervalMonthly,
})

// List plans
list, err := client.SubscriptionPlans.List(ctx, nil)

// Update a plan
plan, err = client.SubscriptionPlans.Update(ctx, "plan_abc123", &billingio.UpdateSubscriptionPlanParams{
	Name: strPtr("Pro Monthly (Updated)"),
})
```

## Subscriptions

```go
// Create a subscription
sub, err := client.Subscriptions.Create(ctx, &billingio.CreateSubscriptionParams{
	CustomerID: "cus_abc123",
	PlanID:     "plan_abc123",
})

// List subscriptions
list, err := client.Subscriptions.List(ctx, nil)

// Cancel a subscription
cancelled := billingio.SubscriptionStatusCancelled
sub, err = client.Subscriptions.Update(ctx, "sub_abc123", &billingio.UpdateSubscriptionParams{
	Status: &cancelled,
})
```

## Subscription renewals

```go
// List renewals
list, err := client.SubscriptionRenewals.List(ctx, &billingio.ListSubscriptionRenewalsParams{
	SubscriptionID: strPtr("sub_abc123"),
})

// Retry a failed renewal
renewal, err := client.SubscriptionRenewals.Retry(ctx, "ren_abc123")
```

## Entitlements

```go
// Create an entitlement
ent, err := client.Entitlements.Create(ctx, &billingio.CreateEntitlementParams{
	SubscriptionID: "sub_abc123",
	FeatureKey:     "api_calls",
	Value:          "10000",
})

// Check if a customer is entitled
check, err := client.Entitlements.Check(ctx, &billingio.CheckEntitlementParams{
	CustomerID: "cus_abc123",
	FeatureKey: "api_calls",
})
fmt.Printf("Entitled: %v, value: %s\n", check.Entitled, check.Value)

// Delete an entitlement
err = client.Entitlements.Delete(ctx, "ent_abc123")
```

## Payouts

```go
// Create a payout intent
payout, err := client.Payouts.Create(ctx, &billingio.CreatePayoutParams{
	AmountUSD:     500.00,
	Chain:         billingio.ChainArbitrum,
	Token:         billingio.TokenUSDC,
	WalletAddress: "0x...",
})

// Execute a pending payout
payout, err = client.Payouts.Execute(ctx, "po_abc123")

// List payouts
list, err := client.Payouts.List(ctx, nil)
```

## Settlements

```go
// List settlements
list, err := client.Settlements.List(ctx, &billingio.ListSettlementsParams{
	PayoutID: strPtr("po_abc123"),
})
```

## Revenue events

```go
// List revenue events
list, err := client.RevenueEvents.List(ctx, nil)

// Get accounting summary
summary, err := client.RevenueEvents.Accounting(ctx, &billingio.AccountingSummaryParams{
	PeriodStart: strPtr("2025-01-01T00:00:00Z"),
	PeriodEnd:   strPtr("2025-02-01T00:00:00Z"),
})
fmt.Printf("Net revenue: $%.2f\n", summary.NetRevenueUSD)
```

## Adjustments

```go
// Create a credit adjustment
adj, err := client.Adjustments.Create(ctx, &billingio.CreateAdjustmentParams{
	Type:        billingio.AdjustmentTypeCredit,
	AmountUSD:   10.00,
	CustomerID:  strPtr("cus_abc123"),
	Description: strPtr("Goodwill credit"),
})

// List adjustments
list, err := client.Adjustments.List(ctx, nil)
```

## License

See [LICENSE](LICENSE) for details.
