package billingio

import "context"

// AdjustmentService handles revenue-adjustment-related API calls.
type AdjustmentService struct {
	client *Client
}

// List returns a paginated list of adjustments.
func (s *AdjustmentService) List(ctx context.Context, params *ListAdjustmentsParams) (*AdjustmentList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		if params.Type != nil {
			qp["type"] = string(*params.Type)
		}
	}
	path := addQueryParams("/revenue/adjustments", qp)

	var list AdjustmentList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Create creates a new revenue adjustment.
func (s *AdjustmentService) Create(ctx context.Context, params *CreateAdjustmentParams) (*Adjustment, error) {
	var adj Adjustment
	err := s.client.post(ctx, "/revenue/adjustments", params, &adj, nil)
	if err != nil {
		return nil, err
	}
	return &adj, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of adjustments. See Iter for usage details.
func (s *AdjustmentService) ListAutoPaginate(ctx context.Context, params *ListAdjustmentsParams) *Iter[Adjustment] {
	if params == nil {
		params = &ListAdjustmentsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]Adjustment, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
