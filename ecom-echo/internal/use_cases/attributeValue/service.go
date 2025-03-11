package attributevalue

import (
	"context"
	"ecom/internal/shared/errors"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateAttributeValueRequest) (*AttributeValueResponse, error) {


	productAttribute := &AttributeValue{
		Name:        				req.Name,
		IDProductAttribute: req.IDProductAttribute,
	}

	if err := s.repo.Create(ctx, productAttribute); err != nil {
		return nil, err
	}

	return s.buildResponse(productAttribute), nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*AttributeValueResponse, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}
	return s.buildResponse(result), nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateAttributeValueRequest) (*AttributeValueResponse, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			result.Name = *req.Name
	}


	if err := s.repo.Update(ctx, result); err != nil {
			return nil, err
	}

	return s.buildResponse(result), nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]AttributeValueResponse, int64, error) {
	data, total, err := s.repo.List(ctx, p)
	if err != nil {
		fmt.Printf("Error en repository.List: %v\n", err)

		return nil, 0, errors.NewInternalError("Error al listar los atributos de productos", err)
	}

	response := make([]AttributeValueResponse, len(data))
	for i, productAttribute := range data {
		response[i] = *s.buildResponse(&productAttribute)
	}

	return response, total, nil
}

// Helper para construir la respuesta
func (s *Service) buildResponse(c *AttributeValue) *AttributeValueResponse {
	return &AttributeValueResponse{
			ID:          				c.ID,
			Name:        				c.Name,
			IDProductAttribute: c.IDProductAttribute,
	}
}