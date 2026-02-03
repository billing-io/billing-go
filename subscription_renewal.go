package billingio

import (
	"context"
	"fmt"
)

// SubscriptionRenewalService handles subscription-renewal-related API calls.
type SubscriptionRenewalService struct {
	client *Client
}

// List returns a paginated list of subscription renewals.
func (s *SubscriptionRenewalService) List(ctx context.Context, params *ListSubscriptionRenewalsParams) (*SubscriptionRenewalList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		qp["subscription_id"] = strOrEmpty(params.SubscriptionID)
		if params.Status != nil {
			qp["status"] = string(*params.Status)
		}
	}
	path := addQueryParams("/subscriptions/renewals", qp)

	var list SubscriptionRenewalList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Retry retries a failed subscription renewal.
func (s *SubscriptionRenewalService) Retry(ctx context.Context, renewalID string) (*SubscriptionRenewal, error) {
	var renewal SubscriptionRenewal
	err := s.client.post(ctx, fmt.Sprintf("/subscriptions/renewals/%s/retry", renewalID), nil, &renewal, nil)
	if err != nil {
		return nil, err
	}
	return &renewal, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of subscription renewals. See Iter for usage details.
func (s *SubscriptionRenewalService) ListAutoPaginate(ctx context.Context, params *ListSubscriptionRenewalsParams) *Iter[SubscriptionRenewal] {
	if params == nil {
		params = &ListSubscriptionRenewalsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]SubscriptionRenewal, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
