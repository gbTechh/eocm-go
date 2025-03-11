// service.go
package taxcategory

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

func (s *Service) Create(ctx context.Context, req *CreateTaxCategoryRequest) (*TaxCategoryResponse, error) {
    tc := &TaxCategory{
        Name:        req.Name,
        Description: req.Description,
        Status:      req.Status,
        IsDefault:   req.IsDefault,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if req.IsDefault == true {
        if err := s.repo.UpdateAllDefaultFalse(ctx); err != nil {
            return nil, err
        }
    } 

    if err := s.repo.Create(ctx, tc); err != nil {
        return nil, err
    }

    return &TaxCategoryResponse{
        ID:          tc.ID,
        Name:        tc.Name,
        Description: tc.Description,
        Status:      tc.Status,
        IsDefault:   tc.IsDefault,
        CreatedAt:   tc.CreatedAt,
        UpdatedAt:   tc.UpdatedAt,
    }, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*TaxCategoryResponse, error) {
    tc, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return &TaxCategoryResponse{
        ID:          tc.ID,
        Name:        tc.Name,
        Description: tc.Description,
        Status:      tc.Status,
        IsDefault:   tc.IsDefault,
        CreatedAt:   tc.CreatedAt,
        UpdatedAt:   tc.UpdatedAt,
    }, nil
}

func (s *Service) GetTaxDefault(ctx context.Context) (*TaxCategoryResponse, error) {
    tc, err := s.repo.GetTaxDefault(ctx)
    
    if err != nil {
        return nil, err
    }

    return &TaxCategoryResponse{
        ID:          tc.ID,
        Name:        tc.Name,
        Description: tc.Description,
        Status:      tc.Status,
        IsDefault:   tc.IsDefault,
        CreatedAt:   tc.CreatedAt,
        UpdatedAt:   tc.UpdatedAt,
    }, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateTaxCategoryRequest) (*TaxCategoryResponse, error) {
    tc, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    if req.Name != nil {
        tc.Name = *req.Name
    }
    if req.Description != nil {
        tc.Description = *req.Description
    }
    if req.Status != nil {
        tc.Status = *req.Status
    }
    tc.UpdatedAt = time.Now()

    if req.IsDefault != nil {
        tc.IsDefault = *req.IsDefault
        if *req.IsDefault {
            if err := s.repo.UpdateAllDefaultFalse(ctx); err != nil {
                return nil, err
            }
        }
    } 

    if err := s.repo.Update(ctx, tc); err != nil {
        return nil, err
    }

    return &TaxCategoryResponse{
        ID:          tc.ID,
        Name:        tc.Name,
        Description: tc.Description,
        Status:      tc.Status,
        IsDefault:   tc.IsDefault,
        CreatedAt:   tc.CreatedAt,
        UpdatedAt:   tc.UpdatedAt,
    }, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
    return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]TaxCategoryResponse, int64, error) {
    categories, total, err := s.repo.List(ctx, p)
    if err != nil {
        return nil, 0, err
    }

    var response []TaxCategoryResponse
    for _, tc := range categories {
        response = append(response, TaxCategoryResponse{
            ID:          tc.ID,
            Name:        tc.Name,
            Description: tc.Description,
            Status:      tc.Status,
            IsDefault:   tc.IsDefault,
            CreatedAt:   tc.CreatedAt,
            UpdatedAt:   tc.UpdatedAt,
        })
    }

    return response, total, nil
}