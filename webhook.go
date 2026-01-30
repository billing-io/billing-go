package billingio

import (
	"context"
	"fmt"
)

// WebhookService handles webhook endpoint API calls.
type WebhookService struct {
	client *Client
}

// Create registers a new webhook endpoint.
// The returned WebhookEndpoint includes the signing secret -- store it securely.
func (s *WebhookService) Create(ctx context.Context, params *CreateWebhookParams) (*WebhookEndpoint, error) {
	var endpoint WebhookEndpoint
	err := s.client.post(ctx, "/webhooks", params, &endpoint, nil)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

// List returns a paginated list of webhook endpoints.
func (s *WebhookService) List(ctx context.Context, params *ListParams) (*WebhookEndpointList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
	}
	path := addQueryParams("/webhooks", qp)

	var list WebhookEndpointList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Get retrieves a single webhook endpoint by ID.
func (s *WebhookService) Get(ctx context.Context, webhookID string) (*WebhookEndpoint, error) {
	var endpoint WebhookEndpoint
	err := s.client.get(ctx, fmt.Sprintf("/webhooks/%s", webhookID), &endpoint)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

// Delete removes a webhook endpoint.
func (s *WebhookService) Delete(ctx context.Context, webhookID string) error {
	return s.client.del(ctx, fmt.Sprintf("/webhooks/%s", webhookID))
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of webhook endpoints. See Iter for usage details.
func (s *WebhookService) ListAutoPaginate(ctx context.Context, params *ListParams) *Iter[WebhookEndpoint] {
	if params == nil {
		params = &ListParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]WebhookEndpoint, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
