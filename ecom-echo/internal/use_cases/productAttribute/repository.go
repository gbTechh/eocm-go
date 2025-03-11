package productattribute

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
	attributevalue "ecom/internal/use_cases/attributeValue"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, c *ProductAttribute) error
	GetByID(ctx context.Context, id int64) (*ProductAttribute, error)
	Update(ctx context.Context, c *ProductAttribute) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]ProductAttribute, int64, error)
	
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

const (
	queryInsertProductAttribute = `
			INSERT INTO product_attributes (name, type, description, created_at, updated_at)
			VALUES (?, ?, ?, NOW(), NOW())
	`

	queryGetProductAttribute = `
			SELECT id, name, type, description, created_at, updated_at
			FROM product_attributes
			WHERE id = ? AND deleted_at IS NULL
	`

	queryUpdateCategory = `
			UPDATE product_attributes 
			SET name = ?, type = ?, description = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`

	queryDeleteCategory = `
			UPDATE product_attributes
			SET deleted_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
)

func (r *MySQLRepository) Create(ctx context.Context, c *ProductAttribute) error {
	result, err := r.db.ExecContext(ctx, queryInsertProductAttribute,
			c.Name, c.Type, c.Description)
	if err != nil {
		return errors.NewMysqlError(err)	
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.NewInternalError("Error al obtener el id creado", err)
	}

	c.ID = id
	return nil
}

func (r *MySQLRepository) GetByID2(ctx context.Context, id int64) (*ProductAttribute, error) {
	c := &ProductAttribute{}

	err := r.db.QueryRowContext(ctx, queryGetProductAttribute, id).Scan(
			&c.ID, &c.Name, &c.Type, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Attributo de producto no encontrado")
	}
	if err != nil {
			return nil, err
	}

	return c, nil
}

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*ProductAttribute, error) {
    // Primero obtenemos el product attribute
    paQuery := `
        SELECT id, name, type, description, created_at, updated_at
        FROM product_attributes
        WHERE id = ? AND deleted_at IS NULL
    `
    pa := &ProductAttribute{}
    err := r.db.QueryRowContext(ctx, paQuery, id).Scan(
        &pa.ID, &pa.Name, &pa.Type, &pa.Description, 
        &pa.CreatedAt, &pa.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, errors.NewNotFoundError("Attributo de producto no encontrado")
    }
    if err != nil {
        return nil, err
    }

    // Luego obtenemos sus valores
    valuesQuery := `
        SELECT id, name, id_attribute_product
        FROM attribute_values
        WHERE id_attribute_product = ?
    `
    rows, err := r.db.QueryContext(ctx, valuesQuery, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var value attributevalue.AttributeValue
        err := rows.Scan(&value.ID, &value.Name, &value.IDProductAttribute)
        if err != nil {
            return nil, err
        }
        pa.Values = append(pa.Values, value)
    }

    return pa, nil
}
func (r *MySQLRepository) Update(ctx context.Context, c *ProductAttribute) error {
	result, err := r.db.ExecContext(ctx, queryUpdateCategory, c.Name, c.Type, c.Description, c.ID)
	if err != nil {
		return errors.NewMysqlError(err)	
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Error al obtener la el producto atributo", err)
	}
	if rows == 0 {
		return errors.NewNotFoundError("Atributo de producto no encontrado")
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, queryDeleteCategory, id)
	if err != nil {
			return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return err
	}
	if rows == 0 {
			return errors.NewNotFoundError("Atributo de producto no encontrado")
	}
	return nil
}


func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]ProductAttribute, int64, error) {
	// Query base para contar
	countQuery := `
			SELECT COUNT(*) 
			FROM product_attributes 
			WHERE deleted_at IS NULL
	`
	
	// Query base para listar con JOIN a media
	listQuery := `
			SELECT id, name, type, description, created_at, updated_at
			FROM product_attributes
			WHERE deleted_at IS NULL
	`

	args := []interface{}{}

	// Agregar condiciones de búsqueda
	if p.Search != "" {
			countQuery += " AND (name LIKE ? OR type LIKE ?)"
			listQuery += " AND (name LIKE ? OR type LIKE ?)"
			searchTerm := "%" + p.Search + "%"
			args = append(args, searchTerm, searchTerm)
	}
	

	// Contar total
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	// Agregar paginación
	listQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	 // Debug: imprimir query y argumentos
	fmt.Printf("Query: %s\n", listQuery)
	fmt.Printf("Args: %v\n", args)

	// Ejecutar query principal
	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var productAttribute []ProductAttribute
	for rows.Next() {
			c := ProductAttribute{}

	
			err := rows.Scan(
					&c.ID, &c.Name, &c.Type, &c.Description,
					&c.CreatedAt, &c.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}		
			// Obtener valores para cada product attribute
			valuesQuery := `
					SELECT id, name 
					FROM attribute_values 
					WHERE id_attribute_product = ?
			`
			valueRows, err := r.db.QueryContext(ctx, valuesQuery, c.ID)
			if err != nil {
					return nil, 0, err
			}
			defer valueRows.Close()

			for valueRows.Next() {
					var value attributevalue.AttributeValue
					err := valueRows.Scan(&value.ID, &value.Name)
					if err != nil {
							return nil, 0, err
					}
					c.Values = append(c.Values, value)
			}	
			productAttribute = append(productAttribute, c)
	}

	if err = rows.Err(); err != nil {
			return nil, 0, err
	}

	return productAttribute, total, nil
}