// service.go
package taxrule

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

func (s *Service) Create(ctx context.Context, req *CreateTaxRuleRequest) (*TaxRuleResponse, error) {
  tr := &TaxRule{
      IDTaxRate: req.IDTaxRate,
      Priority:  req.Priority,
      Status:    req.Status,
      MinAmount: req.MinAmount,
      MaxAmount: req.MaxAmount,
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
  }

  if err := s.repo.Create(ctx, tr); err != nil {
      return nil, err
  }

  return s.GetByID(ctx, tr.ID)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*TaxRuleResponse, error) {
  tr, err := s.repo.GetByID(ctx, id)
  if err != nil {
      return nil, err
  }

  return &TaxRuleResponse{
      ID:        tr.ID,
      IDTaxRate: tr.IDTaxRate,
      Priority:  tr.Priority,
      Status:    tr.Status,
      MinAmount: tr.MinAmount,
      MaxAmount: tr.MaxAmount,
      TaxRate:   tr.TaxRate,
      CreatedAt: tr.CreatedAt,
      UpdatedAt: tr.UpdatedAt,
  }, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateTaxRuleRequest) (*TaxRuleResponse, error) {
  tr, err := s.repo.GetByID(ctx, id)
  if err != nil {
      return nil, err
  }

  if req.IDTaxRate != nil {
      tr.IDTaxRate = *req.IDTaxRate
  }
  if req.Priority != nil {
      tr.Priority = *req.Priority
  }
  if req.Status != nil {
      tr.Status = *req.Status
  }
  if req.MinAmount != nil {
      tr.MinAmount = req.MinAmount
  }
  if req.MaxAmount != nil {
      tr.MaxAmount = req.MaxAmount
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

func (s *Service) List(ctx context.Context, p *Pagination) ([]TaxRuleResponse, int64, error) {
  rules, total, err := s.repo.List(ctx, p)
  if err != nil {
      return nil, 0, err
  }

  var response []TaxRuleResponse
  for _, tr := range rules {
      response = append(response, TaxRuleResponse{
          ID:        tr.ID,
          IDTaxRate: tr.IDTaxRate,
          Priority:  tr.Priority,
          Status:    tr.Status,
          MinAmount: tr.MinAmount,
          MaxAmount: tr.MaxAmount,
          TaxRate:   tr.TaxRate,
          CreatedAt: tr.CreatedAt,
          UpdatedAt: tr.UpdatedAt,
      })
  }

  return response, total, nil
}
