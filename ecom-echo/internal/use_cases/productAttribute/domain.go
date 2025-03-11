package productattribute

import (
	attributeValue "ecom/internal/use_cases/attributeValue"
	"time"
)

type ProductAttribute struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Description string    `json:"description"`   
    Values      []attributeValue.AttributeValue `json:"values"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
