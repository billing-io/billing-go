package billingio

import (
	"context"
	"fmt"
)

// SubscriptionService handles subscription-related API calls.
type SubscriptionService struct {
	client *Client
}

// Create creates a new subscription.
func (s *SubscriptionService) Create(ctx context.Context, params *CreateSubscriptionParams) (*Subscription, error) {
	var sub Subscription
	err := s.client.post(ctx, "/subscriptions", params, &sub, nil)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// List returns a paginated list of subscriptions.
func (s *SubscriptionService) List(ctx context.Context, params *ListSubscriptionsParams) (*SubscriptionList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		qp["customer_id"] = strOrEmpty(params.CustomerID)
		qp["plan_id"] = strOrEmpty(params.PlanID)
		if params.Status != nil {
			qp["status"] = string(*params.Status)
		}
	}
	path := addQueryParams("/subscriptions", qp)

	var list SubscriptionList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Update updates an existing subscription.
func (s *SubscriptionService) Update(ctx context.Context, subscriptionID string, params *UpdateSubscriptionParams) (*Subscription, error) {
	var sub Subscription
	err := s.client.patch(ctx, fmt.Sprintf("/subscriptions/%s", subscriptionID), params, &sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of subscriptions. See Iter for usage details.
func (s *SubscriptionService) ListAutoPaginate(ctx context.Context, params *ListSubscriptionsParams) *Iter[Subscription] {
	if params == nil {
		params = &ListSubscriptionsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]Subscription, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
