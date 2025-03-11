// service.go
package zone

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateZoneRequest) (*ZoneResponse, error) {
	z := &Zone{
			Name:      req.Name,
			Active:    req.Active,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, z); err != nil {
		return nil, err
	}

	return &ZoneResponse{
			ID:        z.ID,
			Name:      z.Name,
			Active:    z.Active,
			CreatedAt: z.CreatedAt,
			UpdatedAt: z.UpdatedAt,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ZoneResponse, error) {
	z, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}

	return &ZoneResponse{
			ID:        z.ID,
			Name:      z.Name,
			Active:    z.Active,
			CreatedAt: z.CreatedAt,
			UpdatedAt: z.UpdatedAt,
	}, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateZoneRequest) (*ZoneResponse, error) {
	z, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			z.Name = *req.Name
	}

	if req.Active != nil {
			z.Active = *req.Active
	}
	z.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, z); err != nil {
		return nil, err
	}

	return &ZoneResponse{
			ID:        z.ID,
			Name:      z.Name,
			Active:    z.Active,
			CreatedAt: z.CreatedAt,
			UpdatedAt: z.UpdatedAt,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]ZoneResponse, int64, error) {
	zones, total, err := s.repo.List(ctx, p)
	if err != nil {
			return nil, 0, err
	}

	var response []ZoneResponse
	for _, z := range zones {
			response = append(response, ZoneResponse{
					ID:        z.ID,
					Name:      z.Name,
					Active:    z.Active,
					CreatedAt: z.CreatedAt,
					UpdatedAt: z.UpdatedAt,
			})
	}

	return response, total, nil
}