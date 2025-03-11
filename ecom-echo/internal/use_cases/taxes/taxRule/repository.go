// repository.go
package taxrule

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
)

type Repository interface {
	Create(ctx context.Context, tr *TaxRule) error
	GetByID(ctx context.Context, id int64) (*TaxRule, error)
	Update(ctx context.Context, tr *TaxRule) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]TaxRule, int64, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(ctx context.Context, tr *TaxRule) error {
	query := `
			INSERT INTO tax_rules 
			(id_tax_rate, priority, status, min_amount, max_amount, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query,
			tr.IDTaxRate, tr.Priority, tr.Status, tr.MinAmount, tr.MaxAmount)
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

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*TaxRule, error) {
	query := `
			SELECT tr.id, tr.id_tax_rate, tr.priority, tr.status, 
							tr.min_amount, tr.max_amount, tr.created_at, tr.updated_at,
							t.id, t.name, t.percentage
			FROM tax_rules tr
			LEFT JOIN tax_rates t ON t.id = tr.id_tax_rate
			WHERE tr.id = ? AND tr.deleted_at IS NULL
	`
	tr := &TaxRule{TaxRate: &TaxRate{}}
	
	var minAmount, maxAmount sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&tr.ID, &tr.IDTaxRate, &tr.Priority, &tr.Status,
			&minAmount, &maxAmount, &tr.CreatedAt, &tr.UpdatedAt,
			&tr.TaxRate.ID, &tr.TaxRate.Name, &tr.TaxRate.Percentage,
	)
	
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Regla de impuesto no encontrada")
	}
	if err != nil {
			return nil, err
	}

	if minAmount.Valid {
			tr.MinAmount = &minAmount.Float64
	}
	if maxAmount.Valid {
			tr.MaxAmount = &maxAmount.Float64
	}
	
	return tr, nil
}

func (r *MySQLRepository) Update(ctx context.Context, tr *TaxRule) error {
	query := `
			UPDATE tax_rules 
			SET id_tax_rate = ?, priority = ?, status = ?,
					min_amount = ?, max_amount = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
			tr.IDTaxRate, tr.Priority, tr.Status,
			tr.MinAmount, tr.MaxAmount, tr.ID)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return errors.NewInternalError("Error al actualizar regla de impuesto", err)
	}
	if rows == 0 {
			return errors.NewNotFoundError("Regla de impuesto no encontrada")
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE tax_rules
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
			return errors.NewNotFoundError("Regla de impuesto no encontrada")
	}
	return nil
}

func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]TaxRule, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM tax_rules 
			WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if p.Status != nil {
			countQuery += " AND status = ?"
			args = append(args, *p.Status)
	}

	if p.IDTaxRate != nil {
			countQuery += " AND id_tax_rate = ?"
			args = append(args, *p.IDTaxRate)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	query := `
			SELECT tr.id, tr.id_tax_rate, tr.priority, tr.status, 
							tr.min_amount, tr.max_amount, tr.created_at, tr.updated_at,
							t.id, t.name, t.percentage
			FROM tax_rules tr
			LEFT JOIN tax_rates t ON t.id = tr.id_tax_rate
			WHERE tr.deleted_at IS NULL
	`

	if p.Status != nil {
			query += " AND tr.status = ?"
	}

	if p.IDTaxRate != nil {
			query += " AND tr.id_tax_rate = ?"
	}

	query += " ORDER BY tr.priority DESC, tr.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var rules []TaxRule
	for rows.Next() {
			tr := TaxRule{TaxRate: &TaxRate{}}
			var minAmount, maxAmount sql.NullFloat64
			
			err := rows.Scan(
					&tr.ID, &tr.IDTaxRate, &tr.Priority, &tr.Status,
					&minAmount, &maxAmount, &tr.CreatedAt, &tr.UpdatedAt,
					&tr.TaxRate.ID, &tr.TaxRate.Name, &tr.TaxRate.Percentage,
			)
			if err != nil {
					return nil, 0, err
			}

			if minAmount.Valid {
					tr.MinAmount = &minAmount.Float64
			}
			if maxAmount.Valid {
					tr.MaxAmount = &maxAmount.Float64
			}

			rules = append(rules, tr)
	}

	return rules, total, nil
}