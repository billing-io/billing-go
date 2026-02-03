package billingio

import (
	"context"
	"fmt"
)

// EntitlementService handles entitlement-related API calls.
type EntitlementService struct {
	client *Client
}

// List returns a paginated list of entitlements.
func (s *EntitlementService) List(ctx context.Context, params *ListEntitlementsParams) (*EntitlementList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		qp["subscription_id"] = strOrEmpty(params.SubscriptionID)
		qp["feature_key"] = strOrEmpty(params.FeatureKey)
	}
	path := addQueryParams("/subscriptions/entitlements", qp)

	var list EntitlementList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Create creates a new entitlement.
func (s *EntitlementService) Create(ctx context.Context, params *CreateEntitlementParams) (*Entitlement, error) {
	var ent Entitlement
	err := s.client.post(ctx, "/subscriptions/entitlements", params, &ent, nil)
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

// Update updates an existing entitlement.
func (s *EntitlementService) Update(ctx context.Context, entitlementID string, params *UpdateEntitlementParams) (*Entitlement, error) {
	var ent Entitlement
	err := s.client.patch(ctx, fmt.Sprintf("/subscriptions/entitlements/%s", entitlementID), params, &ent)
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

// Delete removes an entitlement.
func (s *EntitlementService) Delete(ctx context.Context, entitlementID string) error {
	return s.client.del(ctx, fmt.Sprintf("/subscriptions/entitlements/%s", entitlementID))
}

// Check checks whether a customer is entitled to a specific feature.
func (s *EntitlementService) Check(ctx context.Context, params *CheckEntitlementParams) (*EntitlementCheckResponse, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["customer_id"] = params.CustomerID
		qp["feature_key"] = params.FeatureKey
	}
	path := addQueryParams("/subscriptions/entitlements/check", qp)

	var resp EntitlementCheckResponse
	err := s.client.get(ctx, path, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of entitlements. See Iter for usage details.
func (s *EntitlementService) ListAutoPaginate(ctx context.Context, params *ListEntitlementsParams) *Iter[Entitlement] {
	if params == nil {
		params = &ListEntitlementsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]Entitlement, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
