// repository.go
package product

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
	"fmt"
	"strings"

	at "ecom/internal/use_cases/attributeValue"
	category "ecom/internal/use_cases/category"
	tag "ecom/internal/use_cases/tag"
)

type Repository interface {
	// Operaciones principales de Producto
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]Product, int64, error)
	
	// Operaciones de Variantes
	CreateVariant(ctx context.Context, v *ProductVariant) error
	UpdateVariant(ctx context.Context, v *ProductVariant) error
	DeleteVariant(ctx context.Context, id int64) error
	GetVariantByID(ctx context.Context, id int64) (*ProductVariant, error)
	ListVariants(ctx context.Context, productID int64) ([]ProductVariant, error)
	
	// Métodos auxiliares para relaciones
	AddProductTags(ctx context.Context, tx *sql.Tx, productID int64, tagIDs []int64) error
	AddVariantAttributeValues(ctx context.Context, tx *sql.Tx, variantID int64, values []at.AttributeValue) error
	
	// Métodos para cargar relaciones
	loadProductTags(ctx context.Context, p *Product) error
	loadProductVariants(ctx context.Context, p *Product) error
	loadVariantAttributeValues(ctx context.Context, v *ProductVariant) error
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(ctx context.Context, p *Product) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	
	// Crear producto
	query := `
			INSERT INTO products (name, slug, description, status, id_category, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := tx.ExecContext(ctx, query,
			p.Name, p.Slug, p.Description, p.Status, p.IDCategory)
	if err != nil {
		tx.Rollback()
		return errors.NewMysqlError(err)
	}

	productID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return errors.NewInternalError("Error al obtener el id creado", err)
	}
	p.ID = productID
	
	// Crear relaciones con tags si existen
	if len(p.Tags) > 0 {
		tagIDs := make([]int64, len(p.Tags))
		for i, tag := range p.Tags {
			tagIDs[i] = tag.ID
		}
		if err := r.AddProductTags(ctx, tx, productID, tagIDs); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
	// Consulta principal para obtener el producto
	query := `
			SELECT p.id, p.name, p.slug, p.description, p.status, p.id_category,
							p.created_at, p.updated_at,
							c.id, c.name
			FROM products p
			LEFT JOIN categories c ON c.id = p.id_category
			WHERE p.id = ? AND p.deleted_at IS NULL
	`
	
	p := &Product{
		Category: &category.Category{},
	}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.Status, &p.IDCategory,
		&p.CreatedAt, &p.UpdatedAt,
		&p.Category.ID, &p.Category.Name,
	)
	
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Producto no encontrado")
	}
	if err != nil {
			return nil, err
	}

	// Obtener tags
	tagsQuery := `
			SELECT t.id, t.name, t.code, t.active
			FROM tags t
			INNER JOIN n_products_tags pt ON pt.id_tag = t.id
			WHERE pt.id_product = ? AND t.deleted_at IS NULL
	`
	tagsRows, err := r.db.QueryContext(ctx, tagsQuery, id)
	if err != nil {
			return nil, err
	}
	defer tagsRows.Close()

	for tagsRows.Next() {
			var tag tag.Tag
			if err := tagsRows.Scan(&tag.ID, &tag.Name, &tag.Code, &tag.Active); err != nil {
					return nil, err
			}
			p.Tags = append(p.Tags, tag)
	}

	// Obtener variantes y sus attribute values
	variantsQuery := `
			SELECT v.id, v.name, v.sku, v.stock, v.id_tax_category,
						v.id_product,	v.created_at, v.updated_at
			FROM product_variants v
			WHERE v.id_product = ? AND v.deleted_at IS NULL
	`
	variantsRows, err := r.db.QueryContext(ctx, variantsQuery, id)
	if err != nil {
			return nil, err
	}
	defer variantsRows.Close()

	for variantsRows.Next() {
			var variant ProductVariant
			var idTaxCategory sql.NullInt64
			err := variantsRows.Scan(
					&variant.ID, &variant.Name, &variant.SKU, &variant.Stock,
					&idTaxCategory, &variant.IDProduct, &variant.CreatedAt, &variant.UpdatedAt,
			)
			if err != nil {
					return nil, err
			}

			if idTaxCategory.Valid {
				variant.IDTaxCategory = idTaxCategory.Int64
			} else {
				variant.IDTaxCategory = 0
			}

			// Obtener attribute values para cada variante
			attrQuery := `
					SELECT av.id, av.name, av.id_attribute_product,
									pa.id, pa.name, pa.type, pa.description
					FROM attribute_values av
					INNER JOIN n_product_variant_attribute_values pvav ON pvav.id_attribute_value = av.id
					INNER JOIN product_attributes pa ON pa.id = av.id_attribute_product
					WHERE pvav.id_product_variant = ?
			`
			attrRows, err := r.db.QueryContext(ctx, attrQuery, variant.ID)
			if err != nil {
					return nil, err
			}
			defer attrRows.Close()

			for attrRows.Next() {
					var av at.AttributeValue
        	av.ProductAttribute = &at.ProductAttributeInfo{}
					
					err := attrRows.Scan(
							&av.ID, &av.Name, &av.IDProductAttribute,
							&av.ProductAttribute.ID, &av.ProductAttribute.Name,
							&av.ProductAttribute.Type, &av.ProductAttribute.Description,
					)
					if err != nil {
							return nil, err
					}
					
					variant.AttributeValues = append(variant.AttributeValues, av)
			}

			p.Variants = append(p.Variants, variant)
	}

	return p, nil
}

func (r *MySQLRepository) Update(ctx context.Context, p *Product) error {
	query := `
			UPDATE products 
			SET name = ?, slug = ?, description = ?, status = ?, 
					id_category = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
			p.Name, p.Slug, p.Description, p.Status,
			p.IDCategory, p.ID)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return errors.NewInternalError("Error al actualizar producto", err)
	}
	if rows == 0 {
			return errors.NewNotFoundError("Producto no encontrado")
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE products
			SET deleted_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
			return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return err
	}
	if rows == 0 {
			return errors.NewNotFoundError("Producto no encontrado")
	}
	return nil
}

