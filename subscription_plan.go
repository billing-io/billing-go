package billingio

import (
	"context"
	"fmt"
)

// SubscriptionPlanService handles subscription-plan-related API calls.
type SubscriptionPlanService struct {
	client *Client
}

// Create creates a new subscription plan.
func (s *SubscriptionPlanService) Create(ctx context.Context, params *CreateSubscriptionPlanParams) (*SubscriptionPlan, error) {
	var plan SubscriptionPlan
	err := s.client.post(ctx, "/subscriptions/plans", params, &plan, nil)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// List returns a paginated list of subscription plans.
func (s *SubscriptionPlanService) List(ctx context.Context, params *ListSubscriptionPlansParams) (*SubscriptionPlanList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		if params.Status != nil {
			qp["status"] = string(*params.Status)
		}
	}
	path := addQueryParams("/subscriptions/plans", qp)

	var list SubscriptionPlanList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Update updates an existing subscription plan.
func (s *SubscriptionPlanService) Update(ctx context.Context, planID string, params *UpdateSubscriptionPlanParams) (*SubscriptionPlan, error) {
	var plan SubscriptionPlan
	err := s.client.patch(ctx, fmt.Sprintf("/subscriptions/plans/%s", planID), params, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of subscription plans. See Iter for usage details.
func (s *SubscriptionPlanService) ListAutoPaginate(ctx context.Context, params *ListSubscriptionPlansParams) *Iter[SubscriptionPlan] {
	if params == nil {
		params = &ListSubscriptionPlansParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]SubscriptionPlan, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
