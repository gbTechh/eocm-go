// service.go
package currency

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

// Currency operations
func (s *Service) CreateCurrency(ctx context.Context, req *CreateCurrencyRequest) (*CurrencyResponse, error) {
	// Si es la primera moneda, debe ser base
	currencies, _, err := s.repo.ListCurrencies(ctx, &Pagination{Page: 1, PerPage: 1})
	if err != nil {
			return nil, err
	}
	if len(currencies) == 0 {
			req.IsBase = true
	}

	// Si es base, verificar que no exista otra base
	if req.IsBase {
			baseCurrency, err := s.repo.GetBaseCurrency(ctx)
			if err != nil && !errors.IsNotFound(err) {
					return nil, err
			}
			if baseCurrency != nil {
					return nil, errors.NewBadRequestError("Ya existe una moneda base")
			}
	}

	currency := &Currency{
			Name:      req.Name,
			Code:      req.Code,
			Symbol:    req.Symbol,
			IsBase:    req.IsBase,
			Active:    req.Active,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateCurrency(ctx, currency); err != nil {
			return nil, err
	}

	return &CurrencyResponse{
			ID:        currency.ID,
			Name:      currency.Name,
			Code:      currency.Code,
			Symbol:    currency.Symbol,
			IsBase:    currency.IsBase,
			Active:    currency.Active,
			CreatedAt: currency.CreatedAt,
			UpdatedAt: currency.UpdatedAt,
	}, nil
}

func (s *Service) UpdateCurrency(ctx context.Context, id int64, req *UpdateCurrencyRequest) (*CurrencyResponse, error) {
	currency, err := s.repo.GetCurrencyByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			currency.Name = *req.Name
	}
	if req.Code != nil {
			currency.Code = *req.Code
	}
	if req.Symbol != nil {
			currency.Symbol = *req.Symbol
	}
	if req.Active != nil {
			currency.Active = *req.Active
	}
	if req.IsBase != nil {
			// Si se está cambiando a moneda base
			if !*req.IsBase && currency.IsBase {				
				return nil, errors.NewMessageError("Tiene que haber al menos una moneda base")
			}
			if *req.IsBase && !currency.IsBase {
				if err := s.repo.SetBaseCurrency(ctx, id); err != nil {
					return nil, err
				}
			}
			currency.IsBase = *req.IsBase
	}

	currency.UpdatedAt = time.Now()
	if err := s.repo.UpdateCurrency(ctx, currency); err != nil {
			return nil, err
	}

	return &CurrencyResponse{
			ID:        currency.ID,
			Name:      currency.Name,
			Code:      currency.Code,
			Symbol:    currency.Symbol,
			IsBase:    currency.IsBase,
			Active:    currency.Active,
			CreatedAt: currency.CreatedAt,
			UpdatedAt: currency.UpdatedAt,
	}, nil
}

func (s *Service) DeleteCurrency(ctx context.Context, id int64) error {
	return s.repo.DeleteCurrency(ctx, id)
}

func (s *Service) GetCurrency(ctx context.Context, id int64) (*CurrencyResponse, error) {
	currency, err := s.repo.GetCurrencyByID(ctx, id)
	if err != nil {
			return nil, err
	}

	return &CurrencyResponse{
			ID:        currency.ID,
			Name:      currency.Name,
			Code:      currency.Code,
			Symbol:    currency.Symbol,
			IsBase:    currency.IsBase,
			Active:    currency.Active,
			CreatedAt: currency.CreatedAt,
			UpdatedAt: currency.UpdatedAt,
	}, nil
}

func (s *Service) ListCurrencies(ctx context.Context, p *Pagination) ([]CurrencyResponse, int64, error) {
	currencies, total, err := s.repo.ListCurrencies(ctx, p)
	if err != nil {
			return nil, 0, err
	}

	response := make([]CurrencyResponse, len(currencies))
	for i, currency := range currencies {
			response[i] = CurrencyResponse{
					ID:        currency.ID,
					Name:      currency.Name,
					Code:      currency.Code,
					Symbol:    currency.Symbol,
					IsBase:    currency.IsBase,
					Active:    currency.Active,
					CreatedAt: currency.CreatedAt,
					UpdatedAt: currency.UpdatedAt,
			}
	}

	return response, total, nil
}

