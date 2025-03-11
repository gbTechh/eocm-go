// service.go
package product

import (
	"context"
	"ecom/internal/shared/errors"
	"fmt"
	"hash/fnv"
	"strings"
	"time"
	"unicode"

	at "ecom/internal/use_cases/attributeValue"
	tag "ecom/internal/use_cases/tag"
	taxcategory "ecom/internal/use_cases/taxes/taxCategory"
)

type Service struct {
	repo Repository
	taxCategoryService *taxcategory.Service
}

func NewService(repo Repository, taxCategoryService *taxcategory.Service) *Service {
	return &Service{
		repo: repo,
		taxCategoryService: taxCategoryService,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	product := &Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Status:      req.Status,
		IDCategory:  req.IDCategory,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Si hay tags, prepararlos
	if len(req.TagIDs) > 0 {
		product.Tags = make([]tag.Tag, len(req.TagIDs))
		for i, tagID := range req.TagIDs {
			product.Tags[i] = tag.Tag{ID: tagID}
		}
	}

	// Crear producto base
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	idTaxDefault, err := s.taxCategoryService.GetTaxDefault(ctx)
	if err != nil {
		return nil, err
	}

	// Crear variante por defecto
	variant := &ProductVariant{
		Name:          defaultValue(req.DefaultVariant.Name, product.Name),
		SKU:           generateSKU(req.DefaultVariant.SKU, product.Name,product.Name,product.ID ),
		Stock:         req.DefaultVariant.Stock,
		IDProduct:     product.ID,
		IDTaxCategory: defaultTaxCategory(req.DefaultVariant.IDTaxCategory, idTaxDefault.ID),
		AttributeValues: make([]at.AttributeValue, len(req.DefaultVariant.AttributeValues)),
	}


	// Mapear attribute values
	for i, av := range req.DefaultVariant.AttributeValues {
			variant.AttributeValues[i] = at.AttributeValue{
					ID:                 av.ID,
			}
	}

	if err := s.repo.CreateVariant(ctx, variant); err != nil {
			return nil, err
	}

	return s.GetByID(ctx, product.ID)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {		
		return nil, err
	}

	return mapProductToResponse(product), nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateProductRequest) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
			return nil, err
	}

	if req.Name != nil {
			product.Name = *req.Name
			product.Slug = *req.Slug
	}
	if req.Description != nil {
			product.Description = *req.Description
	}
	if req.Status != nil {
			product.Status = *req.Status
	}
	if req.IDCategory != nil {
			product.IDCategory = *req.IDCategory
	}
	product.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, product); err != nil {
			return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	// Primero obtener todas las variantes y eliminarlas
	variants, err := s.repo.ListVariants(ctx, id)
	if err != nil {
			return err
	}

	for _, variant := range variants {
			if err := s.repo.DeleteVariant(ctx, variant.ID); err != nil {
					return err
			}
	}

	// Luego eliminar el producto
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, p *Pagination) ([]ProductResponse, int64, error) {
	products, total, err := s.repo.List(ctx, p)
	if err != nil {
			return nil, 0, err
	}

	var response []ProductResponse
	for _, product := range products {
			response = append(response, *mapProductToResponse(&product))
	}

	return response, total, nil
}

func (s *Service) CreateVariant(ctx context.Context, req *CreateVariantRequest) (*ProductResponse, error) {
	// Verificar que el producto existe
	_, err := s.repo.GetByID(ctx, req.IDProduct)
	if err != nil {
		return nil, err
	}

	idTaxDefault, err := s.taxCategoryService.GetTaxDefault(ctx)
	if err != nil {
		return nil, err
	}
	
	variant := &ProductVariant{
		Name:          req.Name,
		SKU:           req.SKU,
		Stock:         req.Stock,
		IDProduct:     req.IDProduct,
		IDTaxCategory: defaultTaxCategory(req.IDTaxCategory, idTaxDefault.ID),
		AttributeValues: make([]at.AttributeValue, len(req.AttributeValues)),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Mapear attribute values
	for i, av := range req.AttributeValues {
		variant.AttributeValues[i] = at.AttributeValue{
			ID:                 av.ID,
		}
	}

	if err := s.repo.CreateVariant(ctx, variant); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, req.IDProduct)
}

