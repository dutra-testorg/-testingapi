package micro

import "context"

// Service struct to hold repository
type Service struct {
	repo *Repository
}

// NewService create service struct
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Demo interface to repository
func (s *Service) Demo(ctx context.Context, uid string) (Demo, error) {
	return s.repo.Demo(ctx, uid)
}
