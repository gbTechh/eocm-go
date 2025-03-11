// internal/server/routes.go (NUEVO)
package server

import (
	attributeValue "ecom/internal/use_cases/attributeValue"
	"ecom/internal/use_cases/category"
	currencies "ecom/internal/use_cases/currencies"
	"ecom/internal/use_cases/customer"
	"ecom/internal/use_cases/media"
	"ecom/internal/use_cases/prices"
	product "ecom/internal/use_cases/product"
	productattribute "ecom/internal/use_cases/productAttribute"
	"ecom/internal/use_cases/tag"
	taxcategory "ecom/internal/use_cases/taxes/taxCategory"
	taxrate "ecom/internal/use_cases/taxes/taxRate"
	taxrule "ecom/internal/use_cases/taxes/taxRule"
	taxzone "ecom/internal/use_cases/taxes/taxZone"
	"ecom/internal/use_cases/zone"
	"fmt"
)

// otros imports necesarios

func (s *Server) registerRoutes() {
	// API versioning
		v1 := s.echo.Group("/api/v1")
	
	// Media routes - usa v1.Group("/media")
	media.NewHandler(v1.Group("/media"), s.container.MediaService())
	tag.NewHandler(v1.Group("/tags"), s.container.TagService())
	category.NewHandler(v1.Group("/categories"), s.container.CategoryService())
	productattribute.NewHandler(v1.Group("/product-attributes"), s.container.ProductAttributeService())
	attributeValue.NewHandler(v1.Group("/attributes-value"), s.container.AttributeValueService())
	zone.NewHandler(v1.Group("/zones"), s.container.ZoneService())
	taxcategory.NewHandler(v1.Group("/tax-categories"), s.container.TaxCategoryService())
	taxzone.NewHandler(v1.Group("/tax-zones"), s.container.TaxZoneService())
	taxrate.NewHandler(v1.Group("/tax-rates"), s.container.TaxRateService())
	taxrule.NewHandler(v1.Group("/tax-rules"), s.container.TaxRuleService())
	product.NewHandler(v1.Group("/products"), s.container.ProductService())
	prices.NewHandler(v1.Group("/prices"), s.container.PricesService())
	customer.NewHandler(v1.Group("/customers"), s.container.CustomerService())
	currencies.NewHandler(v1.Group("/currencies"), s.container.CurrencyService())
	
	// Product routes
	//product.NewHandler(v1.Group("/products"), s.container.ProductService())

	// Imprimir rutas registradas
	for _, route := range s.echo.Routes() {
		fmt.Printf("Method: %v, Path: %v\n", route.Method, route.Path)
	}
}