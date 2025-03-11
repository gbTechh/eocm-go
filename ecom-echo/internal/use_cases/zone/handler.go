// handler.go
package zone

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
	var req CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	zone, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Zone creado exitosamente", zone))
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	zone, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Zone obtenido exitosamente", zone))
}

func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	zone, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("zone actualizado exitosamente", zone))
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
			response.Success("Zone eliminado exitosamente", nil))
}

func (h *Handler) List(c echo.Context) error {
	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error en paginación", err.Error()))
	}

	zones, total, err := h.service.List(c.Request().Context(), &p)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Zonas obtenidos exitosamente", 
					map[string]interface{}{
							"items": zones,
							"total": total,
							"page": p.Page,
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
							response.Error("Zona no encontrado", nil))
			}
			return c.JSON(http.StatusInternalServerError, 
					response.Error("Error interno del servidor", nil))
	}
}