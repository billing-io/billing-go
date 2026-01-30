package billingio

import (
	"context"
	"fmt"
)

// EventService handles event-related API calls.
type EventService struct {
	client *Client
}

// List returns a paginated list of events, newest first.
func (s *EventService) List(ctx context.Context, params *ListEventsParams) (*EventList, error) {
	qp := make(map[string]string)
	if params != nil {
		qp["cursor"] = strOrEmpty(params.Cursor)
		qp["limit"] = intToString(params.Limit)
		if params.Type != nil {
			qp["type"] = string(*params.Type)
		}
		qp["checkout_id"] = strOrEmpty(params.CheckoutID)
	}
	path := addQueryParams("/events", qp)

	var list EventList
	err := s.client.get(ctx, path, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Get retrieves a single event by ID.
func (s *EventService) Get(ctx context.Context, eventID string) (*Event, error) {
	var event Event
	err := s.client.get(ctx, fmt.Sprintf("/events/%s", eventID), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// ListAutoPaginate returns an iterator that automatically fetches subsequent
// pages of events. See Iter for usage details.
func (s *EventService) ListAutoPaginate(ctx context.Context, params *ListEventsParams) *Iter[Event] {
	if params == nil {
		params = &ListEventsParams{}
	}
	p := *params

	return newIter(func(cursor *string) ([]Event, bool, *string, error) {
		p.Cursor = cursor
		list, err := s.List(ctx, &p)
		if err != nil {
			return nil, false, nil, err
		}
		return list.Data, list.HasMore, list.NextCursor, nil
	})
}
