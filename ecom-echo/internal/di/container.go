// internal/di/container.go
package di

import (
	"database/sql"
	"ecom/config"
	"ecom/internal/shared/database"
	attributeValue "ecom/internal/use_cases/attributeValue"
	"ecom/internal/use_cases/category"
	currencies "ecom/internal/use_cases/currencies"
	customer "ecom/internal/use_cases/customer"
	"ecom/internal/use_cases/media"
	prices "ecom/internal/use_cases/prices"
	product "ecom/internal/use_cases/product"
	productattribute "ecom/internal/use_cases/productAttribute"
	"ecom/internal/use_cases/tag"
	taxcategory "ecom/internal/use_cases/taxes/taxCategory"
	taxrate "ecom/internal/use_cases/taxes/taxRate"
	taxrule "ecom/internal/use_cases/taxes/taxRule"
	taxzone "ecom/internal/use_cases/taxes/taxZone"
	"ecom/internal/use_cases/zone"
	// otros imports
)

type Container struct {
    db *sql.DB
    
    // Repositories
    mediaRepo media.Repository
    tagRepo tag.Repository
    categoryRepo category.Repository
    productattributeRepo productattribute.Repository
    attributeValueRepo attributeValue.Repository
    zoneRepo zone.Repository
    taxcategoryRepo taxcategory.Repository
    taxzoneRepo taxzone.Repository
    taxrateRepo taxrate.Repository
    taxruleRepo taxrule.Repository
    productRepo product.Repository
    pricesRepo prices.Repository
    customerRepo customer.Repository
    currenciesRepo currencies.Repository
    // Services
    mediaService *media.Service
    tagService *tag.Service
    categoryService *category.Service
    productattributeService *productattribute.Service
    attributeValueService *attributeValue.Service
    zoneService *zone.Service
    taxcategoryService *taxcategory.Service
    taxzoneService *taxzone.Service
    taxrateService *taxrate.Service
    taxruleService *taxrule.Service
    productService *product.Service
    pricesService *prices.Service
    customerService *customer.Service
    currenciesService *currencies.Service
}

func NewContainer(cfg *config.Config) (*Container, error) {
    c := &Container{}
    
    if err := c.initDB(cfg); err != nil {
        return nil, err
    }
    
    if err := c.initRepositories(); err != nil {
        return nil, err
    }
    
    if err := c.initServices(); err != nil {
        return nil, err
    }
    
    return c, nil
}

func (c *Container) initDB(cfg *config.Config) error {
    db, err := database.NewMySQLConnection(cfg)
		
    if err != nil {
        return err
    }
    c.db = db
    return nil
}

// Cleanup método para cerrar conexiones
func (c *Container) Cleanup() error {
    if c.db != nil {
        return c.db.Close()
    }
    return nil
}

// Agregar getter para DB
func (c *Container) DB() *sql.DB {
    return c.db
}

func (c *Container) initRepositories() error {
    c.mediaRepo = media.NewMySQLRepository(c.db)
    c.tagRepo = tag.NewMySQLRepository(c.db)
    c.categoryRepo = category.NewMySQLRepository(c.db)
    c.productattributeRepo = productattribute.NewMySQLRepository(c.db)
    c.attributeValueRepo = attributeValue.NewMySQLRepository(c.db)
    c.zoneRepo = zone.NewMySQLRepository(c.db)
    c.taxcategoryRepo = taxcategory.NewMySQLRepository(c.db)
    c.taxzoneRepo = taxzone.NewMySQLRepository(c.db)
    c.taxrateRepo = taxrate.NewMySQLRepository(c.db)
    c.taxruleRepo = taxrule.NewMySQLRepository(c.db)
    c.productRepo = product.NewMySQLRepository(c.db)
    c.pricesRepo = prices.NewMySQLRepository(c.db)
    c.customerRepo = customer.NewMySQLRepository(c.db)
    c.currenciesRepo = currencies.NewMySQLRepository(c.db)
    return nil
}

func (c *Container) initServices() error {
    uploadStrategy := media.NewLocalUploadStrategy("./uploads")
    
    // Inicializa el servicio de medios con el repositorio y la estrategia de subida
    c.mediaService = media.NewService(c.mediaRepo, uploadStrategy)
    c.tagService = tag.NewService(c.tagRepo)
    c.categoryService = category.NewService(c.categoryRepo)
    c.productattributeService = productattribute.NewService(c.productattributeRepo)
    c.attributeValueService = attributeValue.NewService(c.attributeValueRepo)
    c.zoneService = zone.NewService(c.zoneRepo)
    c.taxcategoryService = taxcategory.NewService(c.taxcategoryRepo)
    c.taxzoneService = taxzone.NewService(c.taxzoneRepo)
    c.taxrateService = taxrate.NewService(c.taxrateRepo)
    c.taxruleService = taxrule.NewService(c.taxruleRepo)
    c.productService = product.NewService(c.productRepo, c.taxcategoryService)
    c.pricesService = prices.NewService(c.pricesRepo)
    c.customerService = customer.NewService(c.customerRepo)
    c.currenciesService = currencies.NewService(c.currenciesRepo)
    return nil
}

// Getters
// func (c *Container) ProductService() *product.Service {
//     return c.productService
// }
func (c *Container) MediaService() *media.Service {
    return c.mediaService
}

func (c *Container) TagService() *tag.Service {
    return c.tagService
}
func (c *Container) CategoryService() *category.Service {
    return c.categoryService
}
func (c *Container) ProductAttributeService() *productattribute.Service {
    return c.productattributeService
}
func (c *Container) AttributeValueService() *attributeValue.Service {
    return c.attributeValueService
}
func (c *Container) ZoneService() *zone.Service {
    return c.zoneService
}
func (c *Container) TaxCategoryService() *taxcategory.Service {
    return c.taxcategoryService
}
func (c *Container) TaxZoneService() *taxzone.Service {
    return c.taxzoneService
}
func (c *Container) TaxRateService() *taxrate.Service {
    return c.taxrateService
}
func (c *Container) TaxRuleService() *taxrule.Service {
    return c.taxruleService
}
func (c *Container) ProductService() *product.Service {
    return c.productService
}
func (c *Container) PricesService() *prices.Service {
    return c.pricesService
}
func (c *Container) CustomerService() *customer.Service {
    return c.customerService
}
func (c *Container) CurrencyService() *currencies.Service {
    return c.currenciesService
}



