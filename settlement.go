package billingio

import "context"

// SettlementService handles settlement-related API calls.
type SettlementService struct {
	client *Client
}

// List returns a paginated list of settlements.
func (s *SettlementService) List(ctx context.Context, params *ListSettlementsParams) (*SettlementList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		qp["payout_id"] = strOrEmpty(params.PayoutID)
	}
	path := addQueryParams("/payouts/settlements", qp)

	var list SettlementList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of settlements. See Iter for usage details.
func (s *SettlementService) ListAutoPaginate(ctx context.Context, params *ListSettlementsParams) *Iter[Settlement] {
	if params == nil {
		params = &ListSettlementsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]Settlement, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
