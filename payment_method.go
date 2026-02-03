package billingio

import (
	"context"
	"fmt"
)

// PaymentMethodService handles payment-method-related API calls.
type PaymentMethodService struct {
	client *Client
}

// Create creates a new payment method.
func (s *PaymentMethodService) Create(ctx context.Context, params *CreatePaymentMethodParams) (*PaymentMethod, error) {
	var pm PaymentMethod
	err := s.client.post(ctx, "/payment-methods", params, &pm, nil)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

// List returns a paginated list of payment methods.
func (s *PaymentMethodService) List(ctx context.Context, params *ListPaymentMethodsParams) (*PaymentMethodList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		qp["customer_id"] = strOrEmpty(params.CustomerID)
	}
	path := addQueryParams("/payment-methods", qp)

	var list PaymentMethodList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Update updates an existing payment method.
func (s *PaymentMethodService) Update(ctx context.Context, paymentMethodID string, params *UpdatePaymentMethodParams) (*PaymentMethod, error) {
	var pm PaymentMethod
	err := s.client.patch(ctx, fmt.Sprintf("/payment-methods/%s", paymentMethodID), params, &pm)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

// Delete removes a payment method.
func (s *PaymentMethodService) Delete(ctx context.Context, paymentMethodID string) error {
	return s.client.del(ctx, fmt.Sprintf("/payment-methods/%s", paymentMethodID))
}

// SetDefault marks a payment method as the default for its customer.
func (s *PaymentMethodService) SetDefault(ctx context.Context, paymentMethodID string) (*PaymentMethod, error) {
	var pm PaymentMethod
	err := s.client.post(ctx, fmt.Sprintf("/payment-methods/%s/default", paymentMethodID), nil, &pm, nil)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of payment methods. See Iter for usage details.
func (s *PaymentMethodService) ListAutoPaginate(ctx context.Context, params *ListPaymentMethodsParams) *Iter[PaymentMethod] {
	if params == nil {
		params = &ListPaymentMethodsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]PaymentMethod, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
