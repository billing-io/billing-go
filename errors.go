package billingio

import (
	"errors"
	"fmt"
)

// Error represents a structured API error returned by billing.io.
type Error struct {
	// Type is the high-level error category (e.g. "invalid_request", "not_found").
	Type string `json:"type"`

	// Code is a machine-readable error code (e.g. "checkout_not_found").
	Code string `json:"code"`

	// StatusCode is the HTTP status code of the response.
	StatusCode int `json:"-"`

	// Message is a human-readable explanation.
	Message string `json:"message"`

	// Param is the request parameter that caused the error, if applicable.
	Param *string `json:"param"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	msg := fmt.Sprintf("billingio: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
	if e.Param != nil {
		msg += fmt.Sprintf(", param=%s", *e.Param)
	}
	return msg
}

// errorResponse mirrors the JSON envelope returned by the API.
type errorResponse struct {
	Err *Error `json:"error"`
}

// IsNotFound reports whether err is a billing.io "not_found" error.
func IsNotFound(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.Type == "not_found"
	}
	return false
}

// IsRateLimited reports whether err is a billing.io "rate_limited" error.
func IsRateLimited(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.Type == "rate_limited"
	}
	return false
}

// IsAuthError reports whether err is a billing.io "authentication_error".
func IsAuthError(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.Type == "authentication_error"
	}
	return false
}
