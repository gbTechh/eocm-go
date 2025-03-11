package category

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

func (s *Service) Create(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error) {
	// Validar que existe el parent_id si se proporciona
	if req.ParentID != nil && *req.ParentID != 0 {
			_, err := s.repo.GetByID(ctx, *req.ParentID)
			if err != nil {
				return nil, errors.NewBadRequestError("Categoría padre no encontrada")
			}
	}

	category := &Category{
			Name:        req.Name,
			Slug:        req.Slug,
			Description: req.Description,
			ParentID:    req.ParentID,
			IDMedia:     req.IDMedia,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return s.buildResponse(category), nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}
	return s.buildResponse(category), nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			category.Name = *req.Name
	}
	if req.Slug != nil {
			category.Slug = *req.Slug
	}
	if req.Description != nil {
			category.Description = *req.Description
	}
	
	if req.ParentID != nil {		
		if *req.ParentID == id {
			return nil, errors.NewValidationError("La categoría no puede tener un parent_id igual al id")
		}
		// Validar que existe la categoría padre
		if req.ParentID != nil && *req.ParentID != 0 {
			_, err := s.repo.GetByID(ctx, *req.ParentID)
			if err != nil {
				return nil, errors.NewBadRequestError("Categoría padre no encontrada")
			}
		}
		category.ParentID = req.ParentID
	} else{		
		category.ParentID = req.ParentID
	}

	if req.IDMedia != nil {
			category.IDMedia = req.IDMedia
	}

	category.UpdatedAt = time.Now()
	category.Media = nil

	if err := s.repo.Update(ctx, category); err != nil {
			return nil, err
	}

	return s.buildResponse(category), nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	// Verificar si hay subcategorías
	hasChildren, err := s.repo.HasChildren(ctx, id)
	if err != nil {
			return errors.NewInternalError("Error al verificar subcategorías", err)
	}
	if hasChildren {
			return errors.NewBadRequestError("No se puede eliminar una categoría con subcategorías")
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]CategoryResponse, int64, error) {
	categories, total, err := s.repo.List(ctx, p)
	if err != nil {
		fmt.Printf("Error en repository.List: %v\n", err)

		return nil, 0, errors.NewInternalError("Error al listar categorías", err)
	}

	response := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		response[i] = *s.buildResponse(&category)
	}

	return response, total, nil
}

// Helper para construir la respuesta
func (s *Service) buildResponse(c *Category) *CategoryResponse {
	return &CategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			ParentID:    c.ParentID,
			IDMedia:     c.IDMedia,
			Media:       c.Media,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
	}
}