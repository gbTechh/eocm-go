// internal/media/repository.go
package media

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("media not found")
)

// Repository interface that defines all media operations
type Repository interface {
	Upload(ctx context.Context, m *Media) error
	GetByID(ctx context.Context, id int64) (*Media, error)
	Update(ctx context.Context, m *Media) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]Media, int64, error)
}

// MySQLRepository implementation of Repository
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a new instance of the repository
func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

// Queries constants
const (
	queryInsertMedia = `
		INSERT INTO media (file_name, file_path, mime_type, size, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	queryGetMediaByID = `
		SELECT id, file_name, file_path, mime_type, size, created_at, updated_at
		FROM media
		WHERE id = ? AND deleted_at IS NULL
	`

	queryUpdateMedia = `
		UPDATE media 
		SET file_name = ?, file_path = ?, mime_type = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	queryDeleteMedia = `
		UPDATE media
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	queryListMedia = `
		SELECT id, file_name, file_path, mime_type, size, created_at, updated_at
		FROM media
		WHERE deleted_at IS NULL
	`

	queryCountMedia = `
		SELECT COUNT(*) 
		FROM media 
		WHERE deleted_at IS NULL
	`
)

// Upload inserts a new media record
func (r *MySQLRepository) Upload(ctx context.Context, m *Media) error {
	result, err := r.db.ExecContext(ctx, queryInsertMedia,
		m.FileName,
		m.FilePath,
		m.MimeType,
		m.Size,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}

// GetByID retrieves a media record by its ID
func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*Media, error) {
	m := &Media{}
	err := r.db.QueryRowContext(ctx, queryGetMediaByID, id).Scan(
		&m.ID,
		&m.FileName,
		&m.FilePath,
		&m.MimeType,
		&m.Size,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return m, nil
}

// Update updates an existing media record
func (r *MySQLRepository) Update(ctx context.Context, m *Media) error {
	result, err := r.db.ExecContext(ctx, queryUpdateMedia,
		m.FileName,
		m.FilePath,
		m.MimeType,
		m.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft deletes a media record
func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, queryDeleteMedia, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}



// List retrieves media records with pagination and filtering
func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]Media, int64, error) {
	// Base queries
	countQuery := queryCountMedia
	listQuery := queryListMedia
	args := []interface{}{}

	// Agregar condición de búsqueda si existe
	if p.Search != "" {
			countQuery += " AND file_name LIKE ?"
			listQuery += " AND file_name LIKE ?"
			args = append(args, "%"+p.Search+"%")
	}

	// Contar total
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	// Agregar orden y paginación
	listQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	offset := (p.Page - 1) * p.PerPage
	args = append(args, p.PerPage, offset)

	// Ejecutar query principal
	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var medias []Media
	for rows.Next() {
			var m Media
			err := rows.Scan(
					&m.ID,
					&m.FileName,
					&m.FilePath,
					&m.MimeType,
					&m.Size,
					&m.CreatedAt,
					&m.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}
			medias = append(medias, m)
	}

	if err = rows.Err(); err != nil {
			return nil, 0, err
	}

	return medias, total, nil
}