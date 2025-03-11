package attributevalue

type AttributeValue struct {
	ID                 int64             `json:"id"`
	Name               string            `json:"name"`
	IDProductAttribute int64             `json:"id_product_attribute"`
	ProductAttribute   *ProductAttributeInfo `json:"product_attribute,omitempty"`
}

// Esta estructura es solo para mostrar info básica y evitar recursión
type ProductAttributeInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}