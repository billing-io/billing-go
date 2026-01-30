package billingio

import "context"

// HealthService handles the health check endpoint.
type HealthService struct {
	client *Client
}

// Get performs a health check. This endpoint does not require authentication.
func (s *HealthService) Get(ctx context.Context) (*HealthResponse, error) {
	var resp HealthResponse
	err := s.client.get(ctx, "/health", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
