package taxzone

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

func (s *Service) Create(ctx context.Context, req *CreateTaxZoneRequest) (*TaxZoneResponse, error) {
  tz := &TaxZone{
      IDZone:      req.IDZone,
      Name:        req.Name,
      Description: req.Description,
      Status:      req.Status,
      CreatedAt:   time.Now(),
      UpdatedAt:   time.Now(),
  }

  if err := s.repo.Create(ctx, tz); err != nil {
      return nil, err
  }

  return s.GetByID(ctx, tz.ID)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*TaxZoneResponse, error) {
  tz, err := s.repo.GetByID(ctx, id)
  if err != nil {
      return nil, err
  }

  return &TaxZoneResponse{
      ID:          tz.ID,
      IDZone:      tz.IDZone,
      Name:        tz.Name,
      Description: tz.Description,
      Status:      tz.Status,
      Zone:        tz.Zone,
      CreatedAt:   tz.CreatedAt,
      UpdatedAt:   tz.UpdatedAt,
  }, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateTaxZoneRequest) (*TaxZoneResponse, error) {
  tz, err := s.repo.GetByID(ctx, id)
  if err != nil {
      return nil, err
  }

  if req.IDZone != nil {
      tz.IDZone = *req.IDZone
  }
  if req.Name != nil {
      tz.Name = *req.Name
  }
  if req.Description != nil {
      tz.Description = *req.Description
  }
  if req.Status != nil {
      tz.Status = *req.Status
  }
  tz.UpdatedAt = time.Now()

  if err := s.repo.Update(ctx, tz); err != nil {
      return nil, err
  }

  return s.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
  return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]TaxZoneResponse, int64, error) {
  zones, total, err := s.repo.List(ctx, p)
  if err != nil {
      return nil, 0, err
  }

  var response []TaxZoneResponse
  for _, tz := range zones {
    response = append(response, TaxZoneResponse{
      ID:          tz.ID,
      IDZone:      tz.IDZone,
      Name:        tz.Name,
      Description: tz.Description,
      Status:      tz.Status,
      Zone:        tz.Zone,
      CreatedAt:   tz.CreatedAt,
      UpdatedAt:   tz.UpdatedAt,
    })
  }

  return response, total, nil
}