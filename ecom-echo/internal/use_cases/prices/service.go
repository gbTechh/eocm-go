// service.go
package prices

import (
	"context"
	"ecom/internal/shared/errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// PriceList methods
func (s *Service) CreatePriceList(ctx context.Context, req *CreatePriceListRequest) (*PriceListResponse, error) {
	
	
	defaultList, err := s.repo.GetDefaultPriceList(ctx)
	if err != nil {
		if _, ok := err.(*errors.AppError); ok && err.(*errors.AppError).Code == 404 {
			if !req.IsDefault {
					return nil, errors.NewBadRequestError("La primera lista de precios debe ser por defecto")
			}
		} else {
			return nil, err
		}
	} else if defaultList != nil && req.IsDefault {
		// Si ya existe una lista por defecto y esta nueva será default
		return nil, errors.NewBadRequestError("Ya existe una lista de precios por defecto")
	}
	
	

	priceList := &PriceList{
			Name:        req.Name,
			Description: req.Description,
			IsDefault:   req.IsDefault,
			Priority:    req.Priority,
			Status:      req.Status,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreatePriceList(ctx, priceList); err != nil {
			return nil, err
	}

	return mapPriceListToResponse(priceList), nil
}

func (s *Service) GetPriceList(ctx context.Context, id int64) (*PriceListResponse, error) {
	priceList, err := s.repo.GetPriceListByID(ctx, id)
	if err != nil {
			return nil, err
	}

	return mapPriceListToResponse(priceList), nil
}

func (s *Service) UpdatePriceList(ctx context.Context, id int64, req *UpdatePriceListRequest) (*PriceListResponse, error) {
	priceList, err := s.repo.GetPriceListByID(ctx, id)
	if err != nil {
		return nil, err
	}

	wasDefault := priceList.IsDefault

	if req.Name != nil {
		priceList.Name = *req.Name
	}
	if req.Description != nil {
		priceList.Description = *req.Description
	}
	if req.Priority != nil {
		priceList.Priority = *req.Priority
	}
	if req.Status != nil {
		priceList.Status = *req.Status
	}
	if req.IsDefault != nil {
		if wasDefault && !*req.IsDefault {
				// Verificar si hay otras listas que podrían ser default
			count, err := s.repo.CountPriceLists(ctx)
			if err != nil {
				return nil, err
			}
			if count > 1 {
				return nil, errors.NewBadRequestError("Debe existir al menos una lista de precios por defecto. Asigne otra lista como default antes de quitar ésta")
			}
		}

		if !wasDefault && *req.IsDefault {
			defaultList, err := s.repo.GetDefaultPriceList(ctx)
			if err != nil && !errors.IsNotFound(err) {
				return nil, err
			}
			if defaultList != nil {
				// Actualizar la lista anterior para quitar el default
				defaultList.IsDefault = false
				if err := s.repo.UpdatePriceList(ctx, defaultList); err != nil {
					return nil, err
				}
			}
		}
		priceList.IsDefault = *req.IsDefault
	}

	// Si está cambiando a default, verificar que no exista otro
	if req.IsDefault != nil && *req.IsDefault && !priceList.IsDefault {
			defaultList, err := s.repo.GetDefaultPriceList(ctx)
			if err != nil && !errors.IsNotFound(err) {
					return nil, err
			}
			if defaultList != nil && defaultList.ID != id {
					return nil, errors.NewBadRequestError("Ya existe una lista de precios por defecto")
			}
	}

	priceList.UpdatedAt = time.Now()
	if err := s.repo.UpdatePriceList(ctx, priceList); err != nil {
			return nil, err
	}

	return mapPriceListToResponse(priceList), nil
}

func (s *Service) DeletePriceList(ctx context.Context, id int64) error {
	// No permitir eliminar la lista por defecto
	priceList, err := s.repo.GetPriceListByID(ctx, id)
	if err != nil {
			return err
	}

	if priceList.IsDefault {
			return errors.NewBadRequestError("No se puede eliminar la lista de precios por defecto")
	}

	return s.repo.DeletePriceList(ctx, id)
}

func (s *Service) ListPriceLists(ctx context.Context, p *Pagination) ([]PriceListResponse, int64, error) {
	priceLists, total, err := s.repo.ListPriceLists(ctx, p)
	if err != nil {
			return nil, 0, err
	}

	var response []PriceListResponse
	for _, pl := range priceLists {
			response = append(response, *mapPriceListToResponse(&pl))
	}

	return response, total, nil
}

// Price methods
func (s *Service) CreatePrice(ctx context.Context, req *CreatePriceRequest) (*PriceResponse, error) {
	// Validar que la lista de precios exista
	_, err := s.repo.GetPriceListByID(ctx, req.IDPriceList)
	if err != nil {
			return nil, err
	}

	// Validar fechas
	if err := validatePriceDates(req.StartsAt, req.EndsAt); err != nil {
			return nil, err
	}

	price := &Price{
			IDPriceList: req.IDPriceList,
			Amount:      req.Amount,
			StartsAt:    req.StartsAt,
			EndsAt:      req.EndsAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreatePrice(ctx, price); err != nil {
			return nil, err
	}

	return s.GetPrice(ctx, price.ID)
}

func (s *Service) GetPrice(ctx context.Context, id int64) (*PriceResponse, error) {
	price, err := s.repo.GetPriceByID(ctx, id)
	if err != nil {
			return nil, err
	}

	return mapPriceToResponse(price), nil
}

func (s *Service) UpdatePrice(ctx context.Context, id int64, req *UpdatePriceRequest) (*PriceResponse, error) {
	price, err := s.repo.GetPriceByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Amount != nil {
			price.Amount = *req.Amount
	}
	if req.StartsAt != nil {
			price.StartsAt = *req.StartsAt
	}
	if req.EndsAt != nil {
			price.EndsAt = req.EndsAt
	}

	// Validar fechas
	if err := validatePriceDates(price.StartsAt, price.EndsAt); err != nil {
			return nil, err
	}

	price.UpdatedAt = time.Now()
	if err := s.repo.UpdatePrice(ctx, price); err != nil {
			return nil, err
	}

	return s.GetPrice(ctx, id)
}

func (s *Service) DeletePrice(ctx context.Context, id int64) error {
	return s.repo.DeletePrice(ctx, id)
}

func (s *Service) ListPrices(ctx context.Context, priceListID int64, p *Pagination) ([]PriceResponse, int64, error) {
	// Validar que la lista de precios exista
	if _, err := s.repo.GetPriceListByID(ctx, priceListID); err != nil {
			return nil, 0, err
	}

	prices, total, err := s.repo.ListPrices(ctx, priceListID, p)
	if err != nil {
			return nil, 0, err
	}

	var response []PriceResponse
	for _, price := range prices {
			response = append(response, *mapPriceToResponse(&price))
	}

	return response, total, nil
}

// ProductVariantPrice methods
func (s *Service) AssignPrice(ctx context.Context, req *AssignPriceRequest) error {
	// Validar que el precio exista y esté vigente
	price, err := s.repo.GetPriceByID(ctx, req.IDPrice)
	if err != nil {
			return err
	}

	now := time.Now()
	if price.StartsAt.After(now) {
			return errors.NewBadRequestError("El precio aún no está vigente")
	}
	if price.EndsAt != nil && price.EndsAt.Before(now) {
			return errors.NewBadRequestError("El precio ha expirado")
	}

	pvp := &ProductVariantPrice{
			IDProductVariant: req.IDProductVariant,
			IDPrice:         req.IDPrice,
			IsActive:        req.IsActive,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
	}

	return s.repo.AssignPrice(ctx, pvp)
}

func (s *Service) UnassignPrice(ctx context.Context, productVariantID, priceID int64) error {
	return s.repo.UnassignPrice(ctx, productVariantID, priceID)
}

func (s *Service) GetActivePrice(ctx context.Context, productVariantID int64) (*PriceResponse, error) {
	price, err := s.repo.GetActivePrice(ctx, productVariantID)
	if err != nil {
			return nil, err
	}

	return mapPriceToResponse(price), nil
}

func (s *Service) GetAllPrices(ctx context.Context, productVariantID int64) ([]ProductVariantPrice, error) {
	return s.repo.GetAllPrices(ctx, productVariantID)
}

// Helper functions
func validatePriceDates(startsAt time.Time, endsAt *time.Time) error {
	if endsAt == nil {
			return nil
	}

	if !startsAt.Before(*endsAt) {
			return errors.NewBadRequestError("La fecha de inicio debe ser anterior a la fecha de fin")
	}

	return nil
}

func mapPriceListToResponse(pl *PriceList) *PriceListResponse {
	return &PriceListResponse{
			ID:          pl.ID,
			Name:        pl.Name,
			Description: pl.Description,
			IsDefault:   pl.IsDefault,
			Priority:    pl.Priority,
			Status:      pl.Status,
			CreatedAt:   pl.CreatedAt,
			UpdatedAt:   pl.UpdatedAt,
	}
}

func mapPriceToResponse(p *Price) *PriceResponse {
	return &PriceResponse{
			ID:          p.ID,
			IDPriceList: p.IDPriceList,
			Amount:      p.Amount,
			StartsAt:    p.StartsAt,
			EndsAt:      p.EndsAt,
			PriceList:   mapPriceListToResponse(p.PriceList),
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
	}
}