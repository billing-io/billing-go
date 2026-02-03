package billingio

import "context"

// PaymentLinkService handles payment-link-related API calls.
type PaymentLinkService struct {
	client *Client
}

// Create creates a new payment link.
func (s *PaymentLinkService) Create(ctx context.Context, params *CreatePaymentLinkParams) (*PaymentLink, error) {
	var link PaymentLink
	err := s.client.post(ctx, "/payment-links", params, &link, nil)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// List returns a paginated list of payment links.
func (s *PaymentLinkService) List(ctx context.Context, params *ListPaymentLinksParams) (*PaymentLinkList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
	}
	path := addQueryParams("/payment-links", qp)

	var list PaymentLinkList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of payment links. See Iter for usage details.
func (s *PaymentLinkService) ListAutoPaginate(ctx context.Context, params *ListPaymentLinksParams) *Iter[PaymentLink] {
	if params == nil {
		params = &ListPaymentLinksParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]PaymentLink, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
