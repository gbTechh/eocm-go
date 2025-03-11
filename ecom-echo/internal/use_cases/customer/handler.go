// handler.go
package customer

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

	// Customer routes
	e.POST("", h.CreateCustomer)
	e.GET("", h.ListCustomers)
	e.GET("/:id", h.GetCustomer)
	e.GET("/email/:email", h.GetCustomerByEmail)
	e.PUT("/:id", h.UpdateCustomer)
	e.DELETE("/:id", h.DeleteCustomer)

	//TODO: Faltan testear todas estas rutas de abajo
	// Group routes
	groups := e.Group("/groups")
	groups.POST("", h.CreateGroup)
	groups.GET("", h.ListGroups)
	groups.GET("/:id", h.GetGroup)
	groups.PUT("/:id", h.UpdateGroup)
	groups.DELETE("/:id", h.DeleteGroup)

	// Customer-Group operations
	groups.POST("/:id/customers", h.AssignCustomers)
	groups.DELETE("/:id/customers", h.RemoveCustomers)
	e.GET("/:id/groups", h.GetCustomerGroups)
}

// Customer handlers
func (h *Handler) CreateCustomer(c echo.Context) error {
	var req CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	customer, err := h.service.CreateCustomer(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Cliente creado exitosamente", customer))
}

func (h *Handler) GetCustomer(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	customer, err := h.service.GetCustomerByID(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Cliente obtenido exitosamente", customer))
}

func (h *Handler) GetCustomerByEmail(c echo.Context) error {
	email := c.Param("email")
	customer, err := h.service.GetCustomerByEmail(c.Request().Context(), email)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Cliente obtenido exitosamente", customer))
}

func (h *Handler) UpdateCustomer(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdateCustomerRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	customer, err := h.service.UpdateCustomer(c.Request().Context(), id, &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Cliente actualizado exitosamente", customer))
}

func (h *Handler) DeleteCustomer(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	if err := h.service.DeleteCustomer(c.Request().Context(), id); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Cliente eliminado exitosamente", nil))
}

func (h *Handler) ListCustomers(c echo.Context) error {
	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error en paginación", err.Error()))
	}

	customers, total, err := h.service.ListCustomers(c.Request().Context(), &p)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Clientes obtenidos exitosamente", 
					map[string]interface{}{
							"items":    customers,
							"total":    total,
							"page":     p.Page,
							"per_page": p.PerPage,
					}))
}

// Group handlers
func (h *Handler) CreateGroup(c echo.Context) error {
	var req CreateGroupRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	group, err := h.service.CreateGroup(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Grupo creado exitosamente", group))
}

func (h *Handler) GetGroup(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	group, err := h.service.GetGroupByID(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Grupo obtenido exitosamente", group))
}

func (h *Handler) UpdateGroup(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	var req UpdateGroupRequest
	if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	group, err := h.service.UpdateGroup(c.Request().Context(), id, &req)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Grupo actualizado exitosamente", group))
}

func (h *Handler) DeleteGroup(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	if err := h.service.DeleteGroup(c.Request().Context(), id); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Grupo eliminado exitosamente", nil))
}

func (h *Handler) ListGroups(c echo.Context) error {
	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error en paginación", err.Error()))
	}

	groups, total, err := h.service.ListGroups(c.Request().Context(), &p)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Grupos obtenidos exitosamente", 
					map[string]interface{}{
							"items":    groups,
							"total":    total,
							"page":     p.Page,
							"per_page": p.PerPage,
					}))
}

// Customer-Group operations
func (h *Handler) AssignCustomers(c echo.Context) error {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID de grupo inválido", err.Error()))
	}

	var customerIDs []int64
	if err := c.Bind(&customerIDs); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := h.service.AssignCustomersToGroup(c.Request().Context(), groupID, customerIDs); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Clientes asignados exitosamente", nil))
}

func (h *Handler) RemoveCustomers(c echo.Context) error {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID de grupo inválido", err.Error()))
	}

	var customerIDs []int64
	if err := c.Bind(&customerIDs); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar la solicitud", err.Error()))
	}

	if err := h.service.RemoveCustomersFromGroup(c.Request().Context(), groupID, customerIDs); err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Clientes removidos exitosamente", nil))
}

func (h *Handler) GetCustomerGroups(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("ID inválido", err.Error()))
	}

	groups, err := h.service.GetCustomerGroups(c.Request().Context(), id)
	if err != nil {
			return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, 
			response.Success("Grupos del cliente obtenidos exitosamente", groups))
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