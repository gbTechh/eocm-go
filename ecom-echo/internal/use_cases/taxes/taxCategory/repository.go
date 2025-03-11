// repository.go
package taxcategory

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
	"fmt"
)



type Repository interface {
	Create(ctx context.Context, t *TaxCategory) error
	GetByID(ctx context.Context, id int64) (*TaxCategory, error)
	GetTaxDefault(ctx context.Context) (*TaxCategory, error)
	Update(ctx context.Context, t *TaxCategory) error
	UpdateAllDefaultFalse(ctx context.Context) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]TaxCategory, int64, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(ctx context.Context, t *TaxCategory) error {
	query := `
			INSERT INTO tax_categories (name, description, status, is_default, created_at, updated_at)
        VALUES (?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query, t.Name, t.Description,t.Status, t.IsDefault)
	if err != nil {
		return errors.NewMysqlError(err)	
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.NewInternalError("Error al obtener el id creado", err)
	}

	t.ID = id
	return nil
}

func (r *MySQLRepository) GetTaxDefault(ctx context.Context) (*TaxCategory, error) {
	query := `
		SELECT id, name, description, status, is_default, created_at, updated_at
		FROM tax_categories
		WHERE is_default = 1 AND deleted_at IS NULL
	`
	tc := &TaxCategory{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&tc.ID, &tc.Name, &tc.Description, &tc.Status, &tc.IsDefault,
		&tc.CreatedAt, &tc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return tc, nil
	}
	if err != nil {
		return nil, err
	}
	return tc, nil
}

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*TaxCategory, error) {
	query := `
		SELECT id, name, description, status, is_default, created_at, updated_at
		FROM tax_categories
		WHERE id = ? AND deleted_at IS NULL
	`
	tc := &TaxCategory{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tc.ID, &tc.Name, &tc.Description, &tc.Status, &tc.IsDefault,
		&tc.CreatedAt, &tc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("Categoría de impuesto no encontrada")
	}
	if err != nil {
		return nil, err
	}
	return tc, nil
}

func (r *MySQLRepository) Update(ctx context.Context, tc *TaxCategory) error {
	query := `
		UPDATE tax_categories 
		SET name = ?, description = ?, status = ?, is_default = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, tc.Name, tc.Description, tc.Status, tc.IsDefault, tc.ID)
	if err != nil {
		return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Error al actualizar categoría de impuesto", err)
	}
	if rows == 0 {
		return errors.NewNotFoundError("Categoría de impuesto no encontrada")
	}
	return nil
}
func (r *MySQLRepository) UpdateAllDefaultFalse(ctx context.Context) error {
	query := `
		UPDATE tax_categories 
		SET is_default = 0
		WHERE deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Error al actualizar el campo is_default de todsa las categorías de impuesto", err)
	}
	fmt.Printf("Filas afectadas: %v\n", rows)
	
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE tax_categories
			SET deleted_at = NOW(), is_default = 0
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
			return errors.NewNotFoundError("Categoría de impuesto no encontrada")
	}
	return nil
}

func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]TaxCategory, int64, error) {
	countQuery := `
		SELECT COUNT(*) 
		FROM tax_categories 
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if p.Search != "" {
		countQuery += " AND name LIKE ?"
		args = append(args, "%"+p.Search+"%")
	}

	if p.Status != nil {
		countQuery += " AND status = ?"
		args = append(args, *p.Status)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, description, status, is_default, created_at, updated_at
		FROM tax_categories
		WHERE deleted_at IS NULL
	`

	if p.Search != "" {
		query += " AND name LIKE ?"
	}

	if p.Status != nil {
		query += " AND status = ?"
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []TaxCategory
	for rows.Next() {
		var tc TaxCategory
		err := rows.Scan(
				&tc.ID, &tc.Name, &tc.Description, &tc.Status, &tc.IsDefault,
				&tc.CreatedAt, &tc.UpdatedAt,
		)
		if err != nil {
				return nil, 0, err
		}
		categories = append(categories, tc)
	}

	return categories, total, nil
}
