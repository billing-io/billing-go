package billingio

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	// SignatureHeader is the HTTP header containing the webhook signature.
	SignatureHeader = "X-Billing-Signature"

	// DefaultTolerance is the default maximum age (in seconds) an event
	// timestamp may differ from the current time.
	DefaultTolerance = 300 // 5 minutes
)

// WebhookVerificationError is returned when webhook signature verification fails.
type WebhookVerificationError struct {
	Message string
}

func (e *WebhookVerificationError) Error() string {
	return fmt.Sprintf("billingio: webhook verification failed: %s", e.Message)
}

// VerifyWebhookSignature verifies the HMAC-SHA256 signature of an incoming
// webhook payload and returns the parsed event.
//
// payload is the raw request body bytes (do NOT parse as JSON first).
// header is the value of the X-Billing-Signature header.
// secret is your webhook endpoint signing secret (prefixed whsec_).
//
// The default timestamp tolerance of 300 seconds (5 minutes) is applied.
// Use VerifyWebhookSignatureWithTolerance for a custom tolerance.
func VerifyWebhookSignature(payload []byte, header string, secret string) (*WebhookEvent, error) {
	return VerifyWebhookSignatureWithTolerance(payload, header, secret, DefaultTolerance)
}

// VerifyWebhookSignatureWithTolerance is like VerifyWebhookSignature but allows
// specifying a custom timestamp tolerance in seconds.
//
// Set tolerance to 0 to disable timestamp checking entirely.
func VerifyWebhookSignatureWithTolerance(payload []byte, header string, secret string, tolerance int) (*WebhookEvent, error) {
	if header == "" {
		return nil, &WebhookVerificationError{Message: "missing signature header"}
	}
	if secret == "" {
		return nil, &WebhookVerificationError{Message: "missing webhook secret"}
	}

	timestamp, signature, err := parseSignatureHeader(header)
	if err != nil {
		return nil, err
	}

	// Check timestamp tolerance (skip if tolerance is 0).
	if tolerance > 0 {
		now := time.Now().Unix()
		diff := math.Abs(float64(now - timestamp))
		if diff > float64(tolerance) {
			return nil, &WebhookVerificationError{
				Message: fmt.Sprintf(
					"timestamp outside tolerance: event=%d, now=%d, tolerance=%ds",
					timestamp, now, tolerance,
				),
			}
		}
	}

	// Build the signed payload: "{timestamp}.{body}"
	signedPayload := fmt.Sprintf("%d.%s", timestamp, string(payload))

	// Compute HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expected := hex.EncodeToString(mac.Sum(nil))

	// Constant-time comparison
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return nil, &WebhookVerificationError{Message: "signature mismatch"}
	}

	// Parse the event payload
	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, &WebhookVerificationError{Message: "invalid JSON in webhook body"}
	}

	return &event, nil
}

// parseSignatureHeader extracts the timestamp and v1 signature from the header.
// Expected format: t={unix_timestamp},v1={hex_hmac_sha256}
func parseSignatureHeader(header string) (int64, string, error) {
	parts := make(map[string]string)
	for _, segment := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(segment), "=", 2)
		if len(kv) == 2 {
			parts[kv[0]] = kv[1]
		}
	}

	tsStr, ok := parts["t"]
	if !ok || tsStr == "" {
		return 0, "", &WebhookVerificationError{
			Message: "invalid signature header format: missing timestamp (t=)",
		}
	}

	timestamp, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return 0, "", &WebhookVerificationError{
			Message: "invalid signature header format: non-numeric timestamp",
		}
	}

	sig, ok := parts["v1"]
	if !ok || sig == "" {
		return 0, "", &WebhookVerificationError{
			Message: "invalid signature header format: missing signature (v1=)",
		}
	}

	return timestamp, sig, nil
}
