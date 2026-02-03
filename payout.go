package billingio

import (
	"context"
	"fmt"
)

// PayoutService handles payout-related API calls.
type PayoutService struct {
	client *Client
}

// Create creates a new payout intent.
func (s *PayoutService) Create(ctx context.Context, params *CreatePayoutParams) (*Payout, error) {
	var payout Payout
	err := s.client.post(ctx, "/payouts", params, &payout, nil)
	if err != nil {
		return nil, err
	}
	return &payout, nil
}

// List returns a paginated list of payouts.
func (s *PayoutService) List(ctx context.Context, params *ListPayoutsParams) (*PayoutList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		if params.Status != nil {
			qp["status"] = string(*params.Status)
		}
	}
	path := addQueryParams("/payouts", qp)

	var list PayoutList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Update updates an existing payout.
func (s *PayoutService) Update(ctx context.Context, payoutID string, params *UpdatePayoutParams) (*Payout, error) {
	var payout Payout
	err := s.client.patch(ctx, fmt.Sprintf("/payouts/%s", payoutID), params, &payout)
	if err != nil {
		return nil, err
	}
	return &payout, nil
}

// Execute triggers execution of a pending payout.
func (s *PayoutService) Execute(ctx context.Context, payoutID string) (*Payout, error) {
	var payout Payout
	err := s.client.post(ctx, fmt.Sprintf("/payouts/%s/execute", payoutID), nil, &payout, nil)
	if err != nil {
		return nil, err
	}
	return &payout, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of payouts. See Iter for usage details.
func (s *PayoutService) ListAutoPaginate(ctx context.Context, params *ListPayoutsParams) *Iter[Payout] {
	if params == nil {
		params = &ListPayoutsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]Payout, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
