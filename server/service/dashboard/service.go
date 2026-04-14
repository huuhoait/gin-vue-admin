package dashboard

import "context"

// Service orchestrates cache-aside reads for the dashboard overview.
type Service struct {
	repo  *Repository
	cache *Cache
}

// NewService creates a dashboard Service.
func NewService(repo *Repository, cache *Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

// GetOverview returns dashboard metrics using cache-aside:
// 1. Check Redis cache  2. On miss → query DB  3. Write cache  4. Return.
func (s *Service) GetOverview(ctx context.Context) (*OverviewMetrics, error) {
	if m := s.cache.Get(ctx); m != nil {
		return m, nil
	}
	m, err := s.repo.FetchOverview()
	if err != nil {
		return nil, err
	}
	s.cache.Set(ctx, m)
	return m, nil
}