func (s *Service) UpdateVariant(ctx context.Context, id int64, req *UpdateVariantRequest) (*ProductVariant, error) {
	variant, err := s.repo.GetVariantByID(ctx, id)
	if err != nil {		
		return nil, err
	}

	if req.Name != nil {
			variant.Name = *req.Name
	}
	if req.SKU != nil {
			variant.SKU = *req.SKU
	}
	if req.Stock != nil {
			variant.Stock = *req.Stock
	}

	if req.IDTaxCategory != nil {
		idTaxDefault, err := s.taxCategoryService.GetTaxDefault(ctx)
		if err != nil {
			return nil, err
		}
		variant.IDTaxCategory = defaultTaxCategory(*req.IDTaxCategory, idTaxDefault.ID)
	}

	// Actualizar AttributeValues si se proporcionaron
	if req.AttributeValues != nil {
			variant.AttributeValues = make([]at.AttributeValue, len(*req.AttributeValues))
			for i, av := range *req.AttributeValues {
					variant.AttributeValues[i] = at.AttributeValue{
							ID: av.ID,
					}
			}
	}
	
	variant.UpdatedAt = time.Now()

	fmt.Printf("variante: %v\n",variant.IDTaxCategory)

	if err := s.repo.UpdateVariant(ctx, variant); err != nil {
		return nil, err
	}

	return s.repo.GetVariantByID(ctx, id)
}

func (s *Service) DeleteVariant(ctx context.Context, id int64) error {
	variant, err := s.repo.GetVariantByID(ctx, id)
	if err != nil {
			return err
	}

	// No permitir eliminar la última variante de un producto
	variants, err := s.repo.ListVariants(ctx, variant.IDProduct)
	if err != nil {
			return err
	}

	if len(variants) <= 1 {
			return errors.NewBadRequestError("No se puede eliminar la única variante del producto")
	}

	return s.repo.DeleteVariant(ctx, id)
}

func (s *Service) GetVariant(ctx context.Context, id int64) (*ProductVariant, error) {
	return s.repo.GetVariantByID(ctx, id)
}

func (s *Service) ListVariants(ctx context.Context, productID int64) ([]ProductVariant, error) {
	return s.repo.ListVariants(ctx, productID)
}



func mapProductToResponse(p *Product) *ProductResponse {

	return &ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			Status:      p.Status,
			IDCategory:  p.IDCategory,
			Category:    p.Category,
			Tags:        p.Tags,
			Variants:    p.Variants,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
	}
}

func defaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func generateSKU(value string, productName string, variantName string, productID int64) (string) {
	// Limpiar y preparar el nombre del producto
	if value != "" {
		return value;
	}
	cleanName := strings.ToUpper(productName)
	cleanName = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
				return r
		}
		return -1
	}, cleanName)

	// Tomar las primeras 3 letras del nombre del producto
	prefix := "PRD"
	if len(cleanName) >= 3 {
		prefix = cleanName[:3]
	}

	// Obtener hash del nombre de la variante para hacerlo único
	hash := fnv.New32a()
	hash.Write([]byte(variantName))
	variantHash := hash.Sum32() % 1000 // Limitar a 3 dígitos

	// Generar timestamp en formato corto (últimos 4 dígitos)
	timestamp := time.Now().UnixNano() / 1000000 % 10000

	// Construir el SKU
	sku := fmt.Sprintf("%s-%d-%03d-%04d", 
			prefix,           // Prefijo del producto (3 letras)
			productID,        // ID del producto
			variantHash,      // Hash de la variante (3 dígitos)
			timestamp)        // Timestamp (4 dígitos)

	return sku
}

func defaultTaxCategory(value, id int64) int64 {
	if id == 0 {
		return id
	}
	if(value > 0) {
		return value
	}
	return id
}

