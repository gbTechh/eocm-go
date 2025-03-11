// repository.go
package zone

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
)



type Repository interface {
	Create(ctx context.Context, t *Zone) error
	GetByID(ctx context.Context, id int64) (*Zone, error)
	Update(ctx context.Context, t *Zone) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]Zone, int64, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(ctx context.Context, t *Zone) error {
	query := `
			INSERT INTO zones (name, active, created_at, updated_at)
			VALUES (?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query, t.Name, t.Active)
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

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*Zone, error) {
	query := `
			SELECT id, name, active, created_at, updated_at
			FROM zones
			WHERE id = ? AND deleted_at IS NULL
	`
	t := &Zone{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&t.ID, &t.Name, &t.Active, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Zona no encontrado")
	}
	if err != nil {
			return nil, err
	}
	return t, nil
}

func (r *MySQLRepository) Update(ctx context.Context, t *Zone) error {
	query := `
			UPDATE zones 
			SET name = ?, active = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, t.Name, t.Active, t.ID)
	if err != nil {
		return errors.NewMysqlError(err)	
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Error al obtener zona", err)
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE zones
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
			return errors.ErrNotFound
	}
	return nil
}

func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]Zone, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM zones 
			WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if p.Search != "" {
			countQuery += " AND (name LIKE ? OR code LIKE ?)"
			searchTerm := "%" + p.Search + "%"
			args = append(args, searchTerm, searchTerm)
	}

	if p.Active != nil {
			countQuery += " AND active = ?"
			args = append(args, *p.Active)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	query := `
			SELECT id, name, active, created_at, updated_at
			FROM zones
			WHERE deleted_at IS NULL
	`

	if p.Search != "" {
			query += " AND (name LIKE ? OR code LIKE ?)"
	}

	if p.Active != nil {
			query += " AND active = ?"
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var zones []Zone
	for rows.Next() {
			var t Zone
			err := rows.Scan(
					&t.ID, &t.Name, &t.Active,
					&t.CreatedAt, &t.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}
			zones = append(zones, t)
	}

	if err = rows.Err(); err != nil {
			return nil, 0, err
	}

	return zones, total, nil
}