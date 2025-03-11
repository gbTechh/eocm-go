// dto.go
package customer

import "time"

type CreateCustomerRequest struct {
	Name     string  `json:"name" validate:"required,min=2,name"`
	LastName string  `json:"last_name" validate:"omitempty,name,min=2"`
	Email    string  `json:"email" validate:"required,email"`
	Phone    string  `json:"phone" validate:"omitempty,e164"`
	Active   bool    `json:"active"`
	IsGuest  bool    `json:"is_guest"`
	GroupIDs []int64 `json:"group_ids,omitempty"`
}

type UpdateCustomerRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=2"`
	LastName *string  `json:"last_name,omitempty"`
	Email    *string  `json:"email,omitempty" validate:"omitempty,email"`
	Phone    *string  `json:"phone,omitempty" validate:"omitempty,e164"`
	Active   *bool    `json:"active,omitempty"`
	IsGuest  *bool    `json:"is_guest,omitempty"`
	GroupIDs []int64  `json:"group_ids,omitempty"`
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	IDPriceList *int64 `json:"id_price_list,omitempty"`
}

type UpdateGroupRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Description *string `json:"description,omitempty"`
	IDPriceList *int64  `json:"id_price_list,omitempty"`
}

type CustomerResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Active    bool      `json:"active"`
	IsGuest   bool      `json:"is_guest"`
	Groups    []Group   `json:"groups,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupResponse struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IDPriceList *int64     `json:"id_price_list,omitempty"`
	PriceList   *PriceList `json:"price_list,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Pagination struct {
	Page     int    `query:"page" validate:"min=1"`
	PerPage  int    `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
	Active   *bool  `query:"active"`
	IsGuest  *bool  `query:"is_guest"`
}