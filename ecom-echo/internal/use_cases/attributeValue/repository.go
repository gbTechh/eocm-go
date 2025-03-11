package attributevalue

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, c *AttributeValue) error
	GetByID(ctx context.Context, id int64) (*AttributeValue, error)
	Update(ctx context.Context, c *AttributeValue) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]AttributeValue, int64, error)
	
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

const (
	queryInsertAttributeValue = `
		INSERT INTO attribute_values 
        (name, id_attribute_product)
        VALUES (?, ?)
	`

	queryGetAttributeValue = `
			SELECT id, name, id_attribute_product
        FROM attribute_values
        WHERE id = ?
	`

	queryUpdateAttributeValue = `
			UPDATE attribute_values 
        SET name = ?
        WHERE id = ? 
	`

	queryDeleteAttributeValue = `
			DELETE FROM attribute_values
			WHERE id = ?
	`
)

func (r *MySQLRepository) Create(ctx context.Context, c *AttributeValue) error {
	result, err := r.db.ExecContext(ctx, queryInsertAttributeValue,
			c.Name, c.IDProductAttribute)
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

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*AttributeValue, error) {
	c := &AttributeValue{}

	err := r.db.QueryRowContext(ctx, queryGetAttributeValue, id).Scan(
			&c.ID, &c.Name, &c.IDProductAttribute,
	)

	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Valor de atributo no encontrado")
	}
	if err != nil {
			return nil, err
	}

	return c, nil
}

func (r *MySQLRepository) Update(ctx context.Context, c *AttributeValue) error {
	result, err := r.db.ExecContext(ctx, queryUpdateAttributeValue, c.Name, c.ID)
	if err != nil {
		return errors.NewMysqlError(err)	
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Error al obtener el valor de atributo", err)
	}
	if rows == 0 {
		return errors.NewNotFoundError("Valor de atributo no encontrado")
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, queryDeleteAttributeValue, id)
	if err != nil {
			return errors.NewMysqlError(err)	
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return err
	}
	if rows == 0 {
			return errors.NewNotFoundError("Valor de atributo no encontrado")
	}
	return nil
}


func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]AttributeValue, int64, error) {
	// Query base para contar
	countQuery := `
			SELECT COUNT(*) 
        FROM attribute_values 
	`
	
	// Query base para listar con JOIN a media
	listQuery := `
			SELECT id, name, id_attribute_product
			FROM attribute_values
	`

	args := []interface{}{}

	// Agregar condiciones de búsqueda
	if p.Search != "" {
			countQuery += " AND (name LIKE ?)"
			listQuery += " AND (name LIKE ?)"
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
	listQuery += "LIMIT ? OFFSET ?"
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

	var attributeValue []AttributeValue
	for rows.Next() {
			c := AttributeValue{}

	
			err := rows.Scan(
					&c.ID, &c.Name, &c.IDProductAttribute,
			)
			if err != nil {
					return nil, 0, err
			}		

			attributeValue = append(attributeValue, c)
	}

	if err = rows.Err(); err != nil {
			return nil, 0, err
	}

	return attributeValue, total, nil
}