func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]Product, int64, error) {
	// Consulta para contar el total
	countQuery := `
			SELECT COUNT(*) 
			FROM products 
			WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if p.Search != "" {
			countQuery += " AND (name LIKE ? OR description LIKE ?)"
			searchTerm := "%" + p.Search + "%"
			args = append(args, searchTerm, searchTerm)
	}

	if p.Status != nil {
			countQuery += " AND status = ?"
			args = append(args, *p.Status)
	}

	if p.IDCategory != nil {
			countQuery += " AND id_category = ?"
			args = append(args, *p.IDCategory)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	// Consulta principal
	query := `
			SELECT DISTINCT p.id, p.name, p.slug, p.description, p.status, 
							p.id_category, p.created_at, p.updated_at,
							c.id, c.name
			FROM products p
			LEFT JOIN categories c ON c.id = p.id_category
			WHERE p.deleted_at IS NULL
	`

	if p.Search != "" {
			query += " AND (p.name LIKE ? OR p.description LIKE ?)"
	}

	if p.Status != nil {
			query += " AND p.status = ?"
	}

	if p.IDCategory != nil {
			query += " AND p.id_category = ?"
	}

	query += " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
			p := Product{
					Category: &category.Category{},
			}
			err := rows.Scan(
					&p.ID, &p.Name, &p.Slug, &p.Description, &p.Status,
					&p.IDCategory, &p.CreatedAt, &p.UpdatedAt,
					&p.Category.ID, &p.Category.Name,
			)
			if err != nil {
					return nil, 0, err
			}

			// Cargar tags
			if err := r.loadProductTags(ctx, &p); err != nil {
					return nil, 0, err
			}

			// Cargar variantes con sus attribute values
			if err := r.loadProductVariants(ctx, &p); err != nil {
					return nil, 0, err
			}

			products = append(products, p)
	}

	return products, total, nil
}

func (r *MySQLRepository) CreateVariant(ctx context.Context, v *ProductVariant) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	var idTaxCategory sql.NullInt64

	// Si el valor de v.IDTaxCategory es 0, lo asignamos como NULL
	if v.IDTaxCategory == 0 {
		idTaxCategory = sql.NullInt64{Valid: false} // Indicamos que es NULL
	} else {
		idTaxCategory = sql.NullInt64{Int64: v.IDTaxCategory, Valid: true} // Valor válido
	}

	query := `
			INSERT INTO product_variants 
			(name, sku, stock, id_product, id_tax_category, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	
	result, err := tx.ExecContext(ctx, query,
			v.Name, v.SKU, v.Stock, v.IDProduct, idTaxCategory)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	variantID, err := result.LastInsertId()
	if err != nil {
			tx.Rollback()
			return errors.NewInternalError("Error al obtener el id creado", err)
	}
	v.ID = variantID

	// Crear attribute values
	if len(v.AttributeValues) > 0 {
			query := `INSERT INTO n_product_variant_attribute_values 
								(id_product_variant, id_attribute_value) VALUES `
			vals := make([]string, len(v.AttributeValues))
			args := make([]interface{}, 0, len(v.AttributeValues)*2)
			
			for i, av := range v.AttributeValues {
				fmt.Printf("av: %v\n",av.ID)
				fmt.Printf("varianteID: %v\n",variantID)
					vals[i] = "(?, ?)"
					args = append(args, variantID, av.ID)
			}
			
			query += strings.Join(vals, ",")
			
			_, err = tx.ExecContext(ctx, query, args...)
			if err != nil {
					tx.Rollback()
					return errors.NewMysqlError(err)
			}
	}

	if err = tx.Commit(); err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	return nil
}
func (r *MySQLRepository) UpdateVariant(ctx context.Context, v *ProductVariant) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	var idTaxCategory sql.NullInt64

	// Si el valor de v.IDTaxCategory es 0, lo asignamos como NULL
	if v.IDTaxCategory == 0 {
		idTaxCategory = sql.NullInt64{Valid: false} // Indicamos que es NULL
	} else {
		idTaxCategory = sql.NullInt64{Int64: v.IDTaxCategory, Valid: true} // Valor válido
	}

	query := `
			UPDATE product_variants 
			SET name = ?, sku = ?, stock = ?, id_tax_category = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := tx.ExecContext(ctx, query,
			v.Name, v.SKU, v.Stock, idTaxCategory, v.ID)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			tx.Rollback()
			return errors.NewInternalError("Error al actualizar variante", err)
	}
	if rows == 0 {
			tx.Rollback()
			return errors.NewNotFoundError("Variante no encontrada")
	}

	// Si hay nuevos attribute values, primero eliminamos los existentes
	if len(v.AttributeValues) > 0 {
		_, err = tx.ExecContext(ctx, 
				"DELETE FROM n_product_variant_attribute_values WHERE id_product_variant = ?", 
				v.ID)
		if err != nil {
				tx.Rollback()
				return err
		}
		
		
		if err := r.AddVariantAttributeValues(ctx, tx, v.ID, v.AttributeValues); err != nil {
			tx.Rollback()
			return err
		}		
	} else {
		_, err = tx.ExecContext(ctx, 
				"DELETE FROM n_product_variant_attribute_values WHERE id_product_variant = ?", 
				v.ID)
		if err != nil {
				tx.Rollback()
				return err
		}
	}


	return tx.Commit()
}

func (r *MySQLRepository) DeleteVariant(ctx context.Context, id int64) error {
	query := `
			UPDATE product_variants
			SET deleted_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
			return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return err
	}
	if rows == 0 {
			return errors.NewNotFoundError("Variante no encontrada")
	}
	return nil
}


