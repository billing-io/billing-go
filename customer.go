package billingio

import (
	"context"
	"fmt"
)

// CustomerService handles customer-related API calls.
type CustomerService struct {
	client *Client
}

// Create creates a new customer.
func (s *CustomerService) Create(ctx context.Context, params *CreateCustomerParams) (*Customer, error) {
	var customer Customer
	err := s.client.post(ctx, "/customers", params, &customer, nil)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// List returns a paginated list of customers.
func (s *CustomerService) List(ctx context.Context, params *ListCustomersParams) (*CustomerList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		if params.Status != nil {
			qp["status"] = string(*params.Status)
		}
	}
	path := addQueryParams("/customers", qp)

	var list CustomerList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Get retrieves a single customer by ID.
func (s *CustomerService) Get(ctx context.Context, customerID string) (*Customer, error) {
	var customer Customer
	err := s.client.get(ctx, fmt.Sprintf("/customers/%s", customerID), &customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// Update updates an existing customer.
func (s *CustomerService) Update(ctx context.Context, customerID string, params *UpdateCustomerParams) (*Customer, error) {
	var customer Customer
	err := s.client.patch(ctx, fmt.Sprintf("/customers/%s", customerID), params, &customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of customers. See Iter for usage details.
func (s *CustomerService) ListAutoPaginate(ctx context.Context, params *ListCustomersParams) *Iter[Customer] {
	if params == nil {
		params = &ListCustomersParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]Customer, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
