package billingio

import "context"

// RevenueEventService handles revenue-event and accounting API calls.
type RevenueEventService struct {
	client *Client
}

// List returns a paginated list of revenue events.
func (s *RevenueEventService) List(ctx context.Context, params *ListRevenueEventsParams) (*RevenueEventList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		if params.Type != nil {
			qp["type"] = string(*params.Type)
		}
	}
	path := addQueryParams("/revenue/events", qp)

	var list RevenueEventList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Accounting returns an aggregated revenue accounting summary.
func (s *RevenueEventService) Accounting(ctx context.Context, params *AccountingSummaryParams) (*AccountingSummary, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["period_start"] = strOrEmpty(params.PeriodStart)
		qp["period_end"] = strOrEmpty(params.PeriodEnd)
	}
	path := addQueryParams("/revenue/accounting", qp)

	var summary AccountingSummary
	err := s.client.get(ctx, path, &summary)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of revenue events. See Iter for usage details.
func (s *RevenueEventService) ListAutoPaginate(ctx context.Context, params *ListRevenueEventsParams) *Iter[RevenueEvent] {
	if params == nil {
		params = &ListRevenueEventsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]RevenueEvent, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