func (r *MySQLRepository) AddProductTags(ctx context.Context, tx *sql.Tx, productID int64, tagIDs []int64) error {
	query := `INSERT INTO n_products_tags (id_product, id_tag) VALUES `
	values := make([]string, len(tagIDs))
	args := make([]interface{}, 0, len(tagIDs)*2)
	
	for i, tagID := range tagIDs {
			values[i] = "(?, ?)"
			args = append(args, productID, tagID)
	}
	
	query += strings.Join(values, ",")
	
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
			return errors.NewMysqlError(err)
	}
	
	return nil
}
func (r *MySQLRepository) AddVariantAttributeValues(ctx context.Context, tx *sql.Tx, variantID int64, values []at.AttributeValue) error {
	query := `INSERT INTO n_product_variant_attribute_values (id_product_variant, id_attribute_value) VALUES `
	vals := make([]string, len(values))
	args := make([]interface{}, 0, len(values)*2)
	
	for i, av := range values {
			vals[i] = "(?, ?)"
			args = append(args, variantID, av.ID)
	}
	
	query += strings.Join(vals, ",")
	
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
			return errors.NewMysqlError(err)
	}
	
	return nil
}

func (r *MySQLRepository) GetVariantByID(ctx context.Context, id int64) (*ProductVariant, error) {
    query := `
        SELECT v.id, v.name, v.sku, v.stock, v.id_product, v.id_tax_category,
               v.created_at, v.updated_at
        FROM product_variants v
        WHERE v.id = ? AND v.deleted_at IS NULL
    `
    
    variant := &ProductVariant{}
		var idTaxCategory sql.NullInt64
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &variant.ID, &variant.Name, &variant.SKU, &variant.Stock,
        &variant.IDProduct, &idTaxCategory,
        &variant.CreatedAt, &variant.UpdatedAt,
    )
		if idTaxCategory.Valid {
			variant.IDTaxCategory = idTaxCategory.Int64
		} else {
			variant.IDTaxCategory = 0
		}
    
    if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Variante no encontrada")
    }
    if err != nil {
			fmt.Printf("ERROR SS: %v\n",err)
			return nil, err
    }

    // Cargar attribute values
    if err := r.loadVariantAttributeValues(ctx, variant); err != nil {
			return nil, err
    }

    return variant, nil
}

