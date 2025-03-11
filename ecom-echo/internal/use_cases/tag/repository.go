// repository.go
package tag

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
)



type Repository interface {
	Create(ctx context.Context, t *Tag) error
	GetByID(ctx context.Context, id int64) (*Tag, error)
	Update(ctx context.Context, t *Tag) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]Tag, int64, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(ctx context.Context, t *Tag) error {
	query := `
			INSERT INTO tags (name, code, active, created_at, updated_at)
			VALUES (?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query, t.Name, t.Code, t.Active)
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

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*Tag, error) {
	query := `
			SELECT id, name, code, active, created_at, updated_at
			FROM tags
			WHERE id = ? AND deleted_at IS NULL
	`
	t := &Tag{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&t.ID, &t.Name, &t.Code, &t.Active, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Tag no encontrado")
	}
	if err != nil {
			return nil, err
	}
	return t, nil
}

func (r *MySQLRepository) Update(ctx context.Context, t *Tag) error {
	query := `
			UPDATE tags 
			SET name = ?, code = ?, active = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, t.Name, t.Code, t.Active, t.ID)
	if err != nil {
		return errors.NewMysqlError(err)	
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Error al obtener tag", err)
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE tags
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

func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]Tag, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM tags 
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
			SELECT id, name, code, active, created_at, updated_at
			FROM tags
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

	var tags []Tag
	for rows.Next() {
			var t Tag
			err := rows.Scan(
					&t.ID, &t.Name, &t.Code, &t.Active,
					&t.CreatedAt, &t.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}
			tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
			return nil, 0, err
	}

	return tags, total, nil
}