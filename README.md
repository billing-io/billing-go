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

## License

See [LICENSE](LICENSE) for details.
