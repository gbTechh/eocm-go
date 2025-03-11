// service.go
package tag

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

func (s *Service) Create(ctx context.Context, req *CreateTagRequest) (*TagResponse, error) {
	tag := &Tag{
			Name:      req.Name,
			Code:      req.Code,
			Active:    req.Active,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, err
	}

	return &TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Code:      tag.Code,
			Active:    tag.Active,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*TagResponse, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}

	return &TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Code:      tag.Code,
			Active:    tag.Active,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
	}, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateTagRequest) (*TagResponse, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			tag.Name = *req.Name
	}
	if req.Code != nil {
			tag.Code = *req.Code
	}
	if req.Active != nil {
			tag.Active = *req.Active
	}
	tag.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, tag); err != nil {
		return nil, err
	}

	return &TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Code:      tag.Code,
			Active:    tag.Active,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]TagResponse, int64, error) {
	tags, total, err := s.repo.List(ctx, p)
	if err != nil {
			return nil, 0, err
	}

	var response []TagResponse
	for _, tag := range tags {
			response = append(response, TagResponse{
					ID:        tag.ID,
					Name:      tag.Name,
					Code:      tag.Code,
					Active:    tag.Active,
					CreatedAt: tag.CreatedAt,
					UpdatedAt: tag.UpdatedAt,
			})
	}

	return response, total, nil
}