// handler.go
package taxrule

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
	e.POST("", h.Create)
	e.GET("", h.List)
	e.GET("/:id", h.GetByID)
	e.PUT("/:id", h.Update)
	e.DELETE("/:id", h.Delete)
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateTaxRuleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	tz, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Regla de impuesto creada exitosamente", tz))
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	tz, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Regla de impuesto actualizada exitosamente", tz))
}

func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdateTaxRuleRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	tz, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Regla de impuesto actualizada exitosamente", tz))
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Regla de impuesto eliminada exitosamente", nil))
}

func (h *Handler) List(c echo.Context) error {
	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error en paginación", err.Error()))
	}

	tzs, total, err := h.service.List(c.Request().Context(), &p)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
        response.Success("Reglas de impuestos obtenidas exitosamente", 
            map[string]interface{}{
                "items":    tzs,
                "total":    total,
                "page":     p.Page,
                "per_page": p.PerPage,
            }))
}

func (h *Handler) handleError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *errors.AppError:
		return c.JSON(e.Code, response.Error(e.Message, nil))
	default:
		if errors.IsNotFound(err) {
				return c.JSON(http.StatusNotFound, 
						response.Error("Regla de impuesto no encontrado", nil))
		}
		return c.JSON(http.StatusInternalServerError, 
				response.Error("Error interno del servidor", nil))
	}
}