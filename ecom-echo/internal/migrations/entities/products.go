package entities

import (
	"database/sql"
	"fmt"
)

type ProductsMigration struct {
    db *sql.DB
}

func NewProductsMigration(db *sql.DB) *ProductsMigration {
    return &ProductsMigration{db: db}
}

func (m *ProductsMigration) Migrate() error {
    // Crear tabla products
    createProductsQuery := `
        CREATE TABLE IF NOT EXISTS products (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(200) NOT NULL,
            slug VARCHAR(220) UNIQUE NOT NULL,
            description TEXT,
            status BOOLEAN,
            id_category INT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP NULL, 
            FOREIGN KEY (id_category) REFERENCES categories(id)
        )
    `
    if _, err := m.db.Exec(createProductsQuery); err != nil {
        return fmt.Errorf("error creating products table: %w", err)
    }

    // Crear tabla products_tags
    createProductsTagsQuery := `
        CREATE TABLE IF NOT EXISTS n_products_tags (
            id_tag INT,
            id_product INT,
            PRIMARY KEY (id_tag, id_product),
            FOREIGN KEY (id_tag) REFERENCES tags(id),
            FOREIGN KEY (id_product) REFERENCES products(id)
        )
    `
    if _, err := m.db.Exec(createProductsTagsQuery); err != nil {
        return fmt.Errorf("error creating products_tags table: %w", err)
    }

    // Crear tabla product_attributes
    createProductAttributesQuery := `
        CREATE TABLE IF NOT EXISTS product_attributes (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(200) UNIQUE NOT NULL,
            type VARCHAR(100) NOT NULL,
            description TEXT, 
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP NULL
        )
    `
    if _, err := m.db.Exec(createProductAttributesQuery); err != nil {
        return fmt.Errorf("error creating product_attributes table: %w", err)
    }

    // Crear tabla attribute_values
    createAttributeValuesQuery := `
        CREATE TABLE IF NOT EXISTS attribute_values (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(200) NOT NULL,
            id_attribute_product INT,
            FOREIGN KEY (id_attribute_product) REFERENCES product_attributes(id)
        )
    `
    if _, err := m.db.Exec(createAttributeValuesQuery); err != nil {
        return fmt.Errorf("error creating attribute_values table: %w", err)
    }

    // Crear tabla product_variants
    createProductVariantsQuery := `
        CREATE TABLE IF NOT EXISTS product_variants (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(200) NOT NULL,
            sku VARCHAR(100) UNIQUE NOT NULL,
            stock INT,
            id_product INT,
            id_tax_category INT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP NULL, 
            FOREIGN KEY (id_product) REFERENCES products(id),
            FOREIGN KEY (id_tax_category) REFERENCES tax_categories(id)
        )
    `
    if _, err := m.db.Exec(createProductVariantsQuery); err != nil {
        return fmt.Errorf("error creating product_variants table: %w", err)
    }

    // Crear tabla product_variant_attribute_values
    createVariantAttributeValuesQuery := `
        CREATE TABLE IF NOT EXISTS n_product_variant_attribute_values (
            id_attribute_value INT,
            id_product_variant INT,
            PRIMARY KEY (id_attribute_value, id_product_variant),
            FOREIGN KEY (id_attribute_value) REFERENCES attribute_values(id),
            FOREIGN KEY (id_product_variant) REFERENCES product_variants(id)
        )
    `
    if _, err := m.db.Exec(createVariantAttributeValuesQuery); err != nil {
        return fmt.Errorf("error creating product_variant_attribute_values table: %w", err)
    }

    // Crear tabla product_variant_media
    createVariantMediaQuery := `
        CREATE TABLE IF NOT EXISTS n_product_variant_media (
            id_product_variant INT,
            id_media INT,
            PRIMARY KEY (id_media, id_product_variant),
            FOREIGN KEY (id_media) REFERENCES media(id),
            FOREIGN KEY (id_product_variant) REFERENCES product_variants(id)
        )
    `
    if _, err := m.db.Exec(createVariantMediaQuery); err != nil {
        return fmt.Errorf("error creating product_variant_media table: %w", err)
    }

    // Crear tabla product_dimensions
    createProductDimensionsQuery := `
        CREATE TABLE IF NOT EXISTS product_dimensions (
            id_product_variant INT PRIMARY KEY,
            length DECIMAL(10,2),
            width DECIMAL(10,2),
            height DECIMAL(10,2),
            weight DECIMAL(10,2),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (id_product_variant) REFERENCES product_variants(id)
        )
    `
    if _, err := m.db.Exec(createProductDimensionsQuery); err != nil {
        return fmt.Errorf("error creating product_dimensions table: %w", err)
    }

    fmt.Println("Products system tables ready")
    return nil
}