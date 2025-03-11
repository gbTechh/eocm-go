// service.go
package customer

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

// Customer methods
func (s *Service) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*CustomerResponse, error) {
	// Verificar si ya existe un cliente con ese email
	existingCustomer, err := s.repo.GetCustomerByEmail(ctx, req.Email)
	if err != nil && !errors.IsNotFound(err) {
			return nil, err
	}
	if existingCustomer != nil {
			return nil, errors.NewBadRequestError("Ya existe un cliente con ese email")
	}

	customer := &Customer{
			Name:      req.Name,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
			Active:    req.Active,
			IsGuest:   req.IsGuest,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
	}

	// Si hay grupos, prepararlos
	if len(req.GroupIDs) > 0 {
			customer.Groups = make([]Group, len(req.GroupIDs))
			for i, id := range req.GroupIDs {
					customer.Groups[i] = Group{ID: id}
			}
	}

	if err := s.repo.CreateCustomer(ctx, customer); err != nil {
			return nil, err
	}

	return s.GetCustomerByID(ctx, customer.ID)
}

func (s *Service) GetCustomerByID(ctx context.Context, id int64) (*CustomerResponse, error) {
	customer, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
			return nil, err
	}

	return mapCustomerToResponse(customer), nil
}

func (s *Service) GetCustomerByEmail(ctx context.Context, email string) (*CustomerResponse, error) {
	customer, err := s.repo.GetCustomerByEmail(ctx, email)
	if err != nil {
			return nil, err
	}

	return mapCustomerToResponse(customer), nil
}

func (s *Service) UpdateCustomer(ctx context.Context, id int64, req *UpdateCustomerRequest) (*CustomerResponse, error) {
	customer, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			customer.Name = *req.Name
	}
	if req.LastName != nil {
			customer.LastName = *req.LastName
	}
	if req.Email != nil && *req.Email != customer.Email {
			// Verificar si el nuevo email ya está en uso
			existingCustomer, err := s.repo.GetCustomerByEmail(ctx, *req.Email)
			if err != nil && !errors.IsNotFound(err) {
					return nil, err
			}
			if existingCustomer != nil {
					return nil, errors.NewBadRequestError("Email ya está en uso")
			}
			customer.Email = *req.Email
	}
	if req.Phone != nil {
			customer.Phone = *req.Phone
	}
	if req.Active != nil {
			customer.Active = *req.Active
	}
	if req.IsGuest != nil {
			customer.IsGuest = *req.IsGuest
	}

	// Actualizar grupos si se proporcionaron
	if req.GroupIDs != nil {
			customer.Groups = make([]Group, len(req.GroupIDs))
			for i, id := range req.GroupIDs {
					customer.Groups[i] = Group{ID: id}
			}
	}

	customer.UpdatedAt = time.Now()
	if err := s.repo.UpdateCustomer(ctx, customer); err != nil {
			return nil, err
	}

	return s.GetCustomerByID(ctx, id)
}

func (s *Service) DeleteCustomer(ctx context.Context, id int64) error {
	return s.repo.DeleteCustomer(ctx, id)
}

func (s *Service) ListCustomers(ctx context.Context, p *Pagination) ([]CustomerResponse, int64, error) {
	customers, total, err := s.repo.ListCustomers(ctx, p)
	if err != nil {
			return nil, 0, err
	}

	var response []CustomerResponse
	for _, customer := range customers {
			response = append(response, *mapCustomerToResponse(&customer))
	}

	return response, total, nil
}

// Group methods
func (s *Service) CreateGroup(ctx context.Context, req *CreateGroupRequest) (*Group, error) {
	group := &Group{
			Name:        req.Name,
			Description: req.Description,
			IDPriceList: req.IDPriceList,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateGroup(ctx, group); err != nil {
			return nil, err
	}

	return s.GetGroupByID(ctx, group.ID)
}

func (s *Service) GetGroupByID(ctx context.Context, id int64) (*Group, error) {
	return s.repo.GetGroupByID(ctx, id)
}

func (s *Service) UpdateGroup(ctx context.Context, id int64, req *UpdateGroupRequest) (*Group, error) {
	group, err := s.repo.GetGroupByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			group.Name = *req.Name
	}
	if req.Description != nil {
			group.Description = *req.Description
	}
	if req.IDPriceList != nil {
			group.IDPriceList = req.IDPriceList
	}

	group.UpdatedAt = time.Now()
	if err := s.repo.UpdateGroup(ctx, group); err != nil {
			return nil, err
	}

	return s.GetGroupByID(ctx, id)
}

func (s *Service) DeleteGroup(ctx context.Context, id int64) error {
	return s.repo.DeleteGroup(ctx, id)
}

func (s *Service) ListGroups(ctx context.Context, p *Pagination) ([]Group, int64, error) {
	return s.repo.ListGroups(ctx, p)
}

// Customer-Group operations
func (s *Service) AssignCustomersToGroup(ctx context.Context, groupID int64, customerIDs []int64) error {
	// Verificar que el grupo existe
	group, err := s.repo.GetGroupByID(ctx, groupID)
	if err != nil {
			return err
	}
	if group == nil {
			return errors.NewNotFoundError("Grupo no encontrado")
	}

	return s.repo.AssignCustomersToGroup(ctx, groupID, customerIDs)
}

func (s *Service) RemoveCustomersFromGroup(ctx context.Context, groupID int64, customerIDs []int64) error {
	return s.repo.RemoveCustomersFromGroup(ctx, groupID, customerIDs)
}

func (s *Service) GetCustomerGroups(ctx context.Context, customerID int64) ([]Group, error) {
	return s.repo.GetCustomerGroups(ctx, customerID)
}

// Helper functions
func mapCustomerToResponse(c *Customer) *CustomerResponse {
	return &CustomerResponse{
			ID:        c.ID,
			Name:      c.Name,
			LastName:  c.LastName,
			Email:     c.Email,
			Phone:     c.Phone,
			Active:    c.Active,
			IsGuest:   c.IsGuest,
			Groups:    c.Groups,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
	}
}