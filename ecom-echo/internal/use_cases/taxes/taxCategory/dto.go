package taxcategory

import "time"

type CreateTaxCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	IsDefault   bool   `json:"is_default"`
}

type TaxCategoryResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateTaxCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,name"`
	Description *string `json:"description,omitempty"`
	Status      *bool   `json:"status,omitempty"`
	IsDefault   *bool   `json:"is_default"`
}

type Pagination struct {
	Page     int    `query:"page" validate:"min=1"`
	PerPage  int    `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
	Status   *bool  `query:"status"`
}