func (r *MySQLRepository) ListVariants(ctx context.Context, productID int64) ([]ProductVariant, error) {
    query := `
        SELECT v.id, v.name, v.sku, v.stock, v.id_product, v.id_tax_category, v.id_product,
               v.created_at, v.updated_at
        FROM product_variants v
        WHERE v.id_product = ? AND v.deleted_at IS NULL
        ORDER BY v.created_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, productID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var variants []ProductVariant
    for rows.Next() {
        var v ProductVariant
				var idTaxCategory sql.NullInt64
        err := rows.Scan(
            &v.ID, &v.Name, &v.SKU, &v.Stock,
            &v.IDProduct, &idTaxCategory, &v.IDProduct,
            &v.CreatedAt, &v.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
				if idTaxCategory.Valid {
					v.IDTaxCategory = idTaxCategory.Int64
				} else {
					v.IDTaxCategory = 0
				}

        // Cargar attribute values para cada variante
        if err := r.loadVariantAttributeValues(ctx, &v); err != nil {
            return nil, err
        }

        variants = append(variants, v)
    }

    return variants, nil
}

func (r *MySQLRepository) loadProductTags(ctx context.Context, p *Product) error {
	query := `
			SELECT t.id, t.name, t.code, t.active
			FROM tags t
			INNER JOIN n_products_tags pt ON pt.id_tag = t.id
			WHERE pt.id_product = ? AND t.deleted_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query, p.ID)
	if err != nil {
			return err
	}
	defer rows.Close()

	for rows.Next() {
			var tag tag.Tag
			if err := rows.Scan(&tag.ID, &tag.Name, &tag.Code, &tag.Active); err != nil {
					return err
			}
			p.Tags = append(p.Tags, tag)
	}
	return nil
}

func (r *MySQLRepository) loadProductVariants(ctx context.Context, p *Product) error {
	query := `
			SELECT v.id, v.name, v.sku, v.stock, v.id_tax_category, v.id_product,
							v.created_at, v.updated_at
			FROM product_variants v
			WHERE v.id_product = ? AND v.deleted_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query, p.ID)
	if err != nil {
			return err
	}
	defer rows.Close()

	for rows.Next() {
			var v ProductVariant
			var idTaxCategory sql.NullInt64
			err := rows.Scan(
					&v.ID, &v.Name, &v.SKU, &v.Stock,
					&idTaxCategory, &v.IDProduct, &v.CreatedAt, &v.UpdatedAt,
			)
			if err != nil {
				return err
			}

			if idTaxCategory.Valid {
				v.IDTaxCategory = idTaxCategory.Int64
			} else {
				v.IDTaxCategory = 0
			}

			// Cargar attribute values para cada variante
			if err := r.loadVariantAttributeValues(ctx, &v); err != nil {
					return err
			}

			p.Variants = append(p.Variants, v)
	}
	return nil
}

func (r *MySQLRepository) loadVariantAttributeValues(ctx context.Context, v *ProductVariant) error {
	query := `
			SELECT av.id, av.name, av.id_attribute_product,
							pa.id, pa.name, pa.type, pa.description
			FROM attribute_values av
			INNER JOIN n_product_variant_attribute_values pvav ON pvav.id_attribute_value = av.id
			INNER JOIN product_attributes pa ON pa.id = av.id_attribute_product
			WHERE pvav.id_product_variant = ?
	`
	rows, err := r.db.QueryContext(ctx, query, v.ID)
	if err != nil {
			return err
	}
	defer rows.Close()

	for rows.Next() {
			var av at.AttributeValue
			av.ProductAttribute = &at.ProductAttributeInfo{}
			
			err := rows.Scan(
					&av.ID, &av.Name, &av.IDProductAttribute,
					&av.ProductAttribute.ID, &av.ProductAttribute.Name,
					&av.ProductAttribute.Type, &av.ProductAttribute.Description,
			)
			if err != nil {
					return err
			}
			
			v.AttributeValues = append(v.AttributeValues, av)
	}
	return nil
}
