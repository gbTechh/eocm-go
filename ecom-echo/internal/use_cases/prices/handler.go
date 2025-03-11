// handler.go
package prices

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

	// PriceList routes
	e.POST("/lists", h.CreatePriceList)
	e.GET("/lists", h.ListPriceLists)
	e.GET("/lists/:id", h.GetPriceList)
	e.PUT("/lists/:id", h.UpdatePriceList)
	e.DELETE("/lists/:id", h.DeletePriceList)

	// Price routes
	e.POST("/lists/:id/value", h.CreatePrice)
	e.GET("/lists/:id/all-values", h.ListPrices) // id price list
	e.GET("/value/:id", h.GetPrice) // id price value
	e.PUT("/value/:id", h.UpdatePrice)
	e.DELETE("/value/:id", h.DeletePrice)

	// ProductVariantPrice routes
	e.POST("/assign", h.AssignPrice)
	e.DELETE("/variants/:variant_id/prices/:price_id", h.UnassignPrice)
	e.GET("/variant/:id/active-price", h.GetActivePrice)
	e.GET("/variant/:id/prices", h.GetAllPrices)
}

// PriceList handlers
func (h *Handler) CreatePriceList(c echo.Context) error {
	var req CreatePriceListRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	priceList, err := h.service.CreatePriceList(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Lista de precios creada exitosamente", priceList))
}

func (h *Handler) GetPriceList(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	priceList, err := h.service.GetPriceList(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Lista de precios obtenida exitosamente", priceList))
}

func (h *Handler) UpdatePriceList(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdatePriceListRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	priceList, err := h.service.UpdatePriceList(c.Request().Context(), id, &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Lista de precios actualizada exitosamente", priceList))
}

func (h *Handler) DeletePriceList(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	if err := h.service.DeletePriceList(c.Request().Context(), id); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Lista de precios eliminada exitosamente", nil))
}

func (h *Handler) ListPriceLists(c echo.Context) error {
	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error en paginación", err.Error()))
	}

	priceLists, total, err := h.service.ListPriceLists(c.Request().Context(), &p)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Listas de precios obtenidas exitosamente", 
					map[string]interface{}{
							"items":    priceLists,
							"total":    total,
							"page":     p.Page,
							"per_page": p.PerPage,
					}))
}

// Price handlers
func (h *Handler) CreatePrice(c echo.Context) error {
	priceListID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID de lista de precios inválido", err.Error()))
	}

	var req CreatePriceRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	req.IDPriceList = priceListID

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	price, err := h.service.CreatePrice(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Precio creado exitosamente", price))
}

func (h *Handler) GetPrice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	price, err := h.service.GetPrice(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Precio obtenido exitosamente", price))
}

func (h *Handler) UpdatePrice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdatePriceRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	price, err := h.service.UpdatePrice(c.Request().Context(), id, &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Precio actualizado exitosamente", price))
}

func (h *Handler) DeletePrice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	if err := h.service.DeletePrice(c.Request().Context(), id); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Precio eliminado exitosamente", nil))
}

func (h *Handler) ListPrices(c echo.Context) error {
	priceListID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID de lista de precios inválido", err.Error()))
	}

	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error en paginación", err.Error()))
	}

	prices, total, err := h.service.ListPrices(c.Request().Context(), priceListID, &p)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Precios obtenidos exitosamente", 
					map[string]interface{}{
							"items":    prices,
							"total":    total,
							"page":     p.Page,
							"per_page": p.PerPage,
					}))
}

// ProductVariantPrice handlers
func (h *Handler) AssignPrice(c echo.Context) error {
	var req AssignPriceRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	if err := h.service.AssignPrice(c.Request().Context(), &req); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Precio asignado exitosamente", nil))
}

func (h *Handler) UnassignPrice(c echo.Context) error {
	variantID, err := strconv.ParseInt(c.Param("variant_id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID de variante inválido", err.Error()))
	}

	priceID, err := strconv.ParseInt(c.Param("price_id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID de precio inválido", err.Error()))
	}

	if err := h.service.UnassignPrice(c.Request().Context(), variantID, priceID); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Precio desasignado exitosamente", nil))
}

func (h *Handler) GetActivePrice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	price, err := h.service.GetActivePrice(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Precio activo obtenido exitosamente", price))
}

func (h *Handler) GetAllPrices(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	prices, err := h.service.GetAllPrices(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Precios obtenidos exitosamente", prices))
}

func (h *Handler) handleError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *errors.AppError:
			return c.JSON(e.Code, response.Error(e.Message, nil))
	default:
			if errors.IsNotFound(err) {
					return c.JSON(http.StatusNotFound, 
							response.Error("Recurso no encontrado", nil))
			}
			return c.JSON(http.StatusInternalServerError, 
					response.Error("Error interno del servidor", nil))
	}
}