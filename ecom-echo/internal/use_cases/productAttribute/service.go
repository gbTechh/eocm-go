package productattribute

import (
	"context"
	"ecom/internal/shared/errors"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateProductAttributeRequest) (*ProductAttributeResponse, error) {


	productAttribute := &ProductAttribute{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, productAttribute); err != nil {
		return nil, err
	}

	return s.buildResponse(productAttribute), nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ProductAttributeResponseList, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}
	return s.buildResponseList(result), nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateProductAttributeRequest) (*ProductAttributeResponse, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			result.Name = *req.Name
	}
	if req.Type != nil {
			result.Type = *req.Type
	}
	if req.Description != nil {
			result.Description = *req.Description
	}
	
	
	result.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, result); err != nil {
			return nil, err
	}

	return s.buildResponse(result), nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]ProductAttributeResponseList, int64, error) {
	productAttributes, total, err := s.repo.List(ctx, p)
	if err != nil {
		fmt.Printf("Error en repository.List: %v\n", err)

		return nil, 0, errors.NewInternalError("Error al listar los atributos de productos", err)
	}

	response := make([]ProductAttributeResponseList, len(productAttributes))
	for i, productAttribute := range productAttributes {
		response[i] = *s.buildResponseList(&productAttribute)
	}

	return response, total, nil
}

// Helper para construir la respuesta
func (s *Service) buildResponse(c *ProductAttribute) *ProductAttributeResponse {
	return &ProductAttributeResponse{
			ID:          c.ID,
			Name:        c.Name,
			Type:        c.Type,
			Description: c.Description,		
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
	}
}
func (s *Service) buildResponseList(c *ProductAttribute) *ProductAttributeResponseList {
	return &ProductAttributeResponseList{
			ID:          c.ID,
			Name:        c.Name,
			Type:        c.Type,
			Description: c.Description,		
			Values: 		 c.Values,		
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
	}
}