package billingio

import (
	"context"
	"fmt"
)

// CheckoutService handles checkout-related API calls.
type CheckoutService struct {
	client *Client
}

// Create creates a new payment checkout.
//
// If params.IdempotencyKey is set it is sent as the Idempotency-Key header.
func (s *CheckoutService) Create(ctx context.Context, params *CreateCheckoutParams) (*Checkout, error) {
	var headers map[string]string
	if params.IdempotencyKey != "" {
		headers = map[string]string{
			"Idempotency-Key": params.IdempotencyKey,
		}
	}

	var checkout Checkout
	err := s.client.post(ctx, "/checkouts", params, &checkout, headers)
	if err != nil {
		return nil, err
	}
	return &checkout, nil
}

// List returns a paginated list of checkouts, newest first.
func (s *CheckoutService) List(ctx context.Context, params *ListCheckoutsParams) (*CheckoutList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		if params.Status != nil {
			qp["status"] = string(*params.Status)
		}
	}
	path := addQueryParams("/checkouts", qp)

	var list CheckoutList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Get retrieves a single checkout by ID.
func (s *CheckoutService) Get(ctx context.Context, checkoutID string) (*Checkout, error) {
	var checkout Checkout
	err := s.client.get(ctx, fmt.Sprintf("/checkouts/%s", checkoutID), &checkout)
	if err != nil {
		return nil, err
	}
	return &checkout, nil
}

// GetStatus returns the lightweight polling status of a checkout.
func (s *CheckoutService) GetStatus(ctx context.Context, checkoutID string) (*CheckoutStatusResponse, error) {
	var status CheckoutStatusResponse
	err := s.client.get(ctx, fmt.Sprintf("/checkouts/%s/status", checkoutID), &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of checkouts. See Iter for usage details.
func (s *CheckoutService) ListAutoPaginate(ctx context.Context, params *ListCheckoutsParams) *Iter[Checkout] {
	if params == nil {
		params = &ListCheckoutsParams{}
	}
	p := *params // shallow copy so we can mutate cursor

	return newIter(func(cursor *string) ([]Checkout, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
