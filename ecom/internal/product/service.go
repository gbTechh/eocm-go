// internal/product/service.go
package product

import (
	"context"
	"ecommerce/pkg/validator"
	"strconv"

	"github.com/google/uuid"
)

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input *CreateProductInput) (*Product, error) {
    if err := validator.ValidateStruct(input); err != nil {
        return nil, err
    }

    product := &Product{
        ID:          uuid.New().String(),
        Name:        input.Name,
        Description: input.Description,
        Price:       input.Price,
        Stock:       input.Stock,
        CategoryID:  input.CategoryID,
    }

    if err := s.repo.Create(ctx, product); err != nil {
        return nil, err
    }

    // Obtener el producto completo con sus relaciones
    return s.repo.GetByID(ctx, product.ID)
}

func (s *Service) Update(ctx context.Context, id string, input *UpdateProductInput) (*Product, error) {
    if err := validator.ValidateStruct(input); err != nil {
        return nil, err
    }

    product, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Actualizar solo los campos proporcionados
    if input.Name != nil {
        product.Name = *input.Name
    }
    if input.Description != nil {
        product.Description = *input.Description
    }
    if input.Price != nil {
        product.Price = *input.Price
    }
    if input.Stock != nil {
        product.Stock = *input.Stock
    }
    if input.CategoryID != nil {
        product.CategoryID = *input.CategoryID
    }

    if err := s.repo.Update(ctx, product); err != nil {
        return nil, err
    }

    return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, id string) (*Product, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, pageStr, limitStr string) (*ProductList, error) {
    page := 1
    limit := 10

    // Convertir y validar parámetros de paginación
    if pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }

    products, err := s.repo.List(ctx, page, limit)
    if err != nil {
        return nil, err
    }

    return &ProductList{
        Products: products,
        Page:     page,
        Limit:    limit,
    }, nil
}

func (s *Service) ListCategories(ctx context.Context) ([]*Category, error) {
    return s.repo.ListCategories(ctx)
}

type ProductList struct {
    Products []*Product `json:"products"`
    Page     int       `json:"page"`
    Limit    int       `json:"limit"`
}