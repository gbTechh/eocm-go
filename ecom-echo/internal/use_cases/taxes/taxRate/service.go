// service.go
package taxrate

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

func (s *Service) Create(ctx context.Context, req *CreateTaxRateRequest) (*TaxRateResponse, error) {
  tr := &TaxRate{
      IDTaxCategory: req.IDTaxCategory,
      IDTaxZone:     req.IDTaxZone,
      Percentage:    req.Percentage,
      IsDefault:     req.IsDefault,
      Status:        req.Status,
      Name:          req.Name,
      CreatedAt:     time.Now(),
      UpdatedAt:     time.Now(),
  }

  if err := s.repo.Create(ctx, tr); err != nil {
      return nil, err
  }

  return s.GetByID(ctx, tr.ID)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*TaxRateResponse, error) {
  tr, err := s.repo.GetByID(ctx, id)
  if err != nil {
      return nil, err
  }

  return &TaxRateResponse{
      ID:            tr.ID,
      IDTaxCategory: tr.IDTaxCategory,
      IDTaxZone:     tr.IDTaxZone,
      Percentage:    tr.Percentage,
      IsDefault:     tr.IsDefault,
      Status:        tr.Status,
      Name:          tr.Name,
      TaxCategory:   tr.TaxCategory,
      TaxZone:      tr.TaxZone,
      CreatedAt:     tr.CreatedAt,
      UpdatedAt:     tr.UpdatedAt,
  }, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateTaxRateRequest) (*TaxRateResponse, error) {
  tr, err := s.repo.GetByID(ctx, id)
  if err != nil {
      return nil, err
  }

  if req.IDTaxCategory != nil {
      tr.IDTaxCategory = *req.IDTaxCategory
  }
  if req.IDTaxZone != nil {
      tr.IDTaxZone = *req.IDTaxZone
  }
  if req.Percentage != nil {
      tr.Percentage = *req.Percentage
  }
  if req.IsDefault != nil {
      tr.IsDefault = *req.IsDefault
  }
  if req.Status != nil {
      tr.Status = *req.Status
  }
  if req.Name != nil {
      tr.Name = *req.Name
  }
  tr.UpdatedAt = time.Now()

  if err := s.repo.Update(ctx, tr); err != nil {
      return nil, err
  }

  return s.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
  return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]TaxRateResponse, int64, error) {
  rates, total, err := s.repo.List(ctx, p)
  if err != nil {
      return nil, 0, err
  }

  var response []TaxRateResponse
  for _, tr := range rates {
      response = append(response, TaxRateResponse{
          ID:            tr.ID,
          IDTaxCategory: tr.IDTaxCategory,
          IDTaxZone:     tr.IDTaxZone,
          Percentage:    tr.Percentage,
          IsDefault:     tr.IsDefault,
          Status:        tr.Status,
          Name:          tr.Name,
          TaxCategory:   tr.TaxCategory,
          TaxZone:      tr.TaxZone,
          CreatedAt:     tr.CreatedAt,
          UpdatedAt:     tr.UpdatedAt,
      })
  }

  return response, total, nil
}
