// handler.go
package product

import (
	"ecom/internal/shared/errors"
	"ecom/internal/shared/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Group, s *Service) {
	h := &Handler{service: s}

	// Rutas de productos
	e.POST("", h.Create)
	e.GET("", h.List)
	e.GET("/:id", h.GetByID)
	e.PUT("/:id", h.Update)
	e.DELETE("/:id", h.Delete)

	// Rutas de variantes
	e.POST("/variant", h.CreateVariant)
	e.GET("/variant/:id", h.GetVariant) // obtener una variante
	e.PUT("/variant/:id", h.UpdateVariant)
	e.DELETE("/variant/:id", h.DeleteVariant)
	e.GET("/:id/variants", h.ListVariants) //obtener todas las variatnes de un producto
}

// Handlers de Productos
func (h *Handler) Create(c echo.Context) error {
	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	product, err := h.service.Create(c.Request().Context(), &req)

	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Producto creado exitosamente", product))
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	product, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Producto obtenido exitosamente", product))
}

func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdateProductRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	product, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Producto actualizado exitosamente", product))
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
			response.Success("Producto eliminado exitosamente", nil))
}

func (h *Handler) List(c echo.Context) error {
	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error en paginación", err.Error()))
	}

	products, total, err := h.service.List(c.Request().Context(), &p)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Productos obtenidos exitosamente", 
					map[string]interface{}{
							"items":    products,
							"total":    total,
							"page":     p.Page,
							"per_page": p.PerPage,
					}))
}

// Handlers de Variantes
func (h *Handler) CreateVariant(c echo.Context) error {
	var req CreateVariantRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	product, err := h.service.CreateVariant(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Variante creada exitosamente", product))
}

func (h *Handler) GetVariant(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	variant, err := h.service.GetVariant(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Variante obtenida exitosamente", variant))
}

func (h *Handler) UpdateVariant(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdateVariantRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	variant, err := h.service.UpdateVariant(c.Request().Context(), id, &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Variante actualizada exitosamente", variant))
}

func (h *Handler) DeleteVariant(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	if err := h.service.DeleteVariant(c.Request().Context(), id); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Variante eliminada exitosamente", nil))
}

func (h *Handler) ListVariants(c echo.Context) error {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID de producto inválido", err.Error()))
	}

	variants, err := h.service.ListVariants(c.Request().Context(), productID)
	if err != nil {
		fmt.Printf("error: %v\n,", err)
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Variantes obtenidas exitosamente", variants))
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