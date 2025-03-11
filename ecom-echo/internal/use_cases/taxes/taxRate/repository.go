// repository.go
package taxrate

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
)

type Repository interface {
	Create(ctx context.Context, tr *TaxRate) error
	GetByID(ctx context.Context, id int64) (*TaxRate, error)
	Update(ctx context.Context, tr *TaxRate) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]TaxRate, int64, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(ctx context.Context, tr *TaxRate) error {
	query := `
			INSERT INTO tax_rates 
			(id_tax_category, id_tax_zone, percentage, is_default, status, name, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query,
			tr.IDTaxCategory, tr.IDTaxZone, tr.Percentage,
			tr.IsDefault, tr.Status, tr.Name)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
			return errors.NewInternalError("Error al obtener el id creado", err)
	}

	tr.ID = id
	return nil
}

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*TaxRate, error) {
	query := `
			SELECT tr.id, tr.id_tax_category, tr.id_tax_zone, tr.percentage, 
							tr.is_default, tr.status, tr.name, tr.created_at, tr.updated_at,
							tc.id, tc.name,
							tz.id, tz.name
			FROM tax_rates tr
			LEFT JOIN tax_categories tc ON tc.id = tr.id_tax_category
			LEFT JOIN tax_zones tz ON tz.id = tr.id_tax_zone
			WHERE tr.id = ? AND tr.deleted_at IS NULL
	`
	tr := &TaxRate{
			TaxCategory: &TaxCategory{},
			TaxZone:    &TaxZone{},
	}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&tr.ID, &tr.IDTaxCategory, &tr.IDTaxZone, &tr.Percentage,
			&tr.IsDefault, &tr.Status, &tr.Name, &tr.CreatedAt, &tr.UpdatedAt,
			&tr.TaxCategory.ID, &tr.TaxCategory.Name,
			&tr.TaxZone.ID, &tr.TaxZone.Name,
	)
	
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Tasa de impuesto no encontrada")
	}
	if err != nil {
			return nil, err
	}
	
	return tr, nil
}

func (r *MySQLRepository) Update(ctx context.Context, tr *TaxRate) error {
	query := `
			UPDATE tax_rates 
			SET id_tax_category = ?, id_tax_zone = ?, percentage = ?,
					is_default = ?, status = ?, name = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
			tr.IDTaxCategory, tr.IDTaxZone, tr.Percentage,
			tr.IsDefault, tr.Status, tr.Name, tr.ID)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return errors.NewInternalError("Error al actualizar tasa de impuesto", err)
	}
	if rows == 0 {
			return errors.NewNotFoundError("Tasa de impuesto no encontrada")
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE tax_rates
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
			return errors.NewNotFoundError("Tasa de impuesto no encontrada")
	}
	return nil
}

func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]TaxRate, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM tax_rates 
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
			SELECT tr.id, tr.id_tax_category, tr.id_tax_zone, tr.percentage, 
							tr.is_default, tr.status, tr.name, tr.created_at, tr.updated_at,
							tc.id, tc.name,
							tz.id, tz.name
			FROM tax_rates tr
			LEFT JOIN tax_categories tc ON tc.id = tr.id_tax_category
			LEFT JOIN tax_zones tz ON tz.id = tr.id_tax_zone
			WHERE tr.deleted_at IS NULL
	`

	if p.Search != "" {
			query += " AND tr.name LIKE ?"
	}

	if p.Status != nil {
			query += " AND tr.status = ?"
	}

	query += " ORDER BY tr.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var rates []TaxRate
	for rows.Next() {
			tr := TaxRate{
					TaxCategory: &TaxCategory{},
					TaxZone:    &TaxZone{},
			}
			err := rows.Scan(
					&tr.ID, &tr.IDTaxCategory, &tr.IDTaxZone, &tr.Percentage,
					&tr.IsDefault, &tr.Status, &tr.Name, &tr.CreatedAt, &tr.UpdatedAt,
					&tr.TaxCategory.ID, &tr.TaxCategory.Name,
					&tr.TaxZone.ID, &tr.TaxZone.Name,
			)
			if err != nil {
					return nil, 0, err
			}
			rates = append(rates, tr)
	}

	return rates, total, nil
}