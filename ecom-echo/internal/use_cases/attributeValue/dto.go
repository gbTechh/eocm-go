package attributevalue

// CreateAttributeValue estructura para crear un nuevo attribute value
type CreateAttributeValueRequest struct {
	Name        				string `json:"name" validate:"required,min=2,max=200,name"`
	IDProductAttribute  int64  `json:"id_product_attribute" validate:"required"`
}

// AttributeValueResponse estructura para respuestas
type AttributeValueResponse struct {
	ID          							int64      `json:"id"`
	Name        							string     `json:"name"`
	IDProductAttribute        int64      `json:"id_product_attribute"`	
}

// UpdateAttributeValueRequest estructura para actualización
type UpdateAttributeValueRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,name"`
}

// Pagination estructura para paginación y búsqueda
type Pagination struct {
	Page     int    `query:"page" validate:"min=1"`
	PerPage  int    `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
}