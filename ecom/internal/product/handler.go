// internal/product/handler.go
package product

import (
	"ecommerce/pkg/apperrors"
	"ecommerce/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// List devuelve una lista paginada de productos
func (h *Handler) List(c *gin.Context) {
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "10")
    
    products, err := h.service.List(c.Request.Context(), page, limit)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Success(c, http.StatusOK, products)
}

// Get obtiene un producto por su ID
func (h *Handler) Get(c *gin.Context) {
    id := c.Param("id")
    
    product, err := h.service.Get(c.Request.Context(), id)
    if err != nil {
        if err, ok := err.(*apperrors.Error); ok && err.Code == apperrors.ErrNotFound {
            response.Error(c, err)
            return
        }
        response.Error(c, err)
        return
    }
    
    response.Success(c, http.StatusOK, product)
}

// Create crea un nuevo producto
func (h *Handler) Create(c *gin.Context) {
    var input CreateProductInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.Error(c, apperrors.NewError(
            apperrors.ErrInvalidInput,
            "Datos de producto inválidos",
        ))
        return
    }
    
    product, err := h.service.Create(c.Request.Context(), &input)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Success(c, http.StatusCreated, product)
}

// Update actualiza un producto existente
func (h *Handler) Update(c *gin.Context) {
    id := c.Param("id")
    
    var input UpdateProductInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.Error(c, apperrors.NewError(
            apperrors.ErrInvalidInput,
            "Datos de actualización inválidos",
        ))
        return
    }
    
    product, err := h.service.Update(c.Request.Context(), id, &input)
    if err != nil {
        if err, ok := err.(*apperrors.Error); ok && err.Code == apperrors.ErrNotFound {
            response.Error(c, err)
            return
        }
        response.Error(c, err)
        return
    }
    
    response.Success(c, http.StatusOK, product)
}

// Delete elimina un producto
func (h *Handler) Delete(c *gin.Context) {
    id := c.Param("id")
    
    err := h.service.Delete(c.Request.Context(), id)
    if err != nil {
        if err, ok := err.(*apperrors.Error); ok && err.Code == apperrors.ErrNotFound {
            response.Error(c, err)
            return
        }
        response.Error(c, err)
        return
    }
    
    response.Success(c, http.StatusNoContent, nil)
}

// ListCategories devuelve todas las categorías de productos
func (h *Handler) ListCategories(c *gin.Context) {
    categories, err := h.service.ListCategories(c.Request.Context())
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Success(c, http.StatusOK, categories)
}