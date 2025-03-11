// handler.go
package currency

import (
	"ecom/internal/shared/errors"
	"ecom/internal/shared/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Group, s *Service) {
	h := &Handler{service: s}

	// Currency routes
	e.POST("", h.CreateCurrency)
	e.GET("", h.ListCurrencies)
	e.GET("/:id", h.GetCurrency)
	e.PUT("/:id", h.UpdateCurrency)
	e.DELETE("/:id", h.DeleteCurrency)
	e.GET("/base", h.GetBaseCurrency)

	// Exchange rate routes
	e.POST("/rates", h.CreateExchangeRate)
	e.GET("/rates/:currency_id", h.GetExchangeRates)
	e.POST("/convert", h.ConvertAmount)
}

func (h *Handler) CreateCurrency(c echo.Context) error {
	var req CreateCurrencyRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	currency, err := h.service.CreateCurrency(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Moneda creada exitosamente", currency))
}

func (h *Handler) GetCurrency(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	currency, err := h.service.GetCurrency(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Moneda obtenida exitosamente", currency))
}

func (h *Handler) GetBaseCurrency(c echo.Context) error {
	currency, err := h.service.GetBaseCurrency(c.Request().Context())
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Moneda base obtenida exitosamente", currency))
}

func (h *Handler) UpdateCurrency(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdateCurrencyRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	currency, err := h.service.UpdateCurrency(c.Request().Context(), id, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Moneda actualizada exitosamente", currency))
}

func (h *Handler) DeleteCurrency(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	if err := h.service.DeleteCurrency(c.Request().Context(), id); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Moneda eliminada exitosamente", nil))
}

func (h *Handler) ListCurrencies(c echo.Context) error {
	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error en paginación", err.Error()))
	}

	currencies, total, err := h.service.ListCurrencies(c.Request().Context(), &p)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Monedas obtenidas exitosamente", 
					map[string]interface{}{
							"items":    currencies,
							"total":    total,
							"page":     p.Page,
							"per_page": p.PerPage,
					}))
}

func (h *Handler) CreateExchangeRate(c echo.Context) error {
	var req CreateExchangeRateRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	rate, err := h.service.CreateExchangeRate(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Tipo de cambio creado exitosamente", rate))
}

func (h *Handler) GetExchangeRates(c echo.Context) error {
	currencyID, err := strconv.ParseInt(c.Param("currency_id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID de moneda inválido", err.Error()))
	}

	rates, err := h.service.GetExchangeRates(c.Request().Context(), currencyID)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Tipos de cambio obtenidos exitosamente", rates))
}

func (h *Handler) ConvertAmount(c echo.Context) error {
	var req ConvertAmountRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	result, err := h.service.ConvertAmount(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Conversión realizada exitosamente", result))
}

func (h *Handler) handleError(c echo.Context, err error) error {
	return errors.HandleError(c, err)
}