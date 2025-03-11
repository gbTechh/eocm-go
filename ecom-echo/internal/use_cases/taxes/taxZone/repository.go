package taxzone

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
)

type Repository interface {
	Create(ctx context.Context, tz *TaxZone) error
	GetByID(ctx context.Context, id int64) (*TaxZone, error)
	Update(ctx context.Context, tz *TaxZone) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]TaxZone, int64, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(ctx context.Context, tz *TaxZone) error {
	query := `
			INSERT INTO tax_zones 
			(id_zone, name, description, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query, 
			tz.IDZone, tz.Name, tz.Description, tz.Status)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
			return errors.NewInternalError("Error al obtener el id creado", err)
	}

	tz.ID = id
	return nil
}

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*TaxZone, error) {
	query := `
			SELECT tz.id, tz.id_zone, tz.name, tz.description, tz.status, 
							tz.created_at, tz.updated_at,
							z.id, z.name
			FROM tax_zones tz
			LEFT JOIN zones z ON z.id = tz.id_zone
			WHERE tz.id = ? AND tz.deleted_at IS NULL
	`
	tz := &TaxZone{Zone: &Zone{}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&tz.ID, &tz.IDZone, &tz.Name, &tz.Description, &tz.Status,
			&tz.CreatedAt, &tz.UpdatedAt,
			&tz.Zone.ID, &tz.Zone.Name,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Zona fiscal no encontrada")
	}
	if err != nil {
			return nil, err
	}
	return tz, nil
}

func (r *MySQLRepository) Update(ctx context.Context, tz *TaxZone) error {
	query := `
			UPDATE tax_zones 
			SET id_zone = ?, name = ?, description = ?, status = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, 
			tz.IDZone, tz.Name, tz.Description, tz.Status, tz.ID)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return errors.NewInternalError("Error al actualizar zona fiscal", err)
	}
	if rows == 0 {
			return errors.NewNotFoundError("Zona fiscal no encontrada")
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE tax_zones
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
			return errors.NewNotFoundError("Zona fiscal no encontrada")
	}
	return nil
}

func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]TaxZone, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM tax_zones 
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
			SELECT tz.id, tz.id_zone, tz.name, tz.description, tz.status, 
							tz.created_at, tz.updated_at,
							z.id, z.name
			FROM tax_zones tz
			LEFT JOIN zones z ON z.id = tz.id_zone
			WHERE tz.deleted_at IS NULL
	`

	if p.Search != "" {
			query += " AND tz.name LIKE ?"
	}

	if p.Status != nil {
			query += " AND tz.status = ?"
	}

	query += " ORDER BY tz.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var zones []TaxZone
	for rows.Next() {
			var tz TaxZone
			tz.Zone = &Zone{}
			err := rows.Scan(
					&tz.ID, &tz.IDZone, &tz.Name, &tz.Description, &tz.Status,
					&tz.CreatedAt, &tz.UpdatedAt,
					&tz.Zone.ID, &tz.Zone.Name,
			)
			if err != nil {
					return nil, 0, err
			}
			zones = append(zones, tz)
	}

	return zones, total, nil
}