// Exchange rate operations
func (s *Service) CreateExchangeRate(ctx context.Context, req *CreateExchangeRateRequest) (*ExchangeRateResponse, error) {
	// Verificar que las monedas existan
	fromCurrency, err := s.repo.GetCurrencyByID(ctx, req.FromCurrencyID)
	if err != nil {
			return nil, err
	}
	toCurrency, err := s.repo.GetCurrencyByID(ctx, req.ToCurrencyID)
	if err != nil {
			return nil, err
	}

	// No permitir tipo de cambio entre la misma moneda
	if req.FromCurrencyID == req.ToCurrencyID {
			return nil, errors.NewBadRequestError("No se puede crear tipo de cambio entre la misma moneda")
	}

	// Crear el tipo de cambio
	rate := &ExchangeRate{
			FromCurrencyID: req.FromCurrencyID,
			ToCurrencyID:   req.ToCurrencyID,
			Rate:           req.Rate,
			FromCurrency:   fromCurrency,
			ToCurrency:     toCurrency,
			CreatedAt:      time.Now(),
	}

	if err := s.repo.CreateExchangeRate(ctx, rate); err != nil {
			return nil, err
	}

	// Crear el tipo de cambio inverso automáticamente
	inverseRate := &ExchangeRate{
			FromCurrencyID: req.ToCurrencyID,
			ToCurrencyID:   req.FromCurrencyID,
			Rate:           1 / req.Rate,
			FromCurrency:   toCurrency,
			ToCurrency:     fromCurrency,
			CreatedAt:      time.Now(),
	}

	if err := s.repo.CreateExchangeRate(ctx, inverseRate); err != nil {
			return nil, err
	}

	return &ExchangeRateResponse{
			ID:           rate.ID,
			FromCurrency: rate.FromCurrency,
			ToCurrency:   rate.ToCurrency,
			Rate:         rate.Rate,
			CreatedAt:    rate.CreatedAt,
	}, nil
}

func (s *Service) ConvertAmount(ctx context.Context, req *ConvertAmountRequest) (*ConvertAmountResponse, error) {
	// Obtener monedas
	fromCurrency, err := s.repo.GetCurrencyByCode(ctx, req.FromCurrencyCode)
	if err != nil {
			return nil, err
	}
	toCurrency, err := s.repo.GetCurrencyByCode(ctx, req.ToCurrencyCode)
	if err != nil {
			return nil, err
	}

	// Si son la misma moneda, retornar el mismo monto
	if fromCurrency.ID == toCurrency.ID {
			return &ConvertAmountResponse{
					OriginalAmount:   req.Amount,
					ConvertedAmount:  req.Amount,
					FromCurrency:     *fromCurrency,
					ToCurrency:       *toCurrency,
					ExchangeRate:     1,
					ConvertedAt:      time.Now(),
			}, nil
	}

	// Obtener tipo de cambio
	rate, err := s.repo.GetLatestExchangeRate(ctx, fromCurrency.ID, toCurrency.ID)
	if err != nil {
			return nil, err
	}

	return &ConvertAmountResponse{
			OriginalAmount:   req.Amount,
			ConvertedAmount:  req.Amount * rate.Rate,
			FromCurrency:     *fromCurrency,
			ToCurrency:       *toCurrency,
			ExchangeRate:     rate.Rate,
			ConvertedAt:      time.Now(),
	}, nil
}

func (s *Service) GetExchangeRates(ctx context.Context, currencyID int64) ([]ExchangeRate, error) {
	return s.repo.ListExchangeRates(ctx, currencyID)
}

func (s *Service) GetBaseCurrency(ctx context.Context) (*CurrencyResponse, error) {
	currency, err := s.repo.GetBaseCurrency(ctx)
	if err != nil {
			return nil, err
	}

	return &CurrencyResponse{
			ID:        currency.ID,
			Name:      currency.Name,
			Code:      currency.Code,
			Symbol:    currency.Symbol,
			IsBase:    currency.IsBase,
			Active:    currency.Active,
			CreatedAt: currency.CreatedAt,
			UpdatedAt: currency.UpdatedAt,
	}, nil